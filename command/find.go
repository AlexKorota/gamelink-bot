package command

import (
	"context"
	"gamelinkBot/prot"
	"gamelinkBot/service"
	"strings"
)

type FindFabric struct{}

type FindCommand struct {
	params   []*prot.OneCriteriaStruct
	res      Responder
	userName string
}

func init() {
	SharedParser().RegisterFabric(FindFabric{})
}

func (f FindFabric) RequireAdmin() bool {
	return false
}

func (f FindFabric) Require() []string {
	return []string{"find"}
}

func (f FindFabric) Command() string {
	return "/find"
}

func (f FindFabric) TryParse(req RequesterResponder) (Command, error) {
	var command FindCommand
	var err error
	ind := strings.Index(req.Request(), " ")
	if ind < 0 {
		return nil, nil
	}
	if req.Request()[:ind] != f.Command() {
		return nil, nil
	}
	command.params, err = service.ParseRequest(req.Request())
	if err != nil {
		return nil, nil
	}
	command.userName = req.UserName()
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
