package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
)

var (
	exampleTeamUUID = uuid.MustParse("d35e93f8-c952-11ed-afa1-0242ac120002")
	exampleTeamName = "Virginia"
)

func TestTeam_NewTeam(t *testing.T) {
	type testCase struct {
		test        string
		group       *entity.Group
		expectedErr error
	}
	testCases := []testCase{
		{
			test:        "Empty name validation",
			group:       &entity.Group{ID: exampleTeamUUID, Name: ""},
			expectedErr: ErrInvalidGroup,
		},
		{
			test:        "Valid name",
			group:       &entity.Group{ID: exampleTeamUUID, Name: exampleTeamName},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := NewTeam(tc.group)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestTeam_NewEvents(t *testing.T) {
	type testCase struct {
		test                string
		events              []event.Event
		expectedIsActivated bool
	}
	testCases := []testCase{
		{
			test: "Team registered",
			events: []event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			},
			expectedIsActivated: true,
		},
		{
			test: "Team registered and deactivated",
			events: []event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
				&event.TeamDeactivated{ID: exampleTeamUUID},
			},
			expectedIsActivated: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			team := NewTeamFromEvents(tc.events)
			if team.GetID() != exampleTeamUUID {
				t.Errorf("Expected %v, got %v", exampleTeamUUID, team.GetID())
			}
			if team.GetName() != exampleTeamName {
				t.Errorf("Expected %v, got %v", exampleTeamName, team.GetName())
			}
			if team.IsActivated() != tc.expectedIsActivated {
				t.Errorf("Expected %v, got %v", tc.expectedIsActivated, team.IsActivated())
			}
		})
	}
}

func TestTeam_Activate(t *testing.T) {
	type testCase struct {
		test        string
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Activate active team",
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: ErrTeamAlreadyActivated,
		},
		{
			test: "Activate deactivated team",
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
				&event.TeamDeactivated{ID: exampleTeamUUID},
			}),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.team.Activate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if !tc.team.IsActivated() {
				t.Fatal("team should always be activated in these cases.")
			}
		})
	}
}

func TestTeam_Deactivate(t *testing.T) {
	type testCase struct {
		test        string
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Deactivate active team",
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: nil,
		},
		{
			test: "Deactivate deactivated team",
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
				&event.TeamDeactivated{ID: exampleTeamUUID},
			}),
			expectedErr: ErrTeamAlreadyDeactivated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.team.Deactivate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if tc.team.IsActivated() {
				t.Fatal("team should always be deactivated in these cases.")
			}
		})
	}
}

func TestTeam_AssignPlayer(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Add player to team",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: nil,
		},
		{
			test: "Add player that is already assigned to team",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
				&event.PlayerAssignedToTeam{ID: exampleTeamUUID, PlayerId: examplePlayerUUID, PlayerName: examplePlayerName},
			}),
			expectedErr: ErrPlayerAlreadyAssigned,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.team.AssignPlayer(tc.player)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestTeam_GetPlayers(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		team        *Team
		playerCount int
	}
	testCases := []testCase{
		{
			test: "No players assigned",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			playerCount: 0,
		},
		{
			test: "Add player that is already assigned to team",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
				&event.PlayerAssignedToTeam{ID: exampleTeamUUID, PlayerId: examplePlayerUUID, PlayerName: examplePlayerName},
			}),
			playerCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			players := tc.team.GetPlayers()
			if len(players) != tc.playerCount {
				t.Errorf("Expected %d, got %d", tc.playerCount, len(players))
			}
		})
	}
}

func TestTeam_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		team, err := NewTeam(&entity.Group{ID: exampleTeamUUID, Name: exampleTeamName})
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		team.Deactivate()
		team.Activate()

		want := 3
		got := len(team.Events())

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

func TestTeam_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		team := NewTeamFromEvents([]event.Event{
			&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			&event.TeamDeactivated{ID: exampleTeamUUID},
			&event.TeamActivated{ID: exampleTeamUUID},
		})

		want := 3
		got := team.Version()

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
