package main

import (
	"fmt"
	"gamelinkBot/common"
	"gamelinkBot/config"
	"gamelinkBot/prot"
	"gamelinkBot/service"
	"github.com/Syfaro/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"reflect"
	"strings"
	"sync"
	"time"
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
	db, err := mgo.Dial(config.MongoAddr)
	if err != nil {
		log.Fatal("can't connect to db. Error:", err)
	}
	defer db.Close()
	fmt.Println(db)
	conn, err := grpc.Dial(config.DialAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := prot.NewAdminServiceClient(conn)
	if c == nil { //Но это не точно!
		log.Fatal("connection error")
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go telegramBot(c, &wg)
	wg.Wait()
	log.Warn("Exiting...")
}

func telegramBot(c prot.AdminServiceClient, wg *sync.WaitGroup) {
	defer wg.Done()
	bot, err := tgbotapi.NewBotAPI(config.TBotToken)
	if err != nil {
		log.Error(err)
		return
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	ctxC, _ := context.WithCancel(context.Background())
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			commands := map[string]int{
				"/start":             0,
				"/send_push":         1,
				"/count":             2,
				"/find":              3,
				"/delete":            4,
				"/update":            5,
				"/get_user":          6,
				"/grant_permission":  7,
				"/revoke_permission": 8,
			}
			log.Println(update.Message.Text)
			arr := strings.Split(strings.Trim(update.Message.Text, " "), " ")
			if _, ok := commands[arr[0]]; !ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Try again")
				bot.Send(msg)
				continue
			}

			//check if user is super admin
			isSuperAdmin := service.SuperAdminCheck(update.Message.From.UserName)
			if !isSuperAdmin && (arr[0] == "/grant_permission" || arr[0] == "/revoke_permission") {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Permission denied")
				bot.Send(msg)
				continue
			}
			//if other - SuperUser middleware or Check from DB

			var req []*prot.OneCriteriaStruct
			if len(arr) > 1 {
				req, err = service.ParseRequest(arr[1:])
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					continue
				}
			}
			rq := common.RequestStruct{Params: req, Command: arr[0]}
			ctxStruct := common.ContextStruct{Request: rq, ChatID: update.Message.Chat.ID, Bot: bot, Client: c}
			ctxV := context.WithValue(ctxC, "ContextStruct", ctxStruct)
			ctxT, _ := context.WithTimeout(ctxV, time.Second*5)
			go service.Sender(ctxT)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}
