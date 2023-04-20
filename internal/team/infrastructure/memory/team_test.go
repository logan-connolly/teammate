package memory

import (
	"testing"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/event"
	"git.sr.ht/~loges/teammate/internal/team/domain/model"
	"git.sr.ht/~loges/teammate/internal/team/domain/repository"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

var (
	exampleTeamUUID    = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherTeamUUID    = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	exampleTeamName    = "Syracuse"
	anotherTeamName    = "Notre Dame"
	teamCreated        = &event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName}
	anotherTeamCreated = &event.TeamCreated{ID: anotherTeamUUID, Name: anotherTeamName}
	playerAssigned     = &event.PlayerAssignedToTeam{ID: exampleTeamUUID, PlayerId: examplePlayerUUID, PlayerName: examplePlayerName}
)

func TestMemoryTeamRepository_Get(t *testing.T) {
	testCases := []struct {
		test        string
		group       *entity.Group
		expectedErr error
	}{
		{
			"No team found with this group",
			&entity.Group{ID: anotherTeamUUID, Name: anotherTeamName},
			repository.ErrTeamNotFound,
		}, {
			"Team found",
			&entity.Group{ID: exampleTeamUUID, Name: exampleTeamName},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			repo := NewMemoryTeamRepository()
			repo.teams[exampleTeamUUID] = []event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}

			_, err := repo.Get(tc.group)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestMemoryTeamRepository_GetPlayers(t *testing.T) {
	testCases := []struct {
		test        string
		teamId      uuid.UUID
		playerCount int
		events      []event.Event
		expectedErr error
	}{
		{"Team not found", anotherTeamUUID, 0, []event.Event{anotherTeamCreated}, repository.ErrTeamNotFound},
		{"No player is assigned", exampleTeamUUID, 0, []event.Event{teamCreated}, nil},
		{"One player is assigned", exampleTeamUUID, 1, []event.Event{teamCreated, playerAssigned}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			repo := NewMemoryTeamRepository()
			repo.teams[tc.teamId] = tc.events

			players, err := repo.GetPlayers(&entity.Group{ID: exampleTeamUUID, Name: exampleTeamName})

			is.Equal(err, tc.expectedErr)
			is.Equal(len(players), tc.playerCount)
		})
	}
}

func TestMemoryTeamRepository_Add(t *testing.T) {
	testCases := []struct {
		test        string
		id          uuid.UUID
		name        string
		expectedErr error
	}{
		{"Successfully add a team", anotherTeamUUID, anotherTeamName, nil},
		{"Team already exists error", exampleTeamUUID, exampleTeamName, repository.ErrTeamAlreadyExists},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			r := NewMemoryTeamRepository()
			team := model.NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: tc.id, Name: tc.name},
			})
			r.teams[exampleTeamUUID] = team.Events()

			err := r.Add(team)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestMemoryTeamRepository_Update(t *testing.T) {
	testCases := []struct {
		test        string
		register    bool
		deactivate  bool
		expectedErr error
	}{
		{"Update team", true, true, nil},
		{"Team has no changes", true, false, repository.ErrTeamHasNoUpdates},
		{"Team not found", false, true, repository.ErrTeamNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			r := NewMemoryTeamRepository()
			team := model.NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			})
			if tc.register {
				r.teams[exampleTeamUUID] = team.Events()
			}
			if tc.deactivate {
				team.Deactivate()
			}

			err := r.Update(team)

			is.Equal(err, tc.expectedErr)
		})
	}
}
