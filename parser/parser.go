package parser

import (
	"errors"
	"gamelinkBot/iface"
	log "github.com/sirupsen/logrus"
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
	log.Debug("parser.SharedParser")
	parserOnce.Do(func() {
		parser = &CommandParser{}
	})
	return parser
}

//SetChecker - add permission checker to command parser
func (p *CommandParser) SetChecker(pc iface.PermChecker) {
	log.Debug("parser.CommandParser.SetChecker")
	if p != nil {
		p.checker = pc
	}
}

//RegisterFabric - append fabric to command parser fabrics array
func (p *CommandParser) RegisterFabric(cf iface.CommandFabric) {
	log.Debug("parser.CommandParser.RegisterFabric")
	if p != nil {
		p.fabrics = append(p.fabrics, cf)
	}
}

//TryParse - tries to parse the request in a loop using array factories
func (p CommandParser) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	log.Debug("parser.CommandParser.TryParse")
	if p.checker == nil {
		return nil, errors.New("permission checker is not defined")
	}
	adm, err := p.checker.IsAdmin(req.UserName())
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"user": req.UserName(), "result": adm}).Debug("admin check")
	for _, v := range p.fabrics {
		if v.RequireAdmin() && !adm {
			continue
		}
		log.WithField("name", v.CommandName()).Debug("trying to parse command")
		cmd, err := v.TryParse(req)
		if err != nil {
			return nil, err
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
