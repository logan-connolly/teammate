package team

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/entity"
)

var (
	ErrTeamNotFound      = errors.New("the team was not found in the repository")
	ErrTeamAlreadyExists = errors.New("team already exists in repository")
	ErrTeamHasNoUpdates  = errors.New("failed to update team in the repository")
)

// TeamRepository defines the interface for the team repository.
type TeamRepository interface {
	Get(*entity.Group) (*Team, error)
	Add(*Team) error
	Update(*Team) error
}
