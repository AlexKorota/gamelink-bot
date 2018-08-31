package command

import (
//"gamelinkBot/config"
//"strings"
)

type (
	MongoChecker struct {
	}
)

func init() {
	SharedParser().SetChecker(NewMongoChecker())
}

func NewMongoChecker() *MongoChecker {
	return &MongoChecker{}
}

func (u MongoChecker) IsAdmin(userName string) (bool, error) {
	//if u.username == "" {
	//	return false
	//}
	//for _, v := range config.SuperAdmin {
	//	if u.username == strings.Trim(v, " ") {
	//		return true
	//	}
	//}
	return false, nil
}

func (u MongoChecker) HasPermissions(userName string, permissions []string) (bool, error) {
	return true, nil
}
