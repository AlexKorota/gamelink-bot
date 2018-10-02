package generalcmd

import (
	"context"
	msg "gamelink-go/proto_msg"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/service"
)

type (
	//SendFabric - struct for send struct
	SendFabric struct{}
	//SendCommand - struct for send command
	SendCommand struct {
		params  []*msg.OneCriteriaStruct
		message string
		res     iface.Responder
	}
)

const (
	sendPushCommand = "send_push"
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(SendFabric{})
}

//RequireAdmin - func for checking if admin permissions required
func (c SendFabric) RequireAdmin() bool {
	return false
}

//Require - return array of needed permissions
func (c SendFabric) Require() []string {
	return []string{sendPushCommand}
}

//CommandName - return human readable command name
func (c SendFabric) CommandName() string {
	return sendPushCommand
}

//TryParse - func for parsing request
func (c SendFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command SendCommand
		err     error
	)
	if command.params, _, command.message, err = service.CompareParseCommand(req.Request(), "/"+sendPushCommand); err != nil {
		if err == service.UnknownCommandError {
			return nil, nil
		}
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (sc SendCommand) Execute(ctx context.Context) {
	r, err := Executor().SendPush(ctx, sc.params, sc.message)
	if err != nil {
		sc.res.Respond(err.Error())
		return
	}
	sc.res.Respond(r.String())
}
