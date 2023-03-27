package registration

import (
	"github.com/logan-connolly/teammate/internal/domain/player"
	playerMemory "github.com/logan-connolly/teammate/internal/domain/player/memory"
	"github.com/logan-connolly/teammate/internal/domain/team"
	teamMemory "github.com/logan-connolly/teammate/internal/domain/team/memory"
)

// RegistrationConfiguration is a function that modifies the service.
type RegistrationConfiguration func(s *RegistrationService) error

// RegistrationService is a implementation of the RegistrationService.
type RegistrationService struct {
	players player.PlayerRepository
	teams   team.TeamRepository
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
	r := playerMemory.NewMemoryPlayerRepository()
	return WithPlayerRepository(r)
}

// WithTeamRepository applies a given team repository to the service.
func WithTeamRepository(r team.TeamRepository) RegistrationConfiguration {
	return func(s *RegistrationService) error {
		s.teams = r
		return nil
	}
}

// WithMemoryTeamRepository applies a memory team repository to the service.
func WithMemoryTeamRepository() RegistrationConfiguration {
	r := teamMemory.NewMemoryTeamRepository()
	return WithTeamRepository(r)
}
