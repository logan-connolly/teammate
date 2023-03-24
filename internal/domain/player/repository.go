package player

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrPlayerNotFound      = errors.New("the player was not found in the repository")
	ErrPlayerAlreadyExists = errors.New("failed to add the player to the repository")
	ErrPlayerHasNoUpdates  = errors.New("failed to update player in the repository")
)

// PlayerRepository defines the interface for the player repository.
type PlayerRepository interface {
	Get(uuid.UUID) (*Player, error)
	Add(*Player) error
	Update(*Player) error
}
