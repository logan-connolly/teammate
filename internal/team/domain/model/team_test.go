package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
	"github.com/matryer/is"
)

var (
	exampleTeamUUID = uuid.MustParse("d35e93f8-c952-11ed-afa1-0242ac120002")
	exampleTeamName = "Virginia"
	teamCreated     = &event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName}
	teamDeactivated = &event.TeamDeactivated{ID: exampleTeamUUID}
	playerAssigned  = &event.PlayerAssignedToTeam{ID: exampleTeamUUID, PlayerId: examplePlayerUUID, PlayerName: examplePlayerName}
)

func TestTeam_NewTeam(t *testing.T) {
	testCases := []struct {
		test        string
		group       *entity.Group
		expectedErr error
	}{
		{"Empty name validation", &entity.Group{ID: exampleTeamUUID, Name: ""}, ErrInvalidGroup},
		{"Valid name", &entity.Group{ID: exampleTeamUUID, Name: exampleTeamName}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			_, err := NewTeam(tc.group)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestTeam_NewEvents(t *testing.T) {
	testCases := []struct {
		test                string
		events              []event.Event
		expectedIsActivated bool
	}{
		{"Team registered", []event.Event{teamCreated}, true},
		{"Team registered and deactivated", []event.Event{teamCreated, teamDeactivated}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			team := NewTeamFromEvents(tc.events)
			is.Equal(team.GetID(), exampleTeamUUID)
			is.Equal(team.GetName(), exampleTeamName)
			is.Equal(team.IsActivated(), tc.expectedIsActivated)
		})
	}
}

func TestTeam_Activate(t *testing.T) {
	testCases := []struct {
		test        string
		team        *Team
		expectedErr error
	}{
		{
			"Activate active team",
			NewTeamFromEvents([]event.Event{teamCreated}),
			ErrTeamUpdateFailed,
		},
		{
			"Activate deactivated team",
			NewTeamFromEvents([]event.Event{teamCreated, teamDeactivated}),
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.team.Activate()
			is.Equal(err, tc.expectedErr)
			is.True(tc.team.IsActivated())
		})
	}
}

func TestTeam_Deactivate(t *testing.T) {
	testCases := []struct {
		test        string
		team        *Team
		expectedErr error
	}{
		{
			"Deactivate active team",
			NewTeamFromEvents([]event.Event{teamCreated}),
			nil,
		},
		{
			"Deactivate deactivated team",
			NewTeamFromEvents([]event.Event{teamCreated, teamDeactivated}),
			ErrTeamUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.team.Deactivate()
			is.Equal(err, tc.expectedErr)
			is.Equal(tc.team.IsActivated(), false)
		})
	}
}

func TestTeam_AssignPlayer(t *testing.T) {
	testCases := []struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}{
		{
			"Add player to team",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			nil,
		},
		{
			"Add player that is already assigned to team",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated, playerAssigned}),
			ErrTeamUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.team.AssignPlayer(tc.player)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestTeam_Unassignplayer(t *testing.T) {
	testCases := []struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}{
		{
			"Unassign player to team",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated, playerAssigned}),
			nil,
		},
		{
			"Try to unassigned a player that is not assigned to team",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			ErrTeamUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.team.UnassignPlayer(tc.player)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestTeam_GetPlayers(t *testing.T) {
	testCases := []struct {
		test        string
		player      *Player
		team        *Team
		playerCount int
	}{
		{
			"No players assigned",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated}),
			0,
		},
		{
			"Add player that is already assigned to team",
			NewPlayerFromEvents([]event.Event{playerCreated}),
			NewTeamFromEvents([]event.Event{teamCreated, playerAssigned}),
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			is.Equal(len(tc.team.GetPlayers()), tc.playerCount)
		})
	}
}

func TestTeam_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		is := is.New(t)
		team, err := NewTeam(&entity.Group{ID: exampleTeamUUID, Name: exampleTeamName})
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		team.Deactivate()
		team.Activate()

		is.Equal(len(team.Events()), 3)
	})
}

func TestTeam_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		is := is.New(t)
		team := NewTeamFromEvents([]event.Event{teamCreated, teamDeactivated})
		is.Equal(team.Version(), 2)
	})
}
