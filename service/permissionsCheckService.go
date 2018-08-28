package service

import (
	"gamelinkBot/config"
	"strings"
)

func SuperAdminCheck(username string) bool {
	for _, v := range config.SuperAdmin {
		if username == strings.Trim(v, " ") {
			return true
		}
	}
	return false
}
