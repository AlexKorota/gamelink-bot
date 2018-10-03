package generalcmd

import (
	"context"
	"fmt"
	msg "gamelink-go/proto_msg"
	"gamelinkBot/command_list"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/service"
	"strings"
)

type (
	//CountFabic - struct for Count fabric
	CountFabric struct{}
	//CountCommand - struct for count command
	CountCommand struct {
		params []*msg.OneCriteriaStruct
		res    iface.Responder
	}
)

//
//const (
//	//commandCount - const for command name
//	commandCount = "count"
//)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(CountFabric{})
}

//RequireAdmin - func for checking if admin permissions required
func (c CountFabric) RequireAdmin() bool {
	return false
}

//Require - return array of needed permissions
func (c CountFabric) Require() []string {
	return []string{command_list.CommandCount}
}

//CommandName - return human readable command name
func (c CountFabric) CommandName() string {
	return command_list.CommandCount
}

//TryParse - func for parsing request
func (c CountFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command CountCommand
		err     error
	)
	if strings.Trim(req.Request(), " ") == "/"+command_list.CommandCount {
		command.res = req
		return command, nil
	}
	if command.params, _, _, err = service.CompareParseCommand(req.Request(), "/"+command_list.CommandCount); err != nil {
		if err == service.UnknownCommandError {
			return nil, nil
		}
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (cc CountCommand) Execute(ctx context.Context) {
	r, err := Executor().Count(ctx, cc.params)
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(fmt.Sprintf("%d", r.Count))
}
