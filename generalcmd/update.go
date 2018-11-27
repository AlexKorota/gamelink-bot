package generalcmd

import (
	"context"
	"gamelink-bot/command_list"
	"gamelink-bot/iface"
	"gamelink-bot/parser"
	"gamelink-bot/service"
	msg "gamelink-go/proto_msg"
)

type (
	//UpdateFabric - struct for update fabric
	UpdateFabric struct{}
	//UpdateCommand - struct for update command
	UpdateCommand struct {
		findParams []*msg.OneCriteriaStruct
		updParams  []*msg.UpdateCriteriaStruct
		res        iface.Responder
	}
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(UpdateFabric{})
}

//CommandName - return human readable command name
func (c UpdateFabric) CommandName() string {
	return command_list.CommandUpdate
}

//RequireAdmin - func for checking if admin permissions required
func (c UpdateFabric) RequireAdmin() bool {
	return false
}

//Require - return array of needed permissions
func (c UpdateFabric) Require() []string {
	return []string{command_list.CommandUpdate}
}

//TryParse - func for parsing request
func (c UpdateFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command UpdateCommand
		err     error
	)
	if command.findParams, command.updParams, _, err = service.CompareParseCommand(req.Request(), "/"+command_list.CommandUpdate); err != nil {
		if err == service.UnknownCommandError {
			return nil, nil
		}
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (cc UpdateCommand) Execute(ctx context.Context) {
	r, err := Executor().Update(ctx, cc.findParams, cc.updParams)
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}
