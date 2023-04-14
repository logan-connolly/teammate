package access

import (
	"testing"

	"github.com/logan-connolly/teammate/internal/access/application/services"
	"github.com/matryer/is"
)

func withInvalidMemoryConfig() services.RegistrationConfiguration {
	return func(s *services.RegistrationService) error {
		return services.ErrInvalidRegistrationConfig
	}
}

func TestMemoryApplication(t *testing.T) {
	t.Run("Init access app", func(t *testing.T) {
		is := is.New(t)
		_, err := NewMemoryApplication()
		is.NoErr(err)
	})

	t.Run("Init failure due to bad registration service config", func(t *testing.T) {
		is := is.New(t)
		originalMemoryConfigs := memoryConfigs
		memoryConfigs = []services.RegistrationConfiguration{withInvalidMemoryConfig()}

		_, err := NewMemoryApplication()

		is.Equal(err, services.ErrInvalidRegistrationConfig)

		// clean up memoryConfigs
		memoryConfigs = originalMemoryConfigs
	})
}
