package team

import (
	"github.com/logan-connolly/teammate/internal/team/services"
)

var memoryConfigs = []services.RosterConfiguration{
	services.WithMemoryPlayerRepository(),
	services.WithMemoryTeamRepository(),
}

type Application struct {
	rosterService services.RosterService
}

func NewMemoryApplication() (*Application, error) {
	rs, err := services.NewRosterService(memoryConfigs...)
	if err != nil {
		return &Application{}, services.ErrInvalidRosterConfig
	}

	return &Application{
		rosterService: *rs,
	}, nil
}
