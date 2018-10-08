package mongo

import (
	"errors"
	"gamelinkBot/admincmd"
	"gamelinkBot/config"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type (
	//MongoWorker - strucnt for work with MongoDB
	MongoWorker struct {
		db *mgo.Session
	}
)

//init - add MongoWorker(permChecker) to parser
func init() {
	w := NewMongoWorker()
	parser.SharedParser().SetChecker(w)
	admincmd.SetExecutor(w)
}

//NewMongoWorker - set connection to mongoDB
func NewMongoWorker() iface.AdminExecutor {
	db, err := mgo.Dial(config.MongoAddr)
	if err != nil {
		log.Fatal("can't connect to db. Error:", err)
	}
	return &MongoWorker{db: db}
}

//IsAdmin - check if user is superAdmin
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

//HasPermissions - check (from mongo) does the user who send command have the necessary permissions
func (u MongoWorker) HasPermissions(userName string, permissions []string) (bool, error) {
	log.WithFields(log.Fields{"userName": userName, "permissions": permissions}).Debug("mongo.HasPermissions call")
	user := iface.AdminRequestStruct{}
	err := u.db.DB(config.MongoDBName).C("admins").Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		if err.Error() == "not found" {
			return false, errors.New("user " + userName + " is not admin approved to access this app")
		}
		return false, err
	}
	log.WithField("permissions", user.Permissions).Debug("user")
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

//GrantPermissions - update/create permissions entry for user (in MongoDB)
func (u MongoWorker) GrantPermissions(userName string, permissions []string) (*iface.AdminRequestStruct, error) {
	selector := bson.M{"name": userName}
	upsertdata := bson.M{"$addToSet": bson.M{"permissions": bson.M{"$each": permissions}}}
	_, err := u.db.DB(config.MongoDBName).C("admins").Upsert(selector, upsertdata)
	if err != nil {
		return nil, err
	}
	return u.findUser(userName)
}

//RevokePermissions - revoke user permissions (delete it from mongo entry)
func (u MongoWorker) RevokePermissions(userName string, permissions []string) (*iface.AdminRequestStruct, error) {
	selector := bson.M{"name": userName}
	revokePermissions := bson.M{"permissions": bson.M{"$in": permissions}}
	revokedata := bson.M{"$pull": revokePermissions}
	err := u.db.DB(config.MongoDBName).C("admins").Update(selector, revokedata)
	if err != nil {
		return nil, err
	}
	return u.findUser(userName)
}

//FindUser - find user entry in mongo
func (u MongoWorker) findUser(userName string) (*iface.AdminRequestStruct, error) {
	user := iface.AdminRequestStruct{}
	err := u.db.DB(config.MongoDBName).C("admins").Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		if err.Error() == "not found" {
			return nil, errors.New("user " + userName + " is not admin approved to access this app")
		}
		return nil, err
	}
	return &user, nil
}
