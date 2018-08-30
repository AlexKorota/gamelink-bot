package command

import (
	"gamelinkBot/config"
	"strings"
)

type (
	user struct {
		username string
	}

	PermChecker interface {
		IsAdmin() bool
		HasPermission() (bool, error)
		UserName() string
	}
)

func (u user) IsAdmin() bool {
	if u.UserName() == "" {
		return false
	}
	for _, v := range config.SuperAdmin {
		if u.UserName() == strings.Trim(v, " ") {
			return true
		}
	}
	return false
}

func (u user) HasPermission() (bool, error) {
	return false, nil
}

func (u user) UserName() string {
	return u.username
}
