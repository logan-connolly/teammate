package access

import (
	"testing"

	"git.sr.ht/~loges/teammate/internal/access/application/services"
	"github.com/matryer/is"
)

var (
	exampleName  = "John"
	exampleEmail = "john@teammate.com"
)

func withInvalidMemoryConfig() services.RegistrationConfiguration {
	return func(s *services.RegistrationService) error {
		return services.ErrInvalidRegistrationConfig
	}
}

func TestMemoryApplication(t *testing.T) {
	t.Run("Init access app", func(t *testing.T) {
		is := is.New(t)
		_, err := NewAccessApplication()
		is.NoErr(err)
	})

	t.Run("Init failure due to bad registration service config", func(t *testing.T) {
		is := is.New(t)
		originalConfigs := services.RegistrationConfigs
		services.RegistrationConfigs = []services.RegistrationConfiguration{withInvalidMemoryConfig()}

		_, err := NewAccessApplication()

		is.Equal(err, services.ErrInvalidRegistrationConfig)

		// clean up configs
		services.RegistrationConfigs = originalConfigs
	})

	t.Run("Registration service workflow", func(t *testing.T) {
		is := is.New(t)
		aa, err := NewAccessApplication()
		rs := aa.GetRegistrationService()

		err = rs.RegisterUser(exampleName, exampleEmail)
		is.NoErr(err)
	})
}
