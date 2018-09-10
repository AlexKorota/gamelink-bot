package generalcmd

import (
	"gamelinkBot/iface"
	"sync"
)

var (
	executor iface.RpcExecutor
	m        sync.RWMutex
)

func SetExecutor(e iface.RpcExecutor) {
	m.Lock()
	defer m.Unlock()
	executor = e
}

func Executor() iface.RpcExecutor {
	m.RLock()
	defer m.RUnlock()
	return executor
}
