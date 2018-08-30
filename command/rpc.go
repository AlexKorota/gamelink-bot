package command

import (
	"context"
	"fmt"
	"gamelinkBot/prot"
)

type (
	RPCExecutable interface {
		ExecuteInContext(ctx context.Context) (fmt.Stringer, error)
		Command() string
		Params() []*prot.OneCriteriaStruct
	}

	RPCCommand struct {
		command string
		params  []*prot.OneCriteriaStruct
	}
)

func (rpcc RPCCommand) ExecuteInContext(ctx context.Context) (fmt.Stringer, error) {
	return nil, nil
}

func (rpcc RPCCommand) Command() string {
	return rpcc.command
}

func (rpcc RPCCommand) Params() []*prot.OneCriteriaStruct {
	return rpcc.params
}
