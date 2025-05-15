package account

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrAccountNotFound      = errors.New("account not found")
	ErrNoAccountFound       = errors.New("no account not found")
	ErrAccountAlreadyLinked = errors.New("account already linked")
	ErrAccountNotLinked     = errors.New("account is not linked")
	ErrLastAccount          = errors.New("cant unlink last account")
)

type UserStore interface {
	GetUserByID(ctx context.Context, id UserID) (User, error)
	GetUserByAccountID(ctx context.Context, id AccountID) (User, error)
	GetAccountsForUserID(ctx context.Context, id UserID) ([]Account, error)
	GetAccountByID(ctx context.Context, accID AccountID) (Account, error)
	GetAccountByProvider(ctx context.Context, id UserID, provider string) (Account, error)

	GetOrCreateUserFromAccount(ctx context.Context, acc Account) (User, error)
	UpdateUser(ctx context.Context, id UserID, user User) error
	UpdateAccount(ctx context.Context, id AccountID, acc Account) error
	LinkAccount(ctx context.Context, id UserID, acc Account) error
	UnLinkAccount(ctx context.Context, id UserID, accID AccountID) error
}
