package command

import (
	"context"
	"fmt"
	"gamelinkBot/prot"
)

type (
	Executable interface {
		ExecuteInContext(ctx context.Context) (fmt.Stringer, error)
		Command() string
		CommandType() string
		Params() []*prot.OneCriteriaStruct
	}

	baseCommand struct {
		command     string
		commandType string // А может лучше енум?
		params      []*prot.OneCriteriaStruct
	}
)

func (bc baseCommand) ExecuteInContext(ctx context.Context) (fmt.Stringer, error) {
	return nil, nil
}

func (bc baseCommand) Command() string {
	return bc.command
}

func (bc baseCommand) Params() []*prot.OneCriteriaStruct {
	return bc.params
}

func (bc baseCommand) CommandType() string {
	return bc.commandType
}
