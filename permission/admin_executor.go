package permission

import (
	"encoding/json"
	"errors"
	"fmt"
	"gamelinkBot/admincmd"
	"gamelinkBot/config"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type (
	//Admin - struct for admin
	Admin iface.AdminRequestStruct
	//Admins - struct for all admins
	Admins []Admin
	//AdminFileWorker - strucnt for work with MongoDB
	AdminFileWorker struct {
		admins Admins `json:"admins"`
		lock   sync.RWMutex
	}
)

//init - add AdminFileWorker(permChecker) to parser, create permfile if not exist
func init() {
	w := &AdminFileWorker{}
	if _, err := os.Stat(config.PermFile); os.IsNotExist(err) {
		w.create()
	} else {
		w.load()
	}
	parser.SharedParser().SetChecker(w)
	admincmd.SetExecutor(w)
}

//create - create new permissions file
func (afw *AdminFileWorker) create() {
	log.Print("create new file")
	jfile, err := os.Create(config.PermFile)
	if err != nil {
		log.Fatal(err)
	}
	jfile.Close()
}

//load - load data from permission file to worker struct
func (afw *AdminFileWorker) load() {
	f, err := os.OpenFile(config.PermFile, os.O_RDWR, os.ModeAppend)
	defer f.Close()
	if err != nil {
		return
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &afw.admins)
	if err != nil {
		return
	}
}

//IsAdmin - check if user is superAdmin
func (w AdminFileWorker) IsAdmin(userName string) (bool, error) {
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

//HasPermissions - check (from permission) does the user who send command have the necessary permissions
func (afw AdminFileWorker) HasPermissions(userName string, permissions []string) (bool, error) {
	afw.lock.RLock()
	defer afw.lock.RUnlock()
	log.WithFields(log.Fields{"userName": userName, "permissions": permissions}).Debug("permission.HasPermissions call")
	a := afw.admins.findAdmin(userName)
	if a == nil {
		return false, errors.New(userName + " isn't admin")
	}
	log.WithField("permissions", permissions).Debug("user")
	return a.checkPermissions(permissions), nil
}

func (a *Admin) checkPermissions(pp []string) bool {
	for _, rp := range pp {
		if !a.checkPermission(rp) {
			return false
		}
	}
	return true
}

func (a *Admin) checkPermission(p string) bool {
	for _, hp := range a.Permissions {
		if p == hp {
			return true
		}
	}
	return false
}

//GrantPermissions - update/create permissions entry for user
func (afw *AdminFileWorker) GrantPermissions(userName string, permissions []string) (*iface.AdminRequestStruct, error) {
	afw.lock.Lock()
	defer afw.lock.Unlock()
	a := afw.admins.findAdmin(userName)
	if a == nil {
		a = &Admin{Name: userName, Permissions: permissions}
		afw.admins = append(afw.admins, *a)
	} else {
		fmt.Println("before", a)
		a.grant(permissions)
		fmt.Println("after", a)
	}
	fmt.Println(afw.admins)
	afw.save()
	return (*iface.AdminRequestStruct)(a), nil
}

//RevokePermissions - revoke user permissions (delete it from permission entry)
func (afw *AdminFileWorker) RevokePermissions(userName string, permissions []string) (*iface.AdminRequestStruct, error) {
	afw.lock.Lock()
	defer afw.lock.Unlock()
	a := afw.admins.findAdmin(userName)
	if a == nil {
		return nil, errors.New(fmt.Sprintf("%s is not admin", userName))
	}
	if len(a.revokePerms(permissions)) == 0 {
		afw.admins.deleteAdmin(a)
	}
	afw.save()
	return (*iface.AdminRequestStruct)(a), nil
}

//grant - add permissions to admin
func (a *Admin) grant(permissions []string) {
	for _, v := range permissions {
		if !a.checkPermission(v) {
			a.Permissions = append(a.Permissions, v)
		}
	}
	fmt.Println(a)
}

//revoke - delete permissions, or delete admin if we delete all admins permissions
func (a *Admin) revokePerms(permissions []string) []string {
	for _, v := range permissions {
		a.revoke(v)
	}
	return a.Permissions
}

func (a *Admin) revoke(p string) []string {
	for i, v := range a.Permissions {
		if v == p {
			a.Permissions[i] = a.Permissions[len(a.Permissions)-1]
			a.Permissions = a.Permissions[:len(a.Permissions)-1]
			break
		}
	}
	return a.Permissions
}

//deleteAdmin - delete admin from json
func (a *Admins) deleteAdmin(admin *Admin) error {
	for i, v := range *a {
		if &v == admin {
			(*a)[i] = (*a)[len(*a)-1]
			(*a) = (*a)[:len(*a)-1]
			break
		}
	}
	return nil
}

//FindUser - find user entry in permission
func (a *Admins) findAdmin(userName string) *Admin {
	for _, v := range *a {
		if v.Name == userName {
			return &v
		}
	}
	return nil
}

//saveFile - save json with admins info inti json file
func (afw *AdminFileWorker) save() error {
	fmt.Println(afw.admins)
	js, err := json.Marshal(afw.admins)
	if err != nil {
		// Заменяем данные в опертаивке на начальный слепок, до редактуры
		return errors.New("marshaling error")
	}
	err = ioutil.WriteFile(config.PermFile, js, 0644)
	if err != nil {
		// Заменяем данные в опертаивке на начальный слепок, до редактуры
		log.Fatal(err)
		return err
	}
	//если все прошло успешно, заменяем слепок на новый с отредактированными данными
	return nil
}
