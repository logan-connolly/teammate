package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
)

var (
	examplePlayerUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	examplePlayerName = "Logan"
)

func TestPlayer_NewPlayer(t *testing.T) {
	type testCase struct {
		test        string
		person      *entity.Person
		expectedErr error
	}
	testCases := []testCase{
		{
			test:        "Empty name validation",
			person:      &entity.Person{ID: examplePlayerUUID, Name: ""},
			expectedErr: ErrInvalidPerson,
		},
		{
			test:        "Valid name",
			person:      &entity.Person{ID: examplePlayerUUID, Name: examplePlayerName},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := NewPlayer(tc.person)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestPlayer_NewEvents(t *testing.T) {
	type testCase struct {
		test                string
		events              []event.Event
		expectedIsActivated bool
	}
	testCases := []testCase{
		{
			test: "Player registered",
			events: []event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			},
			expectedIsActivated: true,
		},
		{
			test: "Player registered and deactivated",
			events: []event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
				&event.PlayerDeactivated{ID: examplePlayerUUID},
			},
			expectedIsActivated: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			p := NewPlayerFromEvents(tc.events)
			if p.GetID() != examplePlayerUUID {
				t.Errorf("Expected %v, got %v", examplePlayerUUID, p.GetID())
			}
			if p.GetName() != examplePlayerName {
				t.Errorf("Expected %v, got %v", examplePlayerName, p.GetName())
			}
			if p.IsActivated() != tc.expectedIsActivated {
				t.Errorf("Expected %v, got %v", tc.expectedIsActivated, p.IsActivated())
			}
		})
	}
}

func TestPlayer_Activate(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Activate active player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			expectedErr: ErrPlayerAlreadyActivated,
		},
		{
			test: "Activate deactivated player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
				&event.PlayerDeactivated{ID: examplePlayerUUID},
			}),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.player.Activate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if !tc.player.IsActivated() {
				t.Fatal("player should always be activated in these cases.")
			}
		})
	}
}

func TestPlayer_Deactivate(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Deactivate active player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			expectedErr: nil,
		},
		{
			test: "Deactivate deactivated player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
				&event.PlayerDeactivated{ID: examplePlayerUUID},
			}),
			expectedErr: ErrPlayerAlreadyDeactivated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.player.Deactivate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if tc.player.IsActivated() {
				t.Fatal("player should always be deactivated in these cases.")
			}
		})
	}
}

func TestPlayer_AssignTeam(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Add team to player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: nil,
		},
		{
			test: "Add team that is already assigned to player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
				&event.TeamAssignedToPlayer{ID: examplePlayerUUID, TeamId: exampleTeamUUID, TeamName: exampleTeamName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: ErrTeamAlreadyAssigned,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.player.AssignTeam(tc.team)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestPlayer_UnassignTeam(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Unassign team to player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
				&event.TeamAssignedToPlayer{ID: examplePlayerUUID, TeamId: exampleTeamUUID, TeamName: exampleTeamName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: nil,
		},
		{
			test: "Try to unassigned a team that is not assigned to player",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: ErrTeamNotAssignedToPlayer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.player.UnassignTeam(tc.team)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestPlayer_GetTeams(t *testing.T) {
	type testCase struct {
		test      string
		player    *Player
		team      *Team
		teamCount int
	}
	testCases := []testCase{
		{
			test: "No teams assigned",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			teamCount: 0,
		},
		{
			test: "Add player that is already assigned to team",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
				&event.TeamAssignedToPlayer{ID: examplePlayerUUID, TeamId: exampleTeamUUID, TeamName: exampleTeamName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			teamCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			team := tc.player.GetTeams()
			if len(team) != tc.teamCount {
				t.Errorf("Expected %d, got %d", tc.teamCount, len(team))
			}
		})
	}
}

func TestPlayer_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		player, err := NewPlayer(&entity.Person{ID: examplePlayerUUID, Name: examplePlayerName})
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		player.Deactivate()
		player.Activate()

		want := 3
		got := len(player.Events())

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

func TestPlayer_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		p := NewPlayerFromEvents([]event.Event{
			&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			&event.PlayerDeactivated{ID: examplePlayerUUID},
			&event.PlayerActivated{ID: examplePlayerUUID},
		})

		want := 3
		got := p.Version()

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
