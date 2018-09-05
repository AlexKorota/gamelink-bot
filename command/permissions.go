package command

import (
	"errors"
	"gamelinkBot/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

type (
	MongoWorker struct {
		db *mgo.Session
	}

	AdminRequestStruct struct {
		Name        string
		Permissions []string
	}
)

func init() {
	SharedParser().SetChecker(NewMongoWorker())
}

func NewMongoWorker() *MongoWorker {
	db, err := mgo.Dial(config.MongoAddr)
	if err != nil {
		log.Fatal("can't connect to db. Error:", err)
	}
	return &MongoWorker{db: db}
}

func (u MongoWorker) IsAdmin(userName string) (bool, error) {
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

func (u MongoWorker) HasPermissions(userName string, permissions []string) (bool, error) {
	user := AdminRequestStruct{}
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

func (u MongoWorker) GrantPermissions(userName string, permissions []string) (*AdminRequestStruct, error) {
	selector := bson.M{"name": userName}
	upsertdata := bson.M{"$addToSet": bson.M{"permissions": bson.M{"$each": permissions}}}
	_, err := u.db.DB(config.MongoDBName).C("admins").Upsert(selector, upsertdata)
	if err != nil {
		return nil, err
	}
	return u.FindUser(userName)
}

func (u MongoWorker) RevokePermissions(userName string, permissions []string) (*AdminRequestStruct, error) {
	selector := bson.M{"name": userName}
	revokePermissions := bson.M{"permissions": bson.M{"$in": permissions}}
	revokedata := bson.M{"$pull": revokePermissions}
	err := u.db.DB(config.MongoDBName).C("admins").Update(selector, revokedata)
	if err != nil {
		return nil, err
	}
	return u.FindUser(userName)
}

func (u MongoWorker) FindUser(userName string) (*AdminRequestStruct, error) {
	user := AdminRequestStruct{}
	err := u.db.DB(config.MongoDBName).C("admins").Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		if err.Error() == "not found" {
			return nil, errors.New("user " + userName + " is not admin approved to access this app")
		}
		return nil, err
	}
	return &user, nil
}
