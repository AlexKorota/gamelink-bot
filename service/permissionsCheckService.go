package service

import (
	"gamelinkBot/config"
	"strings"
)

func SuperAdminCheck(username string) bool {
	arr := strings.Split(config.SuperAdmin, ",")
	for _, v := range arr {
		if username == strings.Trim(v, " ") {
			return true
		}
	}
	return false
}
