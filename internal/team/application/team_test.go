package application

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/application/services"
	"github.com/matryer/is"
)

var (
	exampleGroup  = &entity.Group{ID: uuid.MustParse("a45e93f8-c952-11ed-afa1-0242ac120002"), Name: "Lions"}
	examplePerson = &entity.Person{ID: uuid.MustParse("d17ac10b-58cc-0372-8567-0e02b2c3d479"), Name: "Jen"}
)

func withInvalidConfig() services.RosterConfiguration {
	return func(s *services.RosterService) error {
		return services.ErrInvalidRosterConfig
	}
}

func TestTeamApplication(t *testing.T) {
	t.Run("Init team app", func(t *testing.T) {
		is := is.New(t)
		_, err := NewTeamApplication()
		is.NoErr(err)
	})

	t.Run("Init failure due to bad roster service config", func(t *testing.T) {
		is := is.New(t)
		originalConfigs := services.RosterConfigs
		services.RosterConfigs = []services.RosterConfiguration{withInvalidConfig()}

		_, err := NewTeamApplication()

		is.Equal(err, services.ErrInvalidRosterConfig)
		// clean up configs
		services.RosterConfigs = originalConfigs
	})

	t.Run("Roster service workflow", func(t *testing.T) {
		is := is.New(t)
		ta, err := NewTeamApplication()
		rs := ta.GetRosterService()

		err = rs.AddTeam(exampleGroup)
		is.NoErr(err)

		err = rs.AddPlayer(examplePerson)
		is.NoErr(err)

		err = rs.AssignPlayerToTeam(exampleGroup, examplePerson)
		is.NoErr(err)

		err = rs.UnassignPlayerFromTeam(exampleGroup, examplePerson)
		is.NoErr(err)
	})
}
