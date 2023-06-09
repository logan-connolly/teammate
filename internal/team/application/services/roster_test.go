package services

import (
	"testing"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/model"
	"git.sr.ht/~loges/teammate/internal/team/domain/repository"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

var (
	exampleGroup  = &entity.Group{ID: uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002"), Name: "Tigers"}
	examplePerson = &entity.Person{ID: uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"), Name: "Matt"}
	anotherGroup  = &entity.Group{ID: uuid.MustParse("adbe93f8-c952-11ed-afa1-0242ac120002"), Name: "Bears"}
	anotherPerson = &entity.Person{ID: uuid.MustParse("d38ad10b-58cc-0372-8567-0e02b2c3d479"), Name: "Jackie"}
)

func TestNewRosterService(t *testing.T) {
	t.Run("Create service with defaults", func(t *testing.T) {
		is := is.New(t)
		_, err := NewRosterService()
		is.NoErr(err)
	})

	t.Run("Create service with bad config", func(t *testing.T) {
		is := is.New(t)
		withInvalidConfig := func() RosterConfiguration {
			return func(s *RosterService) error {
				return ErrInvalidRosterConfig
			}
		}
		originalConfigs := RosterConfigs
		RosterConfigs = []RosterConfiguration{withInvalidConfig()}

		_, err := NewRosterService()

		is.Equal(err, ErrInvalidRosterConfig)
		// clean up configs
		RosterConfigs = originalConfigs
	})
}

func TestRosterService_AddPlayer(t *testing.T) {
	testCases := []struct {
		test        string
		person      *entity.Person
		expectedErr error
	}{
		{"Player name required", &entity.Person{ID: anotherPerson.ID, Name: ""}, model.ErrInvalidPerson},
		{"Player already exists", examplePerson, repository.ErrPlayerAlreadyExists},
		{"Player successfully added", anotherPerson, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			s, _ := NewRosterService()
			_ = s.AddPlayer(examplePerson)

			err := s.AddPlayer(tc.person)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestRosterService_AddTeam(t *testing.T) {
	testCases := []struct {
		test        string
		group       *entity.Group
		expectedErr error
	}{
		{"Team name required", &entity.Group{ID: anotherGroup.ID, Name: ""}, model.ErrInvalidGroup},
		{"Team already exists", exampleGroup, repository.ErrTeamAlreadyExists},
		{"Team successfully added", anotherGroup, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			s, _ := NewRosterService()
			_ = s.AddTeam(exampleGroup)

			err := s.AddTeam(tc.group)

			is.Equal(err, tc.expectedErr)
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
		{"Team already assigned to player", exampleGroup, examplePerson, false, true, model.ErrPlayerUpdateFailed},
		{"Player already assigned to Team", exampleGroup, examplePerson, true, false, model.ErrTeamUpdateFailed},
		{"Player assigned to team", exampleGroup, examplePerson, false, false, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)

			// initialize roster service
			s, err := NewRosterService()

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
		{"Team not assigned to player", exampleGroup, examplePerson, true, false, model.ErrPlayerUpdateFailed},
		{"Player not assigned to team", exampleGroup, examplePerson, false, true, model.ErrTeamUpdateFailed},
		{"Player unassigned from team", exampleGroup, examplePerson, true, true, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)

			// initialize roster service
			s, err := NewRosterService()

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
