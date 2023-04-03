package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
	"github.com/logan-connolly/teammate/internal/team/domain/repository"
	"github.com/matryer/is"
)

var (
	exampleGroup  = &entity.Group{ID: uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002"), Name: "Tigers"}
	examplePerson = &entity.Person{ID: uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"), Name: "Matt"}
	anotherGroup  = &entity.Group{ID: uuid.MustParse("adbe93f8-c952-11ed-afa1-0242ac120002"), Name: "Bears"}
	anotherPerson = &entity.Person{ID: uuid.MustParse("d38ad10b-58cc-0372-8567-0e02b2c3d479"), Name: "Jackie"}
)

func withBadConfig() RosterConfiguration {
	return func(s *RosterService) error {
		return ErrInvalidRosterConfig
	}
}

func TestNewRosterService(t *testing.T) {
	testCases := []struct {
		test        string
		testConfig  func() RosterConfiguration
		expectedErr error
	}{
		{"With in-memory player repository", WithMemoryPlayerRepository, nil},
		{"With in-memory team repository", WithMemoryTeamRepository, nil},
		{"With bad config", withBadConfig, ErrInvalidRosterConfig},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			_, err := NewRosterService(tc.testConfig())
			is.Equal(tc.expectedErr, err)
		})
	}
}

func TestRosterService_AssignPlayerToTeam(t *testing.T) {
	testCases := []struct {
		test                string
		group               *entity.Group
		person              *entity.Person
		alreadyAssignPlayer bool
		alreadyAssignTeam   bool
		expectedErr         error
	}{
		{"Team not found", anotherGroup, examplePerson, false, false, repository.ErrTeamNotFound},
		{"Player not found", exampleGroup, anotherPerson, false, false, repository.ErrPlayerNotFound},
		{"Team already assigned to player", exampleGroup, examplePerson, false, true, model.ErrTeamAlreadyAssigned},
		{"Player already assigned to Team", exampleGroup, examplePerson, true, false, model.ErrPlayerAlreadyAssigned},
		{"Player assigned to team", exampleGroup, examplePerson, false, false, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)

			// initialize roster service
			s, err := NewRosterService(WithMemoryPlayerRepository(), WithMemoryTeamRepository())

			// instantiate player and team aggregate
			player, _ := model.NewPlayer(examplePerson)
			team, _ := model.NewTeam(exampleGroup)

			// already assign to team or player aggregate
			if tc.alreadyAssignTeam {
				player.AssignTeam(team)
			}
			if tc.alreadyAssignPlayer {
				team.AssignPlayer(player)
			}

			// store aggregates in repository
			s.players.Add(player)
			s.teams.Add(team)

			err = s.AssignPlayerToTeam(tc.group, tc.person)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestRosterService_UnassignPlayerToTeam(t *testing.T) {
	testCases := []struct {
		test         string
		group        *entity.Group
		person       *entity.Person
		assignPlayer bool
		assignTeam   bool
		expectedErr  error
	}{
		{"Team not found", anotherGroup, examplePerson, false, false, repository.ErrTeamNotFound},
		{"Player not found", exampleGroup, anotherPerson, false, false, repository.ErrPlayerNotFound},
		{"Team not assigned to player", exampleGroup, examplePerson, true, false, model.ErrTeamNotAssignedToPlayer},
		{"Player not assigned to team", exampleGroup, examplePerson, false, true, model.ErrPlayerNotAssignedToTeam},
		{"Player unassigned from team", exampleGroup, examplePerson, true, true, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)

			// initialize roster service
			s, err := NewRosterService(WithMemoryPlayerRepository(), WithMemoryTeamRepository())

			// instantiate player and team aggregate
			player, _ := model.NewPlayer(examplePerson)
			team, _ := model.NewTeam(exampleGroup)

			// already assign to team or player aggregate
			if tc.assignTeam {
				player.AssignTeam(team)
			}
			if tc.assignPlayer {
				team.AssignPlayer(player)
			}

			// store aggregates in repository
			s.players.Add(player)
			s.teams.Add(team)

			err = s.UnassignPlayerFromTeam(tc.group, tc.person)

			is.Equal(err, tc.expectedErr)
		})
	}
}
