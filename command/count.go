package command

import (
	"context"
	"fmt"
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

func (c CountFabric) Command() string {
	return "/count"
}

func (c CountFabric) TryParse(req RequesterResponder) (Command, error) {
	var command CountCommand
	var err error
	ind := strings.Index(req.Request(), " ")
	if ind < 0 {
		return nil, nil
	}
	if req.Request()[:ind] != c.Command() {
		return nil, nil
	}
	command.params, err = service.ParseRequest(req.Request())
	if err != nil {
		return nil, err
	}
	command.userName = req.UserName()
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
