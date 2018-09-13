package generalcmd

import (
	"context"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/service"
)

type (
	//CountFabic - struct for Count fabric
	CountFabric struct{}
	//CountCommand - struct for count command
	CountCommand struct {
		params []*iface.OneCriteriaStruct
		res    iface.Responder
	}
)

const (
	//commandCount - const for command name
	commandCount = "count"
)

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
	return []string{commandCount}
}

//CommandName - return human readable command name
func (c CountFabric) CommandName() string {
	return commandCount
}

//TryParse - func for parsing request
func (c CountFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command CountCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commandCount); err != nil {
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
	cc.res.Respond(r.String())
}
