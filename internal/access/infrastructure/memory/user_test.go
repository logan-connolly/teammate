package memory

import (
	"testing"

	"git.sr.ht/~loges/teammate/internal/access/domain/event"
	"git.sr.ht/~loges/teammate/internal/access/domain/model"
	"git.sr.ht/~loges/teammate/internal/access/domain/repository"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

var (
	exampleUUID  = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherUUID  = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	exampleName  = "Mike Ditka"
	exampleEmail = "ditka@teammate.com"
	anotherName  = "Joe Gibbs"
	anotherEmail = "gibbs@teammate.com"
)

func TestMemoryAccessRepository_GetByEmail(t *testing.T) {
	type testCase struct {
		test        string
		email       string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "No user with this entity",
			email:       anotherEmail,
			expectedErr: repository.ErrUserNotFound,
		},
		{
			test:        "User found",
			email:       exampleEmail,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			repo := NewMemoryUserRepository()
			repo.users[exampleEmail] = []event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			}

			_, err := repo.GetByEmail(tc.email)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestMemoryAccessRepository_Add(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		name        string
		email       string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Successfully add a user",
			id:          anotherUUID,
			name:        anotherName,
			email:       anotherEmail,
			expectedErr: nil,
		},
		{
			test:        "User already exists error",
			id:          exampleUUID,
			name:        exampleName,
			email:       exampleEmail,
			expectedErr: repository.ErrUserAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			r := NewMemoryUserRepository()
			u := model.NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: tc.id, Name: tc.name, Email: tc.email},
			})
			r.users[exampleEmail] = u.Events()

			err := r.Add(u)

			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestMemoryAccessRepository_Update(t *testing.T) {
	type testCase struct {
		test        string
		register    bool
		deactivate  bool
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Update user",
			register:    true,
			deactivate:  true,
			expectedErr: nil,
		},
		{
			test:        "User has no changes",
			register:    true,
			deactivate:  false,
			expectedErr: repository.ErrUserHasNoUpdates,
		},
		{
			test:        "User not found",
			register:    false,
			deactivate:  true,
			expectedErr: repository.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			r := NewMemoryUserRepository()
			u := model.NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			})
			if tc.register {
				r.users[exampleEmail] = u.Events()
			}
			if tc.deactivate {
				u.Deactivate()
			}

			err := r.Update(u)

			is.Equal(err, tc.expectedErr)
		})
	}
}
