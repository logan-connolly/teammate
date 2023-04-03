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
			is := is.New(t)
			_, err := NewTeam(tc.group)
			is.Equal(err, tc.expectedErr)
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
			is := is.New(t)
			team := NewTeamFromEvents(tc.events)
			is.Equal(team.GetID(), exampleTeamUUID)
			is.Equal(team.GetName(), exampleTeamName)
			is.Equal(team.IsActivated(), tc.expectedIsActivated)
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
			is := is.New(t)
			err := tc.team.Activate()
			is.Equal(err, tc.expectedErr)
			is.True(tc.team.IsActivated())
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
			is := is.New(t)
			err := tc.team.Deactivate()
			is.Equal(err, tc.expectedErr)
			is.Equal(tc.team.IsActivated(), false)
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
			is := is.New(t)
			err := tc.team.AssignPlayer(tc.player)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestTeam_Unassignplayer(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Unassign player to team",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
				&event.PlayerAssignedToTeam{ID: exampleTeamUUID, PlayerId: examplePlayerUUID, PlayerName: examplePlayerName},
			}),
			expectedErr: nil,
		},
		{
			test: "Try to unassigned a player that is not assigned to team",
			player: NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: examplePlayerName},
			}),
			team: NewTeamFromEvents([]event.Event{
				&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			}),
			expectedErr: ErrPlayerNotAssignedToTeam,
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
			is := is.New(t)
			got := len(tc.team.GetPlayers())
			is.Equal(got, tc.playerCount)
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

		got := len(team.Events())

		is.Equal(got, 3)
	})
}

func TestTeam_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		is := is.New(t)
		team := NewTeamFromEvents([]event.Event{
			&event.TeamCreated{ID: exampleTeamUUID, Name: exampleTeamName},
			&event.TeamDeactivated{ID: exampleTeamUUID},
			&event.TeamActivated{ID: exampleTeamUUID},
		})

		got := team.Version()

		is.Equal(got, 3)
	})
}
