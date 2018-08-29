package service

import (
	"gamelinkBot/common"
	"gamelinkBot/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func SuperAdminCheck(username string) bool {
	if username == "" {
		return false
	}
	for _, v := range config.SuperAdmin {
		if username == strings.Trim(v, " ") {
			return true
		}
	}
	return false
}

func UserPermissionsCheck(user string, collection *mgo.Collection, command string) (bool, error) {
	admin := common.AdminRequestStruct{}
	err := collection.Find(bson.M{"name": user}).One(&admin)
	if err != nil {
		return false, err
	}
	var success bool
	for _, v := range admin.Permissions {
		if v == strings.Trim(command, "/") {
			success = true
			return success, nil
		}
	}
	return false, nil
}
