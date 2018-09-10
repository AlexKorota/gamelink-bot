package generalcmd

import (
	"context"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/prot"
	"gamelinkBot/service"
)

type (
	//UpdateFabric - struct for update fabric
	UpdateFabric struct{}
	//UpdateCommand - struct for update command
	UpdateCommand struct {
		params []*prot.OneCriteriaStruct
		res    iface.Responder
	}
)

const (
	//commandUpdate - const for command name
	commandUpdate = "update"
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(UpdateFabric{})
}

//RequireAdmin - func for checking if admin permissions required
func (c UpdateFabric) RequireAdmin() bool {
	return false
}

//Require - return array of needed permissions
func (c UpdateFabric) Require() []string {
	return []string{commandUpdate}
}

//TryParse - func for parsing request
func (c UpdateFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
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

//Execute - execute command
func (cc UpdateCommand) Execute(ctx context.Context) {
	r, err := SharedClient().Count(ctx, &prot.MultiCriteriaRequest{Params: cc.params})
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}