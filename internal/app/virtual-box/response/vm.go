package response

import "github.com/terra-farm/go-virtualbox"

type (
	StatusResponse struct {
		Status []Status `json:"status"`
	}

	Status struct {
		Name   string                  `json:"name"`
		CPUs   uint                    `json:"cpus"`
		Memory uint                    `json:"memory"`
		State  virtualbox.MachineState `json:"state"`
	}
)
