package iface

import "context"

type (
	//Requester - interface for getting request info
	Requester interface {
		Request() string
		UserName() string
	}
	//Responder - interface for respond
	Responder interface {
		Respond(message string) error
	}
	//Response - interface for getting response info
	Response interface {
		Response() string
		ChatId() int64
	}
	//RequesterResponder - interface for RequesterResponder
	RequesterResponder interface {
		Requester
		Responder
	}
	//Reactor - interface for react to request
	Reactor interface {
		RequesterResponderWithContext(ctx context.Context) (<-chan RequesterResponder, error)
		Respond(r Response) error
	}

	//Command - command interface
	Command interface {
		Execute(ctx context.Context)
	}

	//CommandFabric - command fabric interface
	CommandFabric interface {
		TryParse(req RequesterResponder) (Command, error)
		RequireAdmin() bool
		Require() []string
	}
	//Parser - parser inerface
	Parser interface {
		SetChecker(pc PermChecker)
		RegisterFabric(cf CommandFabric)
		TryParse(req RequesterResponder) (Command, error)
	}
	//PermChecker - permission checker interface
	PermChecker interface {
		IsAdmin(userName string) (bool, error)
		HasPermissions(userName string, permissions []string) (bool, error)
	}

	//AdminRequestStruc - admin data struct from mongoDB
	AdminRequestStruct struct {
		Name        string
		Permissions []string
	}

	PermissionController interface {
		GrantPermissions(userName string, permissions []string) (*AdminRequestStruct, error)
		RevokePermissions(userName string, permissions []string) (*AdminRequestStruct, error)
	}

	AdminExecutor interface {
		PermChecker
		PermissionController
	}
)
