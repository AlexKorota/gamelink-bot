package main

import (
	"fmt"
	"gamelinkBot/prot"
	"gamelinkBot/service"
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"strings"
	"time"
)

type requestStruct struct {
	params  []*prot.OneCriteriaStruct
	command string
}

func main() {
	conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := prot.NewAdminServiceClient(conn)
	if c == nil { //Но это не точно!
		log.Fatal("connection error")
	}
	telegramBot(c)
}

func telegramBot(c prot.AdminServiceClient) {
	bot, err := tgbotapi.NewBotAPI("643861723:AAHOqxU2GCQ1bqMdqycM1QPCGZEK1ekaH8s")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	var ch chan requestStruct = make(chan requestStruct)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			commands := map[string]int{
				"/start":     0,
				"/send_push": 1,
				"/count":     2,
				"/find":      3,
				"/delete":    4,
				"/update":    5,
				"/get_user":  6,
			}
			ctx := context.Background()
			arr := strings.Split(update.Message.Text, " ")
			if _, ok := commands[arr[0]]; !ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Try again")
				bot.Send(msg)
			}
			var req []*prot.OneCriteriaStruct
			if len(arr) > 1 {
				req, err = service.ParseRequest(arr[1:])
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					continue
				}
			}
			rq := requestStruct{params: req, command: arr[0]}
			fmt.Println(rq)
			ch <- rq

			go sender(ctx, ch, bot, c, update)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}

func sender(ctx context.Context, ch chan requestStruct, bot *tgbotapi.BotAPI, c prot.AdminServiceClient, update tgbotapi.Update) {
	req := <-ch
	fmt.Println("sdfsggfg")
	fmt.Println(req)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Try again")
	bot.Send(msg)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	switch req.command {
	case "/start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I'm ready to serve you Master. Send me /your_command and i'll do it for you. But remember: Be careful what you wish for!")
		bot.Send(msg)
	case "/send_push":
	case "/count":
		resp, err := c.Count(ctx, &prot.MultiCriteriaRequest{Params: req.params})
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp.String()))
	case "/find":
		resp, err := c.Find(ctx, &prot.MultiCriteriaRequest{Params: req.params})
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp.String()))
	case "/delete":
		resp, err := c.Delete(ctx, &prot.MultiCriteriaRequest{Params: req.params})
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp.String()))
	case "/update":
		resp, err := c.Update(ctx, &prot.MultiCriteriaRequest{Params: req.params})
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp.String()))
	case "/get_user":
		resp, err := c.Delete(ctx, &prot.MultiCriteriaRequest{Params: req.params})
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp.String()))
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Try again")
		bot.Send(msg)
	}
}
