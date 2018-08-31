package command

import (
	"errors"
	"gamelinkBot/common"
	"gamelinkBot/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

type (
	MongoChecker struct {
		db *mgo.Session
	}
)

func init() {
	SharedParser().SetChecker(NewMongoChecker())
}

func NewMongoChecker() *MongoChecker {
	db, err := mgo.Dial(config.MongoAddr)
	if err != nil {
		log.Fatal("can't connect to db. Error:", err)
	}
	return &MongoChecker{db: db}
}

func (u MongoChecker) IsAdmin(userName string) (bool, error) {
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

func (u MongoChecker) HasPermissions(userName string, permissions []string) (bool, error) {
	user := common.AdminRequestStruct{}
	err := u.db.DB(config.MongoDBName).C("admins").Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		if err.Error() == "not found" {
			return false, errors.New("user " + userName + " is not admin approved to access this app")
		}
		return false, err
	}
	for _, cmdp := range permissions {
		successOne := false
		for _, up := range user.Permissions {
			if up == cmdp {
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
