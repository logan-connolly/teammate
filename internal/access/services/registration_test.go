package services

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewRegistrationService(t *testing.T) {
	type testCase struct {
		test        string
		testConfig  func() RegistrationConfiguration
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "With in-memory user repository",
			testConfig:  WithMemoryUserRepository,
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
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			_, err := NewRegistrationService(tc.testConfig())
			is.Equal(err, tc.expectedErr)
		})
	}
}
