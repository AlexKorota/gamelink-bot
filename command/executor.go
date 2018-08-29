package command

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type (
	Executor interface {
		Execute(ex Executable) (string, error)
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

func (re RpcExecutor) Execute(ex Executable) (r string, e error) {
	if e = re.checkPermissions(ex); e != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if str, e := ex.ExecuteInContext(ctx); e == nil {
		r = str.String()
	}
	return
}

func (re RpcExecutor) UserName() string {
	return re.userName
}
