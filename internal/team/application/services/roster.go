package services

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
	"github.com/logan-connolly/teammate/internal/team/domain/repository"
	"github.com/logan-connolly/teammate/internal/team/infrastructure/memory"
)

var ErrInvalidRosterConfig = errors.New("services: invalid roster configuration")

// ServiceConfigs defines the configurations to intialize the service with.
var ServiceConfigs = []RosterConfiguration{
	WithMemoryRepositories(),
}

// RosterConfiguration is a function that modifies the service.
type RosterConfiguration func(s *RosterService) error

// WithMemoryRepositories attaches in memory repostories to service.
func WithMemoryRepositories() RosterConfiguration {
	return func(s *RosterService) error {
		s.players = memory.NewMemoryPlayerRepository()
		s.teams = memory.NewMemoryTeamRepository()
		return nil
	}
}

// RosterService is a implementation of the RosterService.
type RosterService struct {
	players repository.PlayerRepository
	teams   repository.TeamRepository
}

// NewRosterService accepts configs and returns a new service.
func NewRosterService() (*RosterService, error) {
	s := &RosterService{}

	for _, cfg := range ServiceConfigs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// AddPlayer initializes a new player to the repository if valid.
func (s *RosterService) AddPlayer(player *entity.Person) error {
	p, err := model.NewPlayer(player)
	if err != nil {
		return err
	}

	return s.players.Add(p)
}

// AddTeam initializes a new team to the repository if valid.
func (s *RosterService) AddTeam(team *entity.Group) error {
	t, err := model.NewTeam(team)
	if err != nil {
		return err
	}

	return s.teams.Add(t)
}

// AssignPlayerToTeam assigns player to team's roster.
func (s *RosterService) AssignPlayerToTeam(team *entity.Group, player *entity.Person) error {
	t, err := s.teams.Get(team)
	if err != nil {
		return err
	}
	p, err := s.players.Get(player)
	if err != nil {
		return err
	}
	err = t.AssignPlayer(p)
	if err != nil {
		return err
	}
	err = p.AssignTeam(t)
	if err != nil {
		return err
	}

	s.teams.Update(t)
	s.players.Update(p)

	return nil
}

// UnassignPlayerToTeam unassigns player from team's roster.
func (s *RosterService) UnassignPlayerFromTeam(team *entity.Group, player *entity.Person) error {
	t, err := s.teams.Get(team)
	if err != nil {
		return err
	}
	p, err := s.players.Get(player)
	if err != nil {
		return err
	}
	err = t.UnassignPlayer(p)
	if err != nil {
		return err
	}
	err = p.UnassignTeam(t)
	if err != nil {
		return err
	}

	s.teams.Update(t)
	s.players.Update(p)

	return nil
}
