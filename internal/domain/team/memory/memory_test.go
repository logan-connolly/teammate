package memory

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/domain/team"
	"github.com/logan-connolly/teammate/internal/entity"
)

var (
	exampleUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherUUID = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	exampleName = "Syracuse"
	anotherName = "Notre Dame"
)

func TestMemoryTeamRepository_Get(t *testing.T) {
	type testCase struct {
		test        string
		group       *entity.Group
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "No team found with this group",
			group:       &entity.Group{ID: anotherUUID, Name: anotherName},
			expectedErr: team.ErrTeamNotFound,
		}, {
			test:        "Team found",
			group:       &entity.Group{ID: exampleUUID, Name: exampleName},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			repo := NewMemoryTeamRepository()
			repo.teams[exampleUUID] = []team.Event{
				&team.TeamRegistered{ID: exampleUUID, Name: exampleName},
			}

			_, err := repo.Get(tc.group)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemoryTeamRepository_Add(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		name        string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Successfully add a team",
			id:          anotherUUID,
			name:        anotherName,
			expectedErr: nil,
		},
		{
			test:        "Team already exists error",
			id:          exampleUUID,
			name:        exampleName,
			expectedErr: team.ErrTeamAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryTeamRepository()
			team := team.NewTeamFromEvents([]team.Event{
				&team.TeamRegistered{ID: tc.id, Name: tc.name},
			})
			r.teams[exampleUUID] = team.Events()

			err := r.Add(team)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemoryTeamRepository_Update(t *testing.T) {
	type testCase struct {
		test        string
		register    bool
		deactivate  bool
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Update team",
			register:    true,
			deactivate:  true,
			expectedErr: nil,
		},
		{
			test:        "Team has no changes",
			register:    true,
			deactivate:  false,
			expectedErr: team.ErrTeamHasNoUpdates,
		},
		{
			test:        "Team not found",
			register:    false,
			deactivate:  true,
			expectedErr: team.ErrTeamNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryTeamRepository()
			team := team.NewTeamFromEvents([]team.Event{
				&team.TeamRegistered{ID: exampleUUID, Name: exampleName},
			})
			if tc.register {
				r.teams[exampleUUID] = team.Events()
			}
			if tc.deactivate {
				team.Deactivate()
			}

			err := r.Update(team)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
