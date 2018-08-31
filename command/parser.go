package command

import (
	"context"
	"errors"
	"sync"
)

type (
	Command interface {
		Execute(ctx context.Context)
	}

	CommandFabric interface {
		TryParse(req RequesterResponder) (Command, error)
		RequireAdmin() bool
		Require() []string
	}

	Parser interface {
		SetChecker(pc PermChecker)
		RegisterFabric(cf CommandFabric)
		TryParse(req RequesterResponder) (Command, error)
	}

	PermChecker interface {
		IsAdmin(userName string) (bool, error)
		HasPermissions(userName string, permissions []string) (bool, error)
	}

	CommandParser struct {
		fabrics []CommandFabric
		checker PermChecker // проверяет привелегии
	}
)

var (
	parser     Parser
	parserOnce sync.Once
)

func SharedParser() Parser {
	parserOnce.Do(func() {
		parser = &CommandParser{}
	})
	return parser
}

func (p *CommandParser) SetChecker(pc PermChecker) {
	if p != nil {
		p.checker = pc
	}
}

func (p *CommandParser) RegisterFabric(cf CommandFabric) {
	if p != nil {
		p.fabrics = append(p.fabrics, cf)
	}
}

func (p CommandParser) TryParse(req RequesterResponder) (Command, error) {
	if p.checker == nil {
		return nil, errors.New("permission checked is not defined")
	}
	for _, v := range p.fabrics {
		cmd, err := v.TryParse(req)
		if err != nil {
			return nil, err
		}
		var adm bool
		if v.RequireAdmin() {
			adm, err = p.checker.IsAdmin(req.UserName())
			if err != nil {
				return nil, err
			}
			if !adm {
				return nil, errors.New("permission denied")
			}
		} else {
			allowed, err := p.checker.HasPermissions(req.UserName(), v.Require())
			if err != nil {
				return nil, err
			}
			if !(adm || allowed) {
				return nil, errors.New("permission denied")
			}
		}
		if cmd != nil {
			return cmd, nil
		}
	}
	return nil, errors.New("can't recognise command")
}
