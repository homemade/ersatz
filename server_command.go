package main

import "fmt"

type ServerCommand struct {
	Command               string                `json:"command"`
	VariableEndpointIndex VariableEndpointIndex `json:"endpoint"`
}

func (c ServerCommand) Execute(s *ServerApp) error {
	switch c.Command {
	case "vary":
		s.EndpointVariationSchedule[c.VariableEndpointIndex.EndpointIndex] = EndpointVariation{
			Variation: c.VariableEndpointIndex.Variant,
			Count:     1,
		}

	default:
		return fmt.Errorf("Unknown command '%s'", c.Command)
	}

	return nil
}
