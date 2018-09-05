package main

import (
	"gamelinkBot/command"
	"gamelinkBot/config"
	"gamelinkBot/prot"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"sync"
)

func init() {
	config.LoadEnvironment()
	if config.IsDevelopmentEnv() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	command.Example()
	//db, err := mgo.Dial(config.MongoAddr)
	//if err != nil {
	//	log.Fatal("can't connect to db. Error:", err)
	//}
	//defer db.Close()
	//conn, err := grpc.Dial(config.DialAddress, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("did not connect: %s", err)
	//}
	//defer conn.Close()
	//c := prot.NewAdminServiceClient(conn)
	//if c == nil { //Но это не точно!
	//	log.Fatal("connection error")
	//}
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go telegramBot(c, &wg, db)
	//wg.Wait()
	//log.Warn("Exiting...")
}

func telegramBot(c prot.AdminServiceClient, wg *sync.WaitGroup, db *mgo.Session) {
	//defer wg.Done()
	//bot, err := tgbotapi.NewBotAPI(config.TBotToken)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//bot.Debug = false
	//log.Printf("Authorized on account %s", bot.Self.UserName)
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//updates, err := bot.GetUpdatesChan(u)
	//ctxC, _ := context.WithCancel(context.Background())
	//for update := range updates {
	//	if update.Message == nil {
	//		continue
	//	}
	//	if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
	//		commands := map[string]int{
	//			"/start":              0,
	//			"/send_push":          1,
	//			"/count":              2,
	//			"/find":               3,
	//			"/delete":             4,
	//			"/update":             5,
	//			"/get_user":           6,
	//			"/grant_permissions":  7,
	//			"/revoke_permissions": 8,
	//			"/show_permissions":   9,
	//		}
	//		arr := strings.Split(strings.Trim(update.Message.Text, " "), " ")
	//		command := arr[0]
	//		if _, ok := commands[command]; !ok {
	//			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Try again")
	//			bot.Send(msg)
	//			continue
	//		}
	//
	//		//check if user is super admin
	//		isSuperAdmin := service.SuperAdminCheck(update.Message.From.UserName)
	//		collection := db.DB(config.MongoDBName).C("admins")
	//		if command == "/grant_permissions" || command == "/revoke_permissions" || command == "/show_permissions" {
	//			if !isSuperAdmin {
	//				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Permission denied")
	//				bot.Send(msg)
	//				continue
	//			}
	//			substr := update.Message.Text[strings.Index(update.Message.Text, " ")+1:]
	//			admin, err := service.ParsePermissionRequest(substr)
	//			if err != nil {
	//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	//				continue
	//			}
	//			if command == "/grant_permissions" {
	//				selector := bson.M{"name": admin.Name}
	//				upsertdata := bson.M{"$addToSet": bson.M{"permissions": bson.M{"$each": admin.Permissions}}}
	//				_, err = collection.Upsert(selector, upsertdata)
	//			} else if command == "/revoke_permissions" {
	//				selector := bson.M{"name": admin.Name}
	//				revokePermissions := bson.M{"permissions": bson.M{"$in": admin.Permissions}}
	//				revokedata := bson.M{"$pull": revokePermissions}
	//				err = collection.Update(selector, revokedata)
	//			}
	//			if err != nil {
	//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	//				continue
	//			}
	//			user := common.AdminRequestStruct{}
	//			err = collection.Find(bson.M{"name": admin.Name}).One(&user)
	//			if err != nil {
	//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	//				continue
	//			}
	//			msg := "user " + user.Name + " has next permissions: " + strings.Join(user.Permissions, ", ")
	//			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	//			continue
	//		} else if !isSuperAdmin {
	//			success, err := service.UserPermissionsCheck(update.Message.From.UserName, collection, command)
	//			if err != nil {
	//				if err.Error() == "not found" {
	//					msg := "user " + update.Message.From.UserName + " is not admin approved to access this app"
	//					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	//					continue
	//				}
	//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	//				continue
	//			}
	//			if !success {
	//				msg := "user " + update.Message.From.UserName + " has no permission to use " + command + " command"
	//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	//				continue
	//			}
	//		}
	//
	//		var req []*prot.OneCriteriaStruct
	//		if len(arr) > 1 {
	//			req, err = service.ParseRequest(arr[1:])
	//			if err != nil {
	//				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	//				continue
	//			}
	//		}
	//		rq := common.RequestStruct{Params: req, Command: command}
	//		ctxStruct := common.ContextStruct{Request: rq, ChatID: update.Message.Chat.ID, Bot: bot, Client: c}
	//		ctxV := context.WithValue(ctxC, "ContextStruct", ctxStruct)
	//		ctxT, _ := context.WithTimeout(ctxV, time.Second*5)
	//		go service.Sender(ctxT)
	//	} else {
	//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
	//		bot.Send(msg)
	//	}
	//}
}
