package repository

import (
	"errors"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/model"
)

var (
	ErrPlayerNotFound      = errors.New("repository: the player was not found")
	ErrPlayerAlreadyExists = errors.New("repository: player already exists")
	ErrPlayerHasNoUpdates  = errors.New("repository: failed to update player")
)

// PlayerRepository defines the interface for the player repository.
type PlayerRepository interface {
	Get(*entity.Person) (*model.Player, error)
	GetTeams(*entity.Person) ([]*entity.Group, error)
	Add(*model.Player) error
	Update(*model.Player) error
}
