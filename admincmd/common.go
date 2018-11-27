package admincmd

import (
	"gamelink-bot/iface"
	"sync"
)

var (
	executor iface.AdminExecutor
	m        sync.RWMutex
)

func SetExecutor(e iface.AdminExecutor) {
	m.Lock()
	defer m.Unlock()
	executor = e
}

func Executor() iface.AdminExecutor {
	m.RLock()
	defer m.RUnlock()
	return executor
}
