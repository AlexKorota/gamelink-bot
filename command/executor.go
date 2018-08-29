package command

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type (
	Executor interface {
		Execute(ex Executable)
		UserName() string
	}

	RpcExecutor struct {
		userName string
	}
)

func NewExecutor(userName string) Executor {
	return &RpcExecutor{userName}
}

func (re RpcExecutor) checkPermissions(ex Executable) error {
	return errors.New(fmt.Sprintf("user %s have not permission to do %s", re.userName, ex.Command()))
}

func (re RpcExecutor) Execute(ex Executable) {
	if e := re.checkPermissions(ex); e != nil {

	}
	go func(ex Executable) {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if _, e := ex.ExecuteInContext(ctx); e == nil {

		}
	}(ex)

}

func (re RpcExecutor) UserName() string {
	return re.userName
}
