package basestore

import (
	"context"
	"errors"
	"slices"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/sentinel/account"
)

type Repository interface {
	CreateUser(ctx context.Context, user account.User) (account.User, error)
	GetUserByID(ctx context.Context, id account.UserID) (account.User, error)
	GetUserByAccountID(ctx context.Context, id account.AccountID) (account.User, error)
	UpdateUser(ctx context.Context, id account.UserID, user account.User) error

	CreateAccount(ctx context.Context, acc account.Account) (account.Account, error)
	GetAccountByID(ctx context.Context, id account.AccountID) (account.Account, error)
	GetAccountsByUserID(ctx context.Context, id account.UserID) ([]account.Account, error)
	UpdateAccount(ctx context.Context, id account.AccountID, acc account.Account) error
	DeleteAccountByID(ctx context.Context, id account.AccountID) error
}

type UserIDAndPRoviderGetter interface {
	GetAccountByUserIDAndProvider(ctx context.Context, id account.UserID, provider string) (account.Account, error)
}

type BaseUserStore struct {
	repo Repository
}

var _ (account.UserStore) = (*BaseUserStore)(nil)

func NewBaseUserStore(repo Repository) *BaseUserStore {
	return &BaseUserStore{
		repo: repo,
	}
}

func (s *BaseUserStore) GetUserByID(ctx context.Context, id account.UserID) (account.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *BaseUserStore) GetAccountsForUserID(ctx context.Context, id account.UserID) ([]account.Account, error) {
	return s.repo.GetAccountsByUserID(ctx, id)
}

func (s *BaseUserStore) GetUserByAccountID(ctx context.Context, id account.AccountID) (account.User, error) {
	return s.repo.GetUserByAccountID(ctx, id)
}

func (s *BaseUserStore) GetAccountByID(ctx context.Context, id account.AccountID) (account.Account, error) {
	return s.repo.GetAccountByID(ctx, id)
}

func (s *BaseUserStore) GetAccountByProvider(ctx context.Context, userID account.UserID, provider string) (account.Account, error) {
	if spec, ok := s.repo.(UserIDAndPRoviderGetter); ok {
		return spec.GetAccountByUserIDAndProvider(ctx, userID, provider)
	}
	accounts, err := s.repo.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return account.Account{}, err
	}
	idx := slices.IndexFunc(accounts, func(acc account.Account) bool { return acc.Provider == provider })
	if idx < 0 {
		return account.Account{}, account.ErrAccountNotFound
	}
	return accounts[idx], nil
}

func (s *BaseUserStore) GetOrCreateUserFromAccount(ctx context.Context, acc account.Account) (account.User, error) {
	oldAcc, err := s.repo.GetAccountByID(ctx, acc.AccountID)
	if err == nil {
		acc.UserID = oldAcc.UserID
		if err = s.repo.UpdateAccount(ctx, acc.AccountID, acc); err != nil {
			return account.User{}, err
		}

		return s.repo.GetUserByID(ctx, acc.UserID)
	} else if !errors.Is(err, account.ErrAccountNotFound) {
		return account.User{}, err
	}

	user, err := s.repo.CreateUser(ctx, account.User{
		Name:          acc.Name,
		Email:         acc.Email,
		EmailVerified: acc.EmailVerified,
		Username:      acc.PreferredUsername,
		Picture:       acc.Picture,
	})
	if err != nil {
		return account.User{}, err
	}
	acc.UserID = user.UserID

	if _, err = s.repo.CreateAccount(ctx, acc); err != nil {
		return account.User{}, err
	}

	return user, nil
}

func (s *BaseUserStore) UpdateUser(ctx context.Context, id account.UserID, user account.User) error {
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *BaseUserStore) UpdateAccount(ctx context.Context, id account.AccountID, acc account.Account) error {
	return s.repo.UpdateAccount(ctx, id, acc)
}

func (s *BaseUserStore) LinkAccount(ctx context.Context, id account.UserID, acc account.Account) error {
	var err error
	if _, err = s.repo.GetUserByID(ctx, id); err != nil {
		return err
	}

	acc.UserID = id

	_, err = s.repo.GetAccountByID(ctx, acc.AccountID)
	if err == nil {
		return s.repo.UpdateAccount(ctx, acc.AccountID, acc)
		// return account.ErrAccountAlreadyLinked
	} else if !errors.Is(err, account.ErrAccountNotFound) {
		return err
	}

	_, err = s.repo.CreateAccount(ctx, acc)

	return err
}

func (s *BaseUserStore) UnLinkAccount(ctx context.Context, id account.UserID, accID account.AccountID) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	oldAcc, err := s.repo.GetAccountByID(ctx, accID)
	if errors.Is(err, account.ErrAccountNotFound) {
		return account.ErrAccountNotLinked
	} else if err != nil {
		return err
	}

	if oldAcc.UserID != id {
		return account.ErrAccountNotLinked
	}

	accounts, err := s.GetAccountsForUserID(ctx, id)
	if err != nil {
		return err
	}

	if len(accounts) == 1 {
		return account.ErrLastAccount
	}

	if err := s.repo.DeleteAccountByID(ctx, accID); err != nil {
		return err
	}

	accounts = slices.Collect(its.Filter(slices.Values(accounts), func(acc account.Account) bool {
		return oldAcc.Provider != acc.Provider || oldAcc.ProviderID != acc.ProviderID
	}))

	newAcc := accounts[0]

	var userNeedsUpdate bool
	if user.Email == oldAcc.Email {
		user.Email = newAcc.Email
		userNeedsUpdate = true
	}
	if user.Name == oldAcc.Name {
		user.Name = newAcc.Name
		userNeedsUpdate = true
	}
	if user.Picture == oldAcc.Picture {
		user.Picture = newAcc.Picture
		userNeedsUpdate = true
	}
	if user.Username == oldAcc.PreferredUsername {
		user.Username = newAcc.PreferredUsername
		userNeedsUpdate = true
	}

	if userNeedsUpdate {
		return s.repo.UpdateUser(ctx, id, user)
	}

	return nil
}
