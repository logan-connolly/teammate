package services

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/team/domain/repository"
	"github.com/logan-connolly/teammate/internal/team/domain/storage/memory"
)

var ErrInvalidRegistrationConfig = errors.New("Invalid registration configuration")

// RegistrationConfiguration is a function that modifies the service.
type RegistrationConfiguration func(s *RegistrationService) error

// RegistrationService is a implementation of the RegistrationService.
type RegistrationService struct {
	players repository.PlayerRepository
	teams   repository.TeamRepository
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
func WithPlayerRepository(r repository.PlayerRepository) RegistrationConfiguration {
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

// WithTeamRepository applies a given team repository to the service.
func WithTeamRepository(r repository.TeamRepository) RegistrationConfiguration {
	return func(s *RegistrationService) error {
		s.teams = r
		return nil
	}
}

// WithMemoryTeamRepository applies a memory team repository to the service.
func WithMemoryTeamRepository() RegistrationConfiguration {
	r := memory.NewMemoryTeamRepository()
	return WithTeamRepository(r)
}
