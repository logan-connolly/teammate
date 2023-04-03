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
	examplePlayerUUID    = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherPlayerUUID    = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	examplePlayerName    = "Logan"
	anotherPlayerName    = "Emily"
	playerCreated        = &event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName}
	anotherPlayerCreated = &event.PlayerCreated{ID: anotherPlayerUUID, Name: anotherPlayerName}
	teamAssigned         = &event.TeamAssignedToPlayer{ID: examplePlayerUUID, TeamId: exampleTeamUUID, TeamName: exampleTeamName}
)

func TestMemoryPlayerRepository_Get(t *testing.T) {
	testCases := []struct {
		test        string
		person      *entity.Person
		expectedErr error
	}{
		{
			"No player with this person",
			&entity.Person{ID: anotherPlayerUUID, Name: anotherTeamName},
			repository.ErrPlayerNotFound,
		},
		{
			"Player found",
			&entity.Person{ID: examplePlayerUUID, Name: exampleTeamName},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			repo := NewMemoryPlayerRepository()
			repo.players[examplePlayerUUID] = []event.Event{playerCreated}

			_, err := repo.Get(tc.person)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestMemoryPlayerRepository_GetTeams(t *testing.T) {
	testCases := []struct {
		test        string
		playerId    uuid.UUID
		events      []event.Event
		teamCount   int
		expectedErr error
	}{
		{
			"Player not found",
			anotherPlayerUUID,
			[]event.Event{anotherPlayerCreated},
			0,
			repository.ErrPlayerNotFound,
		},
		{
			"No team is assigned",
			examplePlayerUUID,
			[]event.Event{playerCreated},
			0,
			nil,
		},
		{
			"One team is assigned",
			examplePlayerUUID,
			[]event.Event{playerCreated, teamAssigned},
			1,
			nil,
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
	testCases := []struct {
		test        string
		id          uuid.UUID
		name        string
		expectedErr error
	}{
		{"Successfully add a player", anotherPlayerUUID, anotherPlayerName, nil},
		{"Player already exists error", examplePlayerUUID, examplePlayerName, repository.ErrPlayerAlreadyExists},
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
	testCases := []struct {
		test        string
		register    bool
		deactivate  bool
		expectedErr error
	}{
		{"Update player", true, true, nil},
		{"Player has no changes", true, false, repository.ErrPlayerHasNoUpdates},
		{"Player not found", false, true, repository.ErrPlayerNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			r := NewMemoryPlayerRepository()
			p := model.NewPlayerFromEvents([]event.Event{playerCreated})
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
