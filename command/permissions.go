package command

import (
//"gamelinkBot/config"
//"strings"
)

type (
	Checker struct {
	}
)

func init() {
	SharedParser().SetChecker(&Checker{})
}

func (u Checker) IsAdmin(userName string) (bool, error) {
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

func (u Checker) HasPermissions(userName string, permissions []string) (bool, error) {
	return true, nil
}
