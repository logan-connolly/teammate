package registration

import (
	"errors"
	"testing"
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
