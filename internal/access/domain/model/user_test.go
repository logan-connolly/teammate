package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/access/domain/event"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/matryer/is"
)

var (
	exampleUUID  = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	exampleName  = "Mike Ditka"
	exampleEmail = "ditka@teammate.com"
	anotherName  = "Joe Gibbs"
	anotherEmail = "gibbs@teammate.com"
)

func TestUser_NewUser(t *testing.T) {
	type testCase struct {
		test        string
		person      *entity.Person
		email       string
		expectedErr error
	}
	testCases := []testCase{
		{
			test:        "Empty name validation",
			person:      &entity.Person{ID: exampleUUID, Name: ""},
			email:       exampleEmail,
			expectedErr: ErrUserNameIsEmpty,
		},
		{
			test:        "Empty name validation",
			person:      &entity.Person{ID: exampleUUID, Name: exampleName},
			email:       "",
			expectedErr: ErrUserEmailIsEmpty,
		},
		{
			test:        "Valid email and name",
			person:      &entity.Person{ID: exampleUUID, Name: exampleName},
			email:       exampleEmail,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			_, err := NewUser(tc.person, tc.email)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestUser_NewUserEvents(t *testing.T) {
	type testCase struct {
		test                string
		events              []event.Event
		expectedIsActivated bool
	}
	testCases := []testCase{
		{
			test: "User registered",
			events: []event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			},
			expectedIsActivated: true,
		},
		{
			test: "User registered and deactivated",
			events: []event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
				&event.UserDeactivated{ID: exampleUUID},
			},
			expectedIsActivated: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			u := NewUserFromEvents(tc.events)
			is.Equal(u.GetID(), exampleUUID)
			is.Equal(u.GetName(), exampleName)
			is.Equal(u.GetEmail(), exampleEmail)
			is.Equal(u.IsActivated(), tc.expectedIsActivated)
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	type testCase struct {
		test        string
		user        *User
		newName     string
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Update name",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			}),
			newName:     anotherName,
			expectedErr: nil,
		},
		{
			test: "Update name with same name",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
				&event.UserDeactivated{ID: exampleUUID},
			}),
			newName:     exampleName,
			expectedErr: ErrUserNameAlreadySetToValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.UpdateName(tc.newName)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	type testCase struct {
		test        string
		user        *User
		newEmail    string
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Update email",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			}),
			newEmail:    anotherEmail,
			expectedErr: nil,
		},
		{
			test: "Update email with same email",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
				&event.UserDeactivated{ID: exampleUUID},
			}),
			newEmail:    exampleEmail,
			expectedErr: ErrUserEmailAlreadySetToValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.UpdateEmail(tc.newEmail)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestUser_Activate(t *testing.T) {
	type testCase struct {
		test        string
		user        *User
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Activate active user",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			}),
			expectedErr: ErrUserAlreadyActivated,
		},
		{
			test: "Activate deactivated user",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
				&event.UserDeactivated{ID: exampleUUID},
			}),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.Activate()
			is.Equal(err, tc.expectedErr)
			is.True(tc.user.IsActivated())
		})
	}
}

func TestUser_Deactivate(t *testing.T) {
	type testCase struct {
		test        string
		user        *User
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Deactivate active user",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			}),
			expectedErr: nil,
		},
		{
			test: "Deactivate deactivated user",
			user: NewUserFromEvents([]event.Event{
				&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
				&event.UserDeactivated{ID: exampleUUID},
			}),
			expectedErr: ErrUserAlreadyDeactivated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.Deactivate()
			is.Equal(err, tc.expectedErr)
			is.Equal(tc.user.IsActivated(), false)
		})
	}
}

func TestUser_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		is := is.New(t)
		u, err := NewUser(&entity.Person{ID: exampleUUID, Name: exampleName}, exampleEmail)
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		u.Deactivate()
		u.Activate()

		got := len(u.Events())

		is.Equal(got, len(u.Events()))
	})
}

func TestUser_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		is := is.New(t)
		u := NewUserFromEvents([]event.Event{
			&event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail},
			&event.UserDeactivated{ID: exampleUUID},
			&event.UserActivated{ID: exampleUUID},
		})

		got := u.Version()

		is.Equal(got, 3)
	})
}
