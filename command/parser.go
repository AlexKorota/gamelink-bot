package command

import (
	"context"
)

type (
	Command interface {
		Execute(ctx context.Context)
	}

	CommandFabric interface {
		TryParse(req RequesterResponder) Command
	}

	Parser interface {
		RegisterFabric(cf CommandFabric)
		TryParse(req RequesterResponder) Command
	}

	CommandParser struct {
		fabrics []CommandFabric
	}
)

var parser Parser

func SharedParser() Parser {
	return parser
}
