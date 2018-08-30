package command

import (
	"fmt"
)

type (
	AdminExecutable interface {
		Execute() (fmt.Stringer, error)
		Command() string
		Params() []string
	}

	AdminCommand struct {
		command string
		params  []string
	}
)

func (admc AdminCommand) Execute() (fmt.Stringer, error) {
	return nil, nil
}

func (admc AdminCommand) Command() string {
	return admc.command
}

func (admc AdminCommand) Params() []string {
	return admc.params
}
