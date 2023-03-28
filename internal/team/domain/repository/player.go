package repository

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
)

var (
	ErrPlayerNotFound      = errors.New("the player was not found in the repository")
	ErrPlayerAlreadyExists = errors.New("player already exists in repository")
	ErrPlayerHasNoUpdates  = errors.New("failed to update player in the repository")
)

// PlayerRepository defines the interface for the player repository.
type PlayerRepository interface {
	Get(*entity.Person) (*model.Player, error)
	Add(*model.Player) error
	Update(*model.Player) error
}
