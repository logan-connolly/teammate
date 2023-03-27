package registration

import (
	"log"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/domain/player"
	"github.com/logan-connolly/teammate/internal/domain/player/memory"
)

// RegistrationConfiguration is a function that modifies the service.
type RegistrationConfiguration func(s *RegistrationService) error

// RegistrationService is a implementation of the RegistrationService.
type RegistrationService struct {
	players player.PlayerRepository
}

// NewRegistrationService accepts configs and returns a new service.
func NewRegistrationService(cfgs ...RegistrationConfiguration) (*RegistrationService, error) {
	s := &RegistrationService{}

	for _, cfg := range cfgs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// WithPlayerRepository applies a given player repository to the service.
func WithPlayerRepository(r player.PlayerRepository) RegistrationConfiguration {
	return func(s *RegistrationService) error {
		s.players = r
		return nil
	}
}

// WithMemoryPlayerRepository applies a memory player repository to the service.
func WithMemoryPlayerRepository() RegistrationConfiguration {
	r := memory.NewMemoryPlayerRepository()
	return WithPlayerRepository(r)
}

// RegisterPlayer will register a player if input is valid.
func (s *RegistrationService) RegisterPlayer(name string) (uuid.UUID, error) {
	p, err := player.NewPlayer(name)
	if err != nil {
		return uuid.Nil, err
	}

	err = s.players.Add(p)
	if err != nil {
		return uuid.Nil, err
	}

	log.Printf("Player: %s has been registered.", p.GetID())
	return p.GetID(), nil
}
