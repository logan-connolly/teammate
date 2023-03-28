package repository

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
)

var (
	ErrTeamNotFound      = errors.New("the team was not found in the repository")
	ErrTeamAlreadyExists = errors.New("team already exists in repository")
	ErrTeamHasNoUpdates  = errors.New("failed to update team in the repository")
)

// TeamRepository defines the interface for the team repository.
type TeamRepository interface {
	Get(*entity.Group) (*model.Team, error)
	Add(*model.Team) error
	Update(*model.Team) error
}
