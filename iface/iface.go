package iface

import (
	"context"
	msg "gamelink-go/protoMsg"
)

type (
	//Requester - interface for getting request info
	Requester interface {
		Request() string  //command name
		UserName() string //sender name
	}
	//Responder - send response message to sender
	Responder interface {
		Respond(message string) error
	}
	//Response - interface for getting response info
	Response interface {
		Response() string //text message
		ChatId() int64
	}
	//RequesterResponder - interface for RequesterResponder
	RequesterResponder interface {
		Requester
		Responder
	}
	//Reactor - interface for react to request.
	Reactor interface {
		RequesterResponderWithContext(ctx context.Context) (<-chan RequesterResponder, error) // Make new chan, get update from telegram chanel, check it(use reflect), then make RoundTrip struct and send it to our new channel.
		Respond(r Response) error                                                             //Respond to sender
	}

	//Command - interface, allows to execute command
	Command interface {
		Execute(ctx context.Context)
	}

	//CommandFabric - command fabric interface
	CommandFabric interface {
		TryParse(req RequesterResponder) (Command, error) //Try to parse request, return command if parse is success
		RequireAdmin() bool                               //Return true if command need admin permissions
		Require() []string                                //Array of needed permissions for this command
		CommandName() string                              //Returns human readable command name
	}
	//Parser - parser interface
	Parser interface {
		SetChecker(pc PermChecker)                        //add permission checker to command parser
		RegisterFabric(cf CommandFabric)                  //Register command fabric
		TryParse(req RequesterResponder) (Command, error) //Try to parse request
	}
	//PermChecker - permission checker interface
	PermChecker interface {
		IsAdmin(userName string) (bool, error)                              //Check if user who send request is admin
		HasPermissions(userName string, permissions []string) (bool, error) //Check if user who send request has required permissions
	}

	//AdminRequestStruct - admin data struct from mongoDB
	AdminRequestStruct struct {
		Name        string //UserName
		Permissions []string
	}
	//PermissionController - interface for adding and revoke users permissions
	PermissionController interface {
		GrantPermissions(userName string, permissions []string) (*AdminRequestStruct, error)
		RevokePermissions(userName string, permissions []string) (*AdminRequestStruct, error)
	}
	//AdminExecutor - interface for admin commands. Allows to check permissions, execute admin commands
	AdminExecutor interface {
		PermChecker
		PermissionController
	}
	//GeneralExecutor - interface for executing general commands
	GeneralExecutor interface {
		Count(ctx context.Context, params []*msg.OneCriteriaStruct) (*msg.CountResponse, error)
		Delete(ctx context.Context, params []*msg.OneCriteriaStruct) (*msg.OneUserResponse, error)
		Find(ctx context.Context, params []*msg.OneCriteriaStruct) (*msg.MultiUserResponse, error)
		Update(ctx context.Context, params []*msg.OneCriteriaStruct) (*msg.MultiUserResponse, error)
	}
)
