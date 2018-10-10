package permission

import (
	"encoding/json"
	"errors"
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
	//PermissionWorker - strucnt for work with MongoDB
	PermissionWorker struct {
		Admins []iface.OneAdminRequestStruct `json:"admins"`
		m      sync.RWMutex
	}
)

//init - add PermissionWorker(permChecker) to parser, create permfile if not exist
func init() {
	if _, err := os.Stat(config.PermFile); os.IsNotExist(err) {
		log.Print("create new file")
		jfile, err := os.Create(config.PermFile)
		if err != nil {
			log.Fatal(err)
		}
		jfile.Close()
	}
	w, err := NewPermissionWorker()
	if err != nil {
		log.Fatal(err)
	}
	parser.SharedParser().SetChecker(w)
	admincmd.SetExecutor(w)
}

//NewPermissionWorker - make new worker with info from file
func NewPermissionWorker() (iface.AdminExecutor, error) {
	f, err := os.OpenFile(config.PermFile, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var admins []iface.OneAdminRequestStruct
	err = json.Unmarshal(bytes, &admins)
	if err != nil {
		return nil, err
	}
	return &PermissionWorker{Admins: admins}, nil
}

//IsAdmin - check if user is superAdmin
func (w PermissionWorker) IsAdmin(userName string) (bool, error) {
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
func (w PermissionWorker) HasPermissions(userName string, permissions []string) (bool, error) {
	w.m.RLock()
	defer w.m.RUnlock()
	log.WithFields(log.Fields{"userName": userName, "permissions": permissions}).Debug("permission.HasPermissions call")
	admin, _ := w.findUser(userName)
	if admin == nil {
		return false, errors.New(userName + " isn't admin")
	}
	log.WithField("permissions", permissions).Debug("user")
	for _, checkPerm := range permissions {
		successOne := false
		for _, ep := range admin.Permissions {
			if checkPerm == ep {
				successOne = true
				break
			}
		}
		if !successOne {
			return false, errors.New(userName + " has't enough permissions")
		}
	}
	return true, nil
}

//GrantPermissions - update/create permissions entry for user
func (w *PermissionWorker) GrantPermissions(userName string, permissions []string) (*iface.OneAdminRequestStruct, error) {
	w.m.Lock()
	defer w.m.Unlock()
	admin, k := w.findUser(userName)
	if admin == nil {
		newAdmin, err := w.addAdmin(userName, permissions)
		if err != nil {
			return nil, err
		}
		return newAdmin, nil
	}
	w.grant(admin, permissions)
	w.Admins[k] = *admin
	err := w.saveFile()
	if err != nil {
		return nil, err
	}
	return &w.Admins[k], nil
}

//RevokePermissions - revoke user permissions (delete it from permission entry)
func (w *PermissionWorker) RevokePermissions(userName string, permissions []string) (*iface.OneAdminRequestStruct, error) {
	w.m.Lock()
	defer w.m.Unlock()
	admin, k := w.findUser(userName)
	if admin == nil {
		return nil, errors.New(userName + " isn't admin")
	}
	perms := w.revoke(admin, permissions)
	if perms == nil {
		err := w.deleteAdmin(admin)
		return nil, err
	}
	w.Admins[k] = *admin
	err := w.saveFile()
	if err != nil {
		return nil, err
	}
	return &w.Admins[k], nil
}

//addAdmin - create new admin
func (w *PermissionWorker) addAdmin(userName string, permissions []string) (*iface.OneAdminRequestStruct, error) {
	newAdmin := iface.OneAdminRequestStruct{Name: userName, Permissions: permissions}
	w.Admins = append(w.Admins, newAdmin)
	err := w.saveFile()
	if err != nil {
		return nil, err
	}
	return &newAdmin, nil
}

//grant - add permissions to admin
func (w *PermissionWorker) grant(admin *iface.OneAdminRequestStruct, permissions []string) {
	for _, newPerm := range permissions {
		alreadyHasPerm := false
		for _, ep := range admin.Permissions {
			if newPerm == ep {
				alreadyHasPerm = true
				break
			}
		}
		if !alreadyHasPerm {
			admin.Permissions = append(admin.Permissions, newPerm)
		}
	}
}

//revoke - delete permissions, or delete admin if we delete all admins permissions
func (w *PermissionWorker) revoke(admin *iface.OneAdminRequestStruct, permissions []string) []string {
	for _, revokePerm := range permissions {
		for i, ep := range admin.Permissions {
			if revokePerm == ep {
				if len(admin.Permissions) == 1 {
					admin.Permissions = nil
					break
				} else {
					admin.Permissions = append(admin.Permissions[:i], admin.Permissions[i+1:]...)
				}
			}
		}
	}
	return admin.Permissions
}

//deleteAdmin - delete admin from json
func (w *PermissionWorker) deleteAdmin(admin *iface.OneAdminRequestStruct) error {
	if len(w.Admins) == 1 {
		w.Admins = []iface.OneAdminRequestStruct{}
	} else {
		_, k := w.findUser(admin.Name)
		w.Admins = append(w.Admins[:k], w.Admins[k+1:]...)
	}
	err := w.saveFile()
	if err != nil {
		return err
	}
	return nil
}

//FindUser - find user entry in permission
func (w PermissionWorker) findUser(userName string) (*iface.OneAdminRequestStruct, int) {
	for k, admin := range w.Admins {
		if admin.Name == userName {
			return &admin, k
		}
	}
	return nil, -1
}

//saveFile - save json with admins info inti json file
func (w *PermissionWorker) saveFile() error {
	js, err := json.Marshal(w.Admins)
	if err != nil {
		return errors.New("marshaling error")
	}
	err = ioutil.WriteFile(config.PermFile, js, 0644)
	if err != nil {
		return err
	}
	return nil
}
