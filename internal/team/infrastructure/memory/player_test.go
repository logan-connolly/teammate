package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
	"github.com/logan-connolly/teammate/internal/team/domain/repository"
	"github.com/matryer/is"
)

var (
	examplePlayerUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherPlayerUUID = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	examplePlayerName = "Logan"
	anotherPlayerName = "Emily"
)

func TestMemoryPlayerRepository_Get(t *testing.T) {
	type testCase struct {
		test        string
		person      *entity.Person
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "No player with this person",
			person:      &entity.Person{ID: anotherPlayerUUID, Name: anotherTeamName},
			expectedErr: repository.ErrPlayerNotFound,
		},
		{
			test:        "Player found",
			person:      &entity.Person{ID: examplePlayerUUID, Name: exampleTeamName},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			repo := NewMemoryPlayerRepository()
			repo.players[examplePlayerUUID] = []event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}

			_, err := repo.Get(tc.person)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestMemoryPlayerRepository_GetTeams(t *testing.T) {
	type testCase struct {
		test        string
		playerId    uuid.UUID
		teamCount   int
		events      []event.Event
		expectedErr error
	}

	testCases := []testCase{
		{
			test:      "Player not found",
			playerId:  anotherPlayerUUID,
			teamCount: 0,
			events: []event.Event{
				&event.PlayerCreated{ID: anotherPlayerUUID, Name: anotherPlayerName},
			},
			expectedErr: repository.ErrPlayerNotFound,
		},
		{
			test:      "No team is assigned",
			playerId:  examplePlayerUUID,
			teamCount: 0,
			events: []event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			},
			expectedErr: nil,
		},
		{
			test:      "One team is assigned",
			playerId:  examplePlayerUUID,
			teamCount: 1,
			events: []event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
				&event.TeamAssignedToPlayer{ID: examplePlayerUUID, TeamId: exampleTeamUUID, TeamName: exampleTeamName},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			repo := NewMemoryPlayerRepository()
			repo.players[tc.playerId] = tc.events

			teams, err := repo.GetTeams(&entity.Person{ID: examplePlayerUUID, Name: examplePlayerName})

			is.Equal(err, tc.expectedErr)
			is.Equal(len(teams), tc.teamCount)
		})
	}
}

func TestMemoryPlayerRepository_Add(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		name        string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Successfully add a player",
			id:          anotherPlayerUUID,
			name:        anotherPlayerName,
			expectedErr: nil,
		},
		{
			test:        "Player already exists error",
			id:          examplePlayerUUID,
			name:        examplePlayerName,
			expectedErr: repository.ErrPlayerAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			r := NewMemoryPlayerRepository()
			p := model.NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: tc.id, Name: tc.name},
			})
			r.players[examplePlayerUUID] = p.Events()

			err := r.Add(p)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestMemoryPlayerRepository_Update(t *testing.T) {
	type testCase struct {
		test        string
		register    bool
		deactivate  bool
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Update player",
			register:    true,
			deactivate:  true,
			expectedErr: nil,
		},
		{
			test:        "Player has no changes",
			register:    true,
			deactivate:  false,
			expectedErr: repository.ErrPlayerHasNoUpdates,
		},
		{
			test:        "Player not found",
			register:    false,
			deactivate:  true,
			expectedErr: repository.ErrPlayerNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			r := NewMemoryPlayerRepository()
			p := model.NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: exampleTeamName},
			})
			if tc.register {
				r.players[examplePlayerUUID] = p.Events()
			}
			if tc.deactivate {
				p.Deactivate()
			}

			err := r.Update(p)

			is.Equal(err, tc.expectedErr)
		})
	}
}
