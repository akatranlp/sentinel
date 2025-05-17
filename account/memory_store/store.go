package memorystore

import (
	"context"
	"encoding/json"
	"maps"
	"os"
	"slices"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/sentinel/account"
	"github.com/google/uuid"
)

type MemoryUserStore struct {
	accounts map[account.AccountID]account.Account
	users    map[account.UserID]account.User
	savePath string
}

type marshal struct {
	Accounts []account.Account `json:"accounts"`
	Users    []account.User    `json:"users"`
}

func (s *MemoryUserStore) MarshalJSON() ([]byte, error) {
	accounts := slices.Collect(maps.Values(s.accounts))
	users := slices.Collect(maps.Values(s.users))

	return json.Marshal(marshal{
		Accounts: accounts,
		Users:    users,
	})
}

func (s *MemoryUserStore) UnmarshalJSON(data []byte) error {
	var store marshal

	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	accounts := maps.Collect(its.Map12(slices.Values(store.Accounts), func(acc account.Account) (account.AccountID, account.Account) {
		return acc.AccountID, acc
	}))

	users := maps.Collect(its.Map12(slices.Values(store.Users), func(u account.User) (account.UserID, account.User) {
		return u.UserID, u
	}))

	*s = MemoryUserStore{
		accounts: accounts,
		users:    users,
	}
	return nil
}

var _ (account.UserStore) = (*MemoryUserStore)(nil)

func NewMemoryUserStore(filePath ...string) (*MemoryUserStore, error) {
	var path string
	if len(filePath) > 0 {
		path = filePath[0]
	}
	if path != "" {
		f, err := os.Open(path)
		if err == nil {
			defer f.Close()
			var store MemoryUserStore
			if err := json.NewDecoder(f).Decode(&store); err != nil {
				return nil, err
			}
			store.savePath = path
			return &store, nil
		}
	}

	return &MemoryUserStore{
		accounts: make(map[account.AccountID]account.Account),
		users:    make(map[account.UserID]account.User),
		savePath: path,
	}, nil
}

func (s *MemoryUserStore) GetUserByID(ctx context.Context, id account.UserID) (account.User, error) {
	user, ok := s.users[id]
	if !ok {
		return user, account.ErrUserNotFound
	}
	return user, nil
}

func (s *MemoryUserStore) GetAccountsForUserID(ctx context.Context, id account.UserID) ([]account.Account, error) {
	if _, ok := s.users[id]; !ok {
		return nil, account.ErrUserNotFound
	}

	accSeq := maps.Values(s.accounts)
	filterSeq := its.Filter(accSeq, func(v account.Account) bool { return v.UserID == id })
	accounts := slices.Collect(filterSeq)
	if len(accounts) == 0 {
		return nil, account.ErrNoAccountFound
	}

	return accounts, nil
}

func (s *MemoryUserStore) GetUserByAccountID(ctx context.Context, id account.AccountID) (account.User, error) {
	acc, ok := s.accounts[id]
	if !ok {
		return account.User{}, account.ErrAccountNotFound
	}

	user, ok := s.users[acc.UserID]
	if !ok {
		return user, account.ErrUserNotFound
	}
	return user, nil
}

func (s *MemoryUserStore) GetAccountByID(ctx context.Context, accID account.AccountID) (account.Account, error) {
	acc, ok := s.accounts[accID]
	if !ok {
		return acc, account.ErrAccountNotFound
	}
	return acc, nil
}

func (s *MemoryUserStore) GetAccountByProvider(ctx context.Context, userID account.UserID, provider string) (account.Account, error) {
	if _, ok := s.users[userID]; !ok {
		return account.Account{}, account.ErrUserNotFound
	}

	for _, acc := range s.accounts {
		if acc.UserID == userID && acc.Provider == provider {
			return acc, nil
		}
	}

	return account.Account{}, account.ErrAccountNotFound
}

func (s *MemoryUserStore) GetOrCreateUserFromAccount(ctx context.Context, acc account.Account) (account.User, error) {
	if oldAcc, ok := s.accounts[acc.AccountID]; ok {
		acc.UserID = oldAcc.UserID
		s.accounts[acc.AccountID] = acc
		s.saveToFile()

		user, ok := s.users[acc.UserID]
		if !ok {
			return user, account.ErrUserNotFound
		}
		return user, nil
	}

	user := account.User{
		UserID:        account.UserID(uuid.NewString()),
		Name:          acc.Name,
		Email:         acc.Email,
		EmailVerified: acc.EmailVerified,
		Username:      acc.PreferredUsername,
		Picture:       acc.Picture,
	}
	acc.UserID = user.UserID

	s.accounts[acc.AccountID] = acc
	s.users[user.UserID] = user

	s.saveToFile()

	return user, nil
}

func (s *MemoryUserStore) UpdateUser(ctx context.Context, id account.UserID, user account.User) error {
	if _, ok := s.users[id]; !ok {
		return account.ErrUserNotFound
	}

	s.users[user.UserID] = user
	s.saveToFile()
	return nil
}

func (s *MemoryUserStore) UpdateAccount(ctx context.Context, id account.AccountID, acc account.Account) error {
	if _, ok := s.accounts[id]; !ok {
		return account.ErrAccountNotFound
	}

	s.accounts[acc.AccountID] = acc
	s.saveToFile()
	return nil
}

func (s *MemoryUserStore) LinkAccount(ctx context.Context, id account.UserID, acc account.Account) error {
	if _, ok := s.users[id]; !ok {
		return account.ErrUserNotFound
	}

	// if _, ok := s.accounts[acc.AccountID]; ok {
	// 	return account.ErrAccountAlreadyLinked
	// }

	acc.UserID = id
	s.accounts[acc.AccountID] = acc
	s.saveToFile()
	return nil
}

func (s *MemoryUserStore) UnLinkAccount(ctx context.Context, id account.UserID, accID account.AccountID) error {
	user, ok := s.users[id]
	if !ok {
		return account.ErrUserNotFound
	}

	acc, ok := s.accounts[accID]
	if !ok {
		return account.ErrAccountNotLinked
	}

	if acc.UserID != id {
		return account.ErrAccountNotLinked
	}

	accounts, err := s.GetAccountsForUserID(ctx, id)
	if err != nil {
		return err
	}

	if len(accounts) == 1 {
		return account.ErrLastAccount
	}
	delete(s.accounts, accID)

	accounts, _ = s.GetAccountsForUserID(ctx, id)
	newAcc := accounts[0]

	var userNeedsUpdate bool
	if user.Email == acc.Email {
		user.Email = newAcc.Email
		userNeedsUpdate = true
	}
	if user.Name == acc.Name {
		user.Name = newAcc.Name
		userNeedsUpdate = true
	}
	if user.Picture == acc.Picture {
		user.Picture = newAcc.Picture
		userNeedsUpdate = true
	}
	if user.Username == acc.PreferredUsername {
		user.Username = newAcc.PreferredUsername
		userNeedsUpdate = true
	}

	if userNeedsUpdate {
		s.users[id] = user
	}
	s.saveToFile()

	return nil
}

func (s *MemoryUserStore) saveToFile() error {
	if s.savePath == "" {
		return nil
	}
	f, err := os.Create(s.savePath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}
