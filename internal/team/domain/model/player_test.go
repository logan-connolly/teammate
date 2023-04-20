package model

import (
	"testing"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/event"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

var (
	examplePlayerUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	examplePlayerName = "Logan"
	playerCreated     = &event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName}
	playerDeactivated = &event.PlayerDeactivated{ID: examplePlayerUUID}
	teamAssigned      = &event.TeamAssignedToPlayer{ID: examplePlayerUUID, TeamId: exampleTeamUUID, TeamName: exampleTeamName}
)

func TestPlayer_NewPlayer(t *testing.T) {
	testCases := []struct {
		test        string
		person      *entity.Person
		expectedErr error
	}{
		{"Empty name validation", &entity.Person{ID: examplePlayerUUID, Name: ""}, ErrInvalidPerson},
		{"Valid name", &entity.Person{ID: examplePlayerUUID, Name: examplePlayerName}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			_, err := NewPlayer(tc.person)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestPlayer_NewEvents(t *testing.T) {
	testCases := []struct {
		test                string
		events              []event.Event
		expectedIsActivated bool
	}{
		{"Player registered", []event.Event{playerCreated}, true},
		{"Player registered and deactivated", []event.Event{playerCreated, playerDeactivated}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			p := NewPlayerFromEvents(tc.events)
			is.Equal(p.GetID(), examplePlayerUUID)
			is.Equal(p.GetName(), examplePlayerName)
			is.Equal(p.IsActivated(), tc.expectedIsActivated)
		})
	}
}

func TestPlayer_Activate(t *testing.T) {
	testCases := []struct {
		test        string
		player      *Player
		expectedErr error
	}{
		{
			"Activate active player",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			ErrPlayerUpdateFailed,
		},
		{
			"Activate deactivated player",
			NewPlayerFromEvents([]event.Event{playerCreated, playerDeactivated}),
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.player.Activate()
			is.Equal(err, tc.expectedErr)
			is.True(tc.player.IsActivated())
		})
	}
}

func TestPlayer_Deactivate(t *testing.T) {
	testCases := []struct {
		test        string
		player      *Player
		expectedErr error
	}{
		{
			"Deactivate active player",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			nil,
		},
		{
			"Deactivate deactivated player",
			NewPlayerFromEvents([]event.Event{playerCreated, playerDeactivated}),
			ErrPlayerUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.player.Deactivate()
			is.Equal(tc.expectedErr, err)
			is.Equal(tc.player.IsActivated(), false)
		})
	}
}

func TestPlayer_AssignTeam(t *testing.T) {
	testCases := []struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}{
		{
			"Add team to player",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			nil,
		},
		{
			"Add team that is already assigned to player",
			NewPlayerFromEvents([]event.Event{playerCreated, teamAssigned}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			ErrPlayerUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.player.AssignTeam(tc.team)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestPlayer_UnassignTeam(t *testing.T) {
	type testCase struct{}
	testCases := []struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}{
		{
			"Unassign team to player",
			NewPlayerFromEvents([]event.Event{playerCreated, teamAssigned}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			nil,
		},
		{
			"Try to unassigned a team that is not assigned to player",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			ErrPlayerUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.player.UnassignTeam(tc.team)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestPlayer_GetTeams(t *testing.T) {
	testCases := []struct {
		test      string
		player    *Player
		team      *Team
		teamCount int
	}{
		{
			"No teams assigned",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			0,
		},
		{
			"Add player that is already assigned to team",
			NewPlayerFromEvents([]event.Event{playerCreated, teamAssigned}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			team := tc.player.GetTeams()
			is.Equal(len(team), tc.teamCount)
		})
	}
}

func TestPlayer_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		is := is.New(t)
		player, err := NewPlayer(&entity.Person{ID: examplePlayerUUID, Name: examplePlayerName})
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		player.Deactivate()
		player.Activate()

		is.Equal(len(player.Events()), 3)
	})
}

func TestPlayer_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		is := is.New(t)
		p := NewPlayerFromEvents([]event.Event{playerCreated, playerDeactivated})
		is.Equal(p.Version(), 2)
	})
}
