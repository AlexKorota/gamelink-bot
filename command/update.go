package command

import (
	"context"
	"gamelinkBot/prot"
	"gamelinkBot/service"
)

type (
	UpdateFabric  struct{}
	UpdateCommand struct {
		params []*prot.OneCriteriaStruct
		res    Responder
	}
)

const (
	commandUpdate = "update"
)

func init() {
	SharedParser().RegisterFabric(UpdateFabric{})
}

func (c UpdateFabric) RequireAdmin() bool {
	return false
}

func (c UpdateFabric) Require() []string {
	return []string{commandUpdate}
}

func (c UpdateFabric) TryParse(req RequesterResponder) (Command, error) {
	var (
		command UpdateCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commandUpdate); err != nil {
		return nil, err
	}
	command.res = req
	return command, nil
}

func (cc UpdateCommand) Execute(ctx context.Context) {
	r, err := SharedClient().Count(ctx, &prot.MultiCriteriaRequest{Params: cc.params})
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}
