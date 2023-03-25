package registration

import (
	"errors"
	"testing"

	"github.com/logan-connolly/teammate/internal/domain/player"
)

const exampleName = "Matt"

func TestNewRegistrationService(t *testing.T) {
	ErrInvalidConfig := errors.New("Invalid configuration")

	type testCase struct {
		test        string
		testConfig  func() RegistrationConfiguration
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "With in-memory repository",
			testConfig:  WithMemoryPlayerRepository,
			expectedErr: nil,
		},
		{
			test: "With bad config",
			testConfig: func() RegistrationConfiguration {
				return func(s *RegistrationService) error {
					return ErrInvalidConfig
				}
			},
			expectedErr: ErrInvalidConfig,
		},
	}

	for _, tc := range testCases {
		_, err := NewRegistrationService(tc.testConfig())
		if !errors.Is(err, tc.expectedErr) {
			t.Errorf("Expected %v, but got %v.", tc.expectedErr, err)
		}
	}
}

func TestRegistrationService_RegisterPlayer(t *testing.T) {
	type testCase struct {
		test         string
		testRegister func(s *RegistrationService) error
		expectedErr  error
	}

	testCases := []testCase{
		{
			test: "Register player",
			testRegister: func(s *RegistrationService) error {
				_, err := s.RegisterPlayer(exampleName)
				return err
			},
			expectedErr: nil,
		},
		{
			test: "Register player with invalid name",
			testRegister: func(s *RegistrationService) error {
				_, err := s.RegisterPlayer("")
				return err
			},
			expectedErr: player.ErrInvalidPerson,
		},
	}

	s, err := NewRegistrationService(
		WithMemoryPlayerRepository(),
	)
	if err != nil {
		t.Error(err)
	}

	for _, tc := range testCases {
		got := tc.testRegister(s)
		if !errors.Is(got, tc.expectedErr) {
			t.Errorf("Expected %v, but got %v.", tc.expectedErr, got)
		}
	}
}
