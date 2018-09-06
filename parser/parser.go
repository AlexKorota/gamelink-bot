package parser

import (
	"context"
	"errors"
	"gamelinkBot/bot"
	"sync"
)

type (
	//Command - command interface
	Command interface {
		Execute(ctx context.Context)
	}
	//CommandFabric - command fabric interface
	CommandFabric interface {
		TryParse(req bot.RequesterResponder) (Command, error)
		RequireAdmin() bool
		Require() []string
	}
	//Parser - parser inerface
	Parser interface {
		SetChecker(pc PermChecker)
		RegisterFabric(cf CommandFabric)
		TryParse(req bot.RequesterResponder) (Command, error)
	}
	//PermChecker - permission checker interface
	PermChecker interface {
		IsAdmin(userName string) (bool, error)
		HasPermissions(userName string, permissions []string) (bool, error)
	}
	//CommandParser - struct for command parser
	CommandParser struct {
		fabrics []CommandFabric
		checker PermChecker // проверяет привелегии
	}
)

var (
	parser     Parser
	parserOnce sync.Once
)

//SharedParser - singleton for parser
func SharedParser() Parser {
	parserOnce.Do(func() {
		parser = &CommandParser{}
	})
	return parser
}

//SetChecker - add permission checker to command parser
func (p *CommandParser) SetChecker(pc PermChecker) {
	if p != nil {
		p.checker = pc
	}
}

//RegisterFabric - append fabric to command parser fabrics array
func (p *CommandParser) RegisterFabric(cf CommandFabric) {
	if p != nil {
		p.fabrics = append(p.fabrics, cf)
	}
}

//TryParse - tries to parse the request in a loop using array factories
func (p CommandParser) TryParse(req bot.RequesterResponder) (Command, error) {
	if p.checker == nil {
		return nil, errors.New("permission checked is not defined")
	}
	adm, err := p.checker.IsAdmin(req.UserName())
	if err != nil {
		return nil, err
	}
	for _, v := range p.fabrics {
		if v.RequireAdmin() && !adm {
			return nil, errors.New("permission denied")
		}
		cmd, err := v.TryParse(req)
		if err != nil {
			continue
			//return nil, err // А зачем нам тут ретурн?
		}
		if cmd != nil {
			if adm {
				return cmd, nil
			}
			allowed, err := p.checker.HasPermissions(req.UserName(), v.Require())
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, errors.New("permission denied")
			}
			return cmd, nil
		}
	}
	return nil, errors.New("can't recognise command")
}
