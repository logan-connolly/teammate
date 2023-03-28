package services

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/team/domain/repository"
	"github.com/logan-connolly/teammate/internal/team/domain/storage/memory"
)

var ErrInvalidRosterConfig = errors.New("services: invalid roster configuration")

// RosterConfiguration is a function that modifies the service.
type RosterConfiguration func(s *RosterService) error

// RosterService is a implementation of the RosterService.
type RosterService struct {
	players repository.PlayerRepository
	teams   repository.TeamRepository
}

// NewRosterService accepts configs and returns a new service.
func NewRosterService(cfgs ...RosterConfiguration) (*RosterService, error) {
	s := &RosterService{}

	for _, cfg := range cfgs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// WithPlayerRepository applies a given player repository to the service.
func WithPlayerRepository(r repository.PlayerRepository) RosterConfiguration {
	return func(s *RosterService) error {
		s.players = r
		return nil
	}
}

// WithMemoryPlayerRepository applies a memory player repository to the service.
func WithMemoryPlayerRepository() RosterConfiguration {
	r := memory.NewMemoryPlayerRepository()
	return WithPlayerRepository(r)
}

// WithTeamRepository applies a given team repository to the service.
func WithTeamRepository(r repository.TeamRepository) RosterConfiguration {
	return func(s *RosterService) error {
		s.teams = r
		return nil
	}
}

// WithMemoryTeamRepository applies a memory team repository to the service.
func WithMemoryTeamRepository() RosterConfiguration {
	r := memory.NewMemoryTeamRepository()
	return WithTeamRepository(r)
}
