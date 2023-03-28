package services

import (
	"errors"
	"testing"
)

const exampleName = "Matt"

func TestNewRegistrationService(t *testing.T) {
	type testCase struct {
		test        string
		testConfig  func() RegistrationConfiguration
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "With in-memory player repository",
			testConfig:  WithMemoryPlayerRepository,
			expectedErr: nil,
		},
		{
			test:        "With in-memory team repository",
			testConfig:  WithMemoryTeamRepository,
			expectedErr: nil,
		},
		{
			test: "With bad config",
			testConfig: func() RegistrationConfiguration {
				return func(s *RegistrationService) error {
					return ErrInvalidRegistrationConfig
				}
			},
			expectedErr: ErrInvalidRegistrationConfig,
		},
	}

	for _, tc := range testCases {
		_, err := NewRegistrationService(tc.testConfig())
		if !errors.Is(err, tc.expectedErr) {
			t.Errorf("Expected %v, but got %v.", tc.expectedErr, err)
		}
	}
}
