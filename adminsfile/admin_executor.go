package adminsfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"gamelinkBot/admincmd"
	"gamelinkBot/config"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type (
	Admin           iface.AdminRequestStruct
	Admins          []Admin
	AdminFileWorker struct {
		admins Admins
		lock   sync.RWMutex
	}
)

func init() {
	w := &AdminFileWorker{}
	w.Load()
	parser.SharedParser().SetChecker(w)
	admincmd.SetExecutor(w)
}

func (a *Admin) grantPermission(p []string) {
	for _, v := range p {
		if !a.checkPermission(v) {
			a.Permissions = append(a.Permissions, p...)
		}
	}
}

func (a *Admin) revokePermission(p string) []string {
	for i, v := range a.Permissions {
		if v == p {
			a.Permissions[i] = a.Permissions[len(a.Permissions)-1]
			a.Permissions = a.Permissions[:len(a.Permissions)-1]
			break
		}
	}
	return a.Permissions
}

func (a *Admin) revokePermissions(p []string) []string {
	for _, v := range p {
		a.revokePermission(v)
	}
	return a.Permissions
}

func (a *Admin) checkPermission(p string) bool {
	for _, hp := range a.Permissions {
		if p == hp {
			return true
		}
	}
	return false
}

func (a *Admin) checkPermissions(pp []string) bool {
	for _, rp := range pp {
		if !a.checkPermission(rp) {
			return false
		}
	}
	return true
}

func (a *Admins) Find(name string) *Admin {
	for _, v := range *a {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func (a *Admins) Delete(admin *Admin) {
	for i, v := range *a {
		if &v == admin {
			(*a)[i] = (*a)[len(*a)-1]
			(*a) = (*a)[:len(*a)-1]
			break
		}
	}
}

func (afw *AdminFileWorker) IsAdmin(userName string) (bool, error) {
	if userName == "" {
		return false, nil
	}
	for _, v := range config.SuperAdmin {
		if userName == strings.Trim(v, " ") {
			return true, nil
		}
	}
	return false, nil
}

func (afw *AdminFileWorker) HasPermissions(userName string, permissions []string) (bool, error) {
	afw.lock.RLock()
	defer afw.lock.RUnlock()
	a := afw.admins.Find(userName)
	if a == nil {
		return false, errors.New(fmt.Sprintf("%s is not admin", userName))
	}
	return a.checkPermission(permissions), nil
}

func (afw *AdminFileWorker) GrantPermissions(userName string, permissions []string) (*iface.AdminRequestStruct, error) {
	afw.lock.Lock()
	defer afw.lock.Unlock()
	a := afw.admins.Find(userName)
	if a == nil {
		a = &Admin{userName, permissions}
		afw.admins = append(afw.admins, *a)
	} else {
		a.grantPermission(permissions)
	}
	afw.Save()
	return (*iface.AdminRequestStruct)(a), nil
}

func (afw *AdminFileWorker) RevokePermissions(userName string, permissions []string) (*iface.AdminRequestStruct, error) {
	afw.lock.Lock()
	defer afw.lock.Unlock()
	a := afw.admins.Find(userName)
	if a == nil {
		return nil, errors.New(fmt.Sprintf("%s is not admin", userName))
	}
	if len(a.revokePermissions(permissions)) == 0 {
		afw.admins.Delete(a)
	}
	afw.Save()
	return (*iface.AdminRequestStruct)(a), nil
}

func (afw *AdminFileWorker) Load() {
	f, err := os.OpenFile("файло", os.O_RDWR, os.ModeAppend)
	if err != nil {
		return
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	json.Unmarshal(bytes, &afw.admins)
}

func (afw *AdminFileWorker) Save() {
	js, err := json.Marshal(afw.admins)
	if err != nil {
		return
	}
	ioutil.WriteFile("файло", js, 0644)
}
