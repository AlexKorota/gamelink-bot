package command

import (
	"context"
)

type (
	Command interface {
		Execute(ctx context.Context) error
	}

	CommandFabric interface {
		TryParse(req RequesterResponder) (Command, error)
		RequireAdmin() bool
	}

	Parser interface {
		RegisterFabric(cf CommandFabric)
		TryParse(req RequesterResponder) (Command, error)
	}

	CommandParser struct {
		fabrics []CommandFabric
	}
)

var parser Parser

func SharedParser() Parser {
	return parser
}
