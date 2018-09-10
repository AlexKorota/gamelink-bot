package command

import (
	"context"
	"gamelinkBot/bot"
	"gamelinkBot/parser"
	"gamelinkBot/prot"
	"gamelinkBot/rpc"
	"gamelinkBot/service"
)

type (
	//FindFabric - strucet for find fabric
	FindFabric struct{}
	//FindCommand - struct for find command
	FindCommand struct {
		params []*prot.OneCriteriaStruct
		res    bot.Responder
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

//RequireAdmin - func for checking if admin permissions required
func (f FindFabric) RequireAdmin() bool {
	return false
}

//Require - return array of needed permissions
func (f FindFabric) Require() []string {
	return []string{commandFind}
}

//TryParse - func for parsing request
func (c FindFabric) TryParse(req bot.RequesterResponder) (parser.Command, error) {
	var (
		command FindCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commandFind); err != nil {
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (fc FindCommand) Execute(ctx context.Context) {
	r, err := rpc.SharedClient().Find(ctx, &prot.MultiCriteriaRequest{Params: fc.params})
	if err != nil {
		fc.res.Respond(err.Error())
		return
	}
	fc.res.Respond(r.String())
}
