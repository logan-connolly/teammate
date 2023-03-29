package repository

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/access/domain/model"
	"github.com/logan-connolly/teammate/internal/entity"
)

var (
	ErrUserNotFound      = errors.New("the user was not found in the repository")
	ErrUserAlreadyExists = errors.New("user already exists in repository")
	ErrUserHasNoUpdates  = errors.New("failed to update user in the repository")
)

// UserRepository defines the interface for the user repository.
type UserRepository interface {
	Get(*entity.Person) (*model.User, error)
	Add(*model.User) error
	Update(*model.User) error
}
