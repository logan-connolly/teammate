package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
	"github.com/logan-connolly/teammate/internal/team/domain/repository"
)

var (
	exampleGroup  = &entity.Group{ID: uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002"), Name: "Tigers"}
	examplePerson = &entity.Person{ID: uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"), Name: "Matt"}
	anotherGroup  = &entity.Group{ID: uuid.MustParse("adbe93f8-c952-11ed-afa1-0242ac120002"), Name: "Bears"}
	anotherPerson = &entity.Person{ID: uuid.MustParse("d38ad10b-58cc-0372-8567-0e02b2c3d479"), Name: "Jackie"}
)

func TestNewRosterService(t *testing.T) {
	type testCase struct {
		test        string
		testConfig  func() RosterConfiguration
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "With in-memory player repository",
			testConfig:  WithMemoryPlayerRepository,
			expectedErr: nil,
		},
		{
			test:        "With in-memory team repository",
			testConfig:  WithMemoryTeamRepository,
			expectedErr: nil,
		},
		{
			test: "With bad config",
			testConfig: func() RosterConfiguration {
				return func(s *RosterService) error {
					return ErrInvalidRosterConfig
				}
			},
			expectedErr: ErrInvalidRosterConfig,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := NewRosterService(tc.testConfig())
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected %v, but got %v.", tc.expectedErr, err)
			}
		})
	}
}

func TestRosterService_AssignPlayerToTeam(t *testing.T) {
	type testCase struct {
		test                string
		group               *entity.Group
		person              *entity.Person
		alreadyAssignPlayer bool
		alreadyAssignTeam   bool
		expectedErr         error
	}

	testCases := []testCase{
		{
			test:        "Team not found",
			group:       anotherGroup,
			person:      examplePerson,
			expectedErr: repository.ErrTeamNotFound,
		},
		{
			test:        "Player not found",
			group:       exampleGroup,
			person:      anotherPerson,
			expectedErr: repository.ErrPlayerNotFound,
		},
		{
			test:              "Team already assigned to player",
			group:             exampleGroup,
			person:            examplePerson,
			alreadyAssignTeam: true,
			expectedErr:       model.ErrTeamAlreadyAssigned,
		},
		{
			test:                "Player already assigned to Team",
			group:               exampleGroup,
			person:              examplePerson,
			alreadyAssignPlayer: true,
			expectedErr:         model.ErrPlayerAlreadyAssigned,
		},
		{
			test:        "Player assigned to team",
			group:       exampleGroup,
			person:      examplePerson,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
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

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected %v, but got %v.", tc.expectedErr, err)
			}
		})
	}
}
