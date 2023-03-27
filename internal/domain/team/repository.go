package team

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTeamNotFound      = errors.New("the team was not found in the repository")
	ErrTeamAlreadyExists = errors.New("team already exists in repository")
	ErrTeamHasNoUpdates  = errors.New("failed to update team in the repository")
)

// TeamRepository defines the interface for the team repository.
type TeamRepository interface {
	Get(uuid.UUID) (*Team, error)
	Add(*Team) error
	Update(*Team) error
}
