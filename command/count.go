package command

import (
	"context"
	"fmt"
	"gamelinkBot/prot"
	"gamelinkBot/service"
)

type CountFabric struct{}

type CountCommand struct {
	params []*prot.OneCriteriaStruct
	res    Responder
}

func init() {
	SharedParser().RegisterFabric(CountFabric{})
}

func (c CountFabric) RequireAdmin() bool {
	return false
}

func (c CountFabric) Require() []string {
	return []string{"count"}
}

func (c CountFabric) Command() string {
	return "/count"
}

func (c CountFabric) TryParse(req RequesterResponder) (Command, error) {
	var (
		command CountCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/count"); err != nil {
		return nil, err
	}

	command.res = req
	return command, nil
}

func (cc CountCommand) Execute(ctx context.Context) {
	fmt.Println(cc)
	r, err := SharedClient().Count(ctx, &prot.MultiCriteriaRequest{Params: cc.params})
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}
