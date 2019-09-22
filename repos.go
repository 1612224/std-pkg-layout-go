package app

import "errors"

var (
	// ErrNotFound is an implementation-independent error
	// that should be return by any repo implementation
	// when a record is not found
	ErrNotFound = errors.New("app: the requested resource is not found")
)

// UserRepo is an interface for interact with users in database
type UserRepo interface {
	ByEmail(email string) (*User, error)
	ByToken(token int) (*User, error)
	UpdateToken(userID int, newToken int) error
}

// ItemRepo is an interface for interact with items in database
type ItemRepo interface {
	ByUser(userID int) ([]Item, error)
	Create(item *Item) error
}
