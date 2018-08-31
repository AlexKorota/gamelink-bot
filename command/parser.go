package command

import (
	"context"
	"errors"
	"gamelinkBot/config"
	"strings"
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
		RegisterFabric(cf CommandFabric)
		TryParse(req RequesterResponder) (Command, error)
	}

	PermChecker interface {
		IsAdmin(username string) bool
		HasPermission(username string) bool
	}

	CommandParser struct {
		fabrics []CommandFabric
		cheker  PermChecker // проверяет привелегии
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

func (p *CommandParser) RegisterFabric(cf CommandFabric) {
	if p != nil {
		p.fabrics = append(p.fabrics, cf)
	}
}

func (p CommandParser) TryParse(req RequesterResponder) (Command, error) {
	for _, v := range p.fabrics {
		//TODO: добавить метод проверки привелегий
		cmd, err := v.TryParse(req)
		if err != nil {
			return nil, err
		}
		if v.RequireAdmin() {
			if !p.IsAdmin(req.UserName()) {
				return nil, errors.New("permission denied")
			}
		} else {
			if !(p.IsAdmin(req.UserName()) || p.HasPermission(req.UserName())) {
				return nil, errors.New("permission denied")
			}
		}
		if cmd != nil {
			return cmd, nil
		}
	}
	return nil, errors.New("can't recognise command")
}

func (p CommandParser) IsAdmin(username string) bool {
	if username == "" {
		return false
	}
	for _, v := range config.SuperAdmin {
		if username == strings.Trim(v, " ") {
			return true
		}
	}
	return false
}

func (p CommandParser) HasPermission(username string) bool {
	return false
}
