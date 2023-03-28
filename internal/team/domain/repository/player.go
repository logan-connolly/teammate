package repository

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
)

var (
	ErrPlayerNotFound      = errors.New("repository: the player was not found")
	ErrPlayerAlreadyExists = errors.New("repository: player already exists")
	ErrPlayerHasNoUpdates  = errors.New("repository: failed to update player")
)

// PlayerRepository defines the interface for the player repository.
type PlayerRepository interface {
	Get(*entity.Person) (*model.Player, error)
	Add(*model.Player) error
	Update(*model.Player) error
}
