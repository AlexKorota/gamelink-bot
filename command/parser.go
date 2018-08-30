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
		RegisterFabric(cf CommandFabric)
		TryParse(req RequesterResponder) (Command, error)
	}

	Cheker interface {
		// как вариант
	}

	CommandParser struct {
		fabrics []CommandFabric
		cheker  Cheker // проверяет привелегии
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
		if cmd != nil {
			return cmd, nil
		}
	}
	return nil, errors.New("can't recognise command")
}
