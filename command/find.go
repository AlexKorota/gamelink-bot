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

func (c FindFabric) TryParse(req RequesterResponder) (Command, error) {
	var command FindCommand
	ind := strings.Index(req.Request(), " ")
	if ind < 0 || req.Request()[:ind] != "/find" {
		return nil, nil
	}
	params := strings.Split(req.Request()[ind+1:], " ")
	service.ParseRequest(params)
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
