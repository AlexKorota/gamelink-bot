package command

import (
	"context"
	"gamelinkBot/prot"
	"gamelinkBot/service"
)

type (
	FindFabric struct{}

	FindCommand struct {
		params []*prot.OneCriteriaStruct
		res    Responder
	}
)

const (
	commandFind = "find"
)

func init() {
	SharedParser().RegisterFabric(FindFabric{})
}

func (f FindFabric) RequireAdmin() bool {
	return false
}

func (f FindFabric) Require() []string {
	return []string{commandFind}
}

func (c FindFabric) TryParse(req RequesterResponder) (Command, error) {
	var (
		command FindCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commandFind); err != nil {
		return nil, err
	}
	command.res = req
	return command, nil
}

func (fc FindCommand) Execute(ctx context.Context) {
	r, err := SharedClient().Find(ctx, &prot.MultiCriteriaRequest{Params: fc.params})
	if err != nil {
		fc.res.Respond(err.Error())
		return
	}
	fc.res.Respond(r.String())
}
