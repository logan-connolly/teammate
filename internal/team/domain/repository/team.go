package repository

import (
	"errors"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/model"
)

var (
	ErrTeamNotFound      = errors.New("repository: the team was not found")
	ErrTeamAlreadyExists = errors.New("repository: team already exists")
	ErrTeamHasNoUpdates  = errors.New("repository: failed to update team")
)

// TeamRepository defines the interface for the team repository.
type TeamRepository interface {
	Get(*entity.Group) (*model.Team, error)
	GetPlayers(*entity.Group) ([]*entity.Person, error)
	Add(*model.Team) error
	Update(*model.Team) error
}
