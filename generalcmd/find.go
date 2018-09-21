package generalcmd

import (
	"bytes"
	"context"
	msg "gamelink-go/proto_msg"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/service"
	"html/template"
)

type (
	//FindFabric - strucet for find fabric
	FindFabric struct{}
	//FindCommand - struct for find command
	FindCommand struct {
		params []*msg.OneCriteriaStruct
		res    iface.Responder
	}
)

const (
	//commandFind - const for command
	commandFind = "find"
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(FindFabric{})
}

//CommandName - return human readable command name
func (c FindFabric) CommandName() string {
	return commandFind
}

//RequireAdmin - func for checking if admin permissions required
func (f FindFabric) RequireAdmin() bool {
	return false
}

//Require - return array of needed permissions
func (f FindFabric) Require() []string {
	return []string{commandFind}
}

//TryParse - func for parsing request
func (c FindFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command FindCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commandFind); err != nil {
		if err == service.UnknownCommandError {
			return nil, nil
		}
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (fc FindCommand) Execute(ctx context.Context) {
	r, err := Executor().Find(ctx, fc.params)
	if err != nil {
		fc.res.Respond(err.Error())
		return
	}
	t := template.Must(template.New("user.html").ParseFiles("generalcmd/template/user.html"))
	buf := new(bytes.Buffer)
	err = t.Execute(buf, r)
	if err != nil {
		fc.res.Respond(err.Error())
		return
	}
	fc.res.Respond(buf.String())
}
