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
		Params() []*prot.OneCriteriaStruct
	}

	baseCommand struct {
		params  []*prot.OneCriteriaStruct
		command string
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
