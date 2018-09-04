package command

import (
	"context"
	"gamelinkBot/prot"
	"gamelinkBot/service"
)

type (
	CountFabric  struct{}
	CountCommand struct {
		params []*prot.OneCriteriaStruct
		res    Responder
	}
)

const (
	commanStr = "count"
)

func init() {
	SharedParser().RegisterFabric(CountFabric{})
}

func (c CountFabric) RequireAdmin() bool {
	return false
}

func (c CountFabric) Require() []string {
	return []string{commanStr}
}

func (c CountFabric) TryParse(req RequesterResponder) (Command, error) {
	var (
		command CountCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commanStr); err != nil {
		return nil, err
	}
	command.res = req
	return command, nil
}

func (cc CountCommand) Execute(ctx context.Context) {
	r, err := SharedClient().Count(ctx, &prot.MultiCriteriaRequest{Params: cc.params})
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}
