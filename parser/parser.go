package parser

import (
	"errors"
	"gamelinkBot/iface"
	"sync"
)

type (

	//CommandParser - struct for command parser
	CommandParser struct {
		fabrics []iface.CommandFabric
		checker iface.PermChecker // проверяет привелегии
	}
)

var (
	parser     iface.Parser
	parserOnce sync.Once
)

//SharedParser - singleton for parser
func SharedParser() iface.Parser {
	parserOnce.Do(func() {
		parser = &CommandParser{}
	})
	return parser
}

//SetChecker - add permission checker to command parser
func (p *CommandParser) SetChecker(pc iface.PermChecker) {
	if p != nil {
		p.checker = pc
	}
}

//RegisterFabric - append fabric to command parser fabrics array
func (p *CommandParser) RegisterFabric(cf iface.CommandFabric) {
	if p != nil {
		p.fabrics = append(p.fabrics, cf)
	}
}

//TryParse - tries to parse the request in a loop using array factories
func (p CommandParser) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	if p.checker == nil {
		return nil, errors.New("permission checker is not defined")
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
