package command

import (
	"context"
	"gamelinkBot/prot"
	"gamelinkBot/service"
)

type (
	DeleteFabric  struct{}
	DeleteCommand struct {
		params []*prot.OneCriteriaStruct
		res    Responder
	}
)

const (
	commandDelete = "count"
)

func init() {
	SharedParser().RegisterFabric(DeleteFabric{})
}

func (c DeleteFabric) RequireAdmin() bool {
	return false
}

func (c DeleteFabric) Require() []string {
	return []string{commandDelete}
}

func (c DeleteFabric) TryParse(req RequesterResponder) (Command, error) {
	var (
		command DeleteCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commandDelete); err != nil {
		return nil, err
	}
	command.res = req
	return command, nil
}

func (cc DeleteCommand) Execute(ctx context.Context) {
	r, err := SharedClient().Count(ctx, &prot.MultiCriteriaRequest{Params: cc.params})
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}
