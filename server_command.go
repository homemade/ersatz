package main

import (
	"fmt"
	"strings"
)

type ServerCommand struct {
	Command               string                `json:"command"`
	VariableEndpointIndex VariableEndpointIndex `json:"endpoint"`
}

const (
	SERVER_COMMAND_VARY = "vary"
)

func NewServerCommand() *ServerCommand {
	return &ServerCommand{}
}

func (c ServerCommand) Execute(s *ServerApp) error {
	switch c.Command {
	case SERVER_COMMAND_VARY:

		// Trim the leading slash to prevent confusing collisions
		c.VariableEndpointIndex.EndpointIndex.URL = strings.TrimLeft(c.VariableEndpointIndex.EndpointIndex.URL, "/")

		s.EndpointVariationSchedule[c.VariableEndpointIndex.EndpointIndex] = EndpointVariation{
			Variation: c.VariableEndpointIndex.Variant,
			Count:     1,
		}

	default:
		return fmt.Errorf("Unknown command '%s'", c.Command)
	}

	return nil
}
