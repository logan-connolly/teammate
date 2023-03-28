package memory

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
	"github.com/logan-connolly/teammate/internal/team/domain/repository"
)

var (
	exampleTeamUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherTeamUUID = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	exampleTeamName = "Syracuse"
	anotherTeamName = "Notre Dame"
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
			group:       &entity.Group{ID: anotherTeamUUID, Name: anotherTeamName},
			expectedErr: repository.ErrTeamNotFound,
		}, {
			test:        "Team found",
			group:       &entity.Group{ID: exampleTeamUUID, Name: exampleTeamName},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			repo := NewMemoryTeamRepository()
			repo.teams[exampleTeamUUID] = []event.Event{
				&event.TeamRegistered{ID: exampleTeamUUID, Name: exampleTeamName},
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
			id:          anotherTeamUUID,
			name:        anotherTeamName,
			expectedErr: nil,
		},
		{
			test:        "Team already exists error",
			id:          exampleTeamUUID,
			name:        exampleTeamName,
			expectedErr: repository.ErrTeamAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryTeamRepository()
			team := model.NewTeamFromEvents([]event.Event{
				&event.TeamRegistered{ID: tc.id, Name: tc.name},
			})
			r.teams[exampleTeamUUID] = team.Events()

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
			expectedErr: repository.ErrTeamHasNoUpdates,
		},
		{
			test:        "Team not found",
			register:    false,
			deactivate:  true,
			expectedErr: repository.ErrTeamNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryTeamRepository()
			team := model.NewTeamFromEvents([]event.Event{
				&event.TeamRegistered{ID: exampleTeamUUID, Name: exampleTeamName},
			})
			if tc.register {
				r.teams[exampleTeamUUID] = team.Events()
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
