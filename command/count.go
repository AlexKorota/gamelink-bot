package command

import (
	"context"
	"gamelinkBot/prot"
	"gamelinkBot/service"
	"strings"
)

type CountFabric struct{}

type CountCommand struct {
	params   []*prot.OneCriteriaStruct
	res      Responder
	userName string
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

func (c CountFabric) TryParse(req RequesterResponder) (Command, error) {
	var command CountCommand
	ind := strings.Index(req.Request(), " ")
	if ind < 0 || req.Request()[:ind] != "/count" {
		return nil, nil
	}
	params := strings.Split(req.Request()[ind+1:], " ")
	service.ParseRequest(params)
	command.userName = req.UserName()
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
