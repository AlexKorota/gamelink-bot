package generalcmd

import (
	"gamelinkBot/iface"
	"sync"
)

var (
	executor iface.GeneralExecutor
	m        sync.RWMutex
)

func SetExecutor(e iface.GeneralExecutor) {
	m.Lock()
	defer m.Unlock()
	executor = e
}

func Executor() iface.GeneralExecutor {
	m.RLock()
	defer m.RUnlock()
	return executor
}
