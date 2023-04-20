package repository

import (
	"errors"

	"git.sr.ht/~loges/teammate/internal/access/domain/model"
)

var (
	ErrUserNotFound      = errors.New("repository: the user was not found")
	ErrUserAlreadyExists = errors.New("repository: user already exists")
	ErrUserHasNoUpdates  = errors.New("repository: failed to update user")
)

// UserRepository defines the interface for the user repository.
type UserRepository interface {
	GetByEmail(string) (*model.User, error)
	Add(*model.User) error
	Update(*model.User) error
}
