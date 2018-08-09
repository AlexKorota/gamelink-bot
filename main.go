package main

import (
	"gamelinkBot/service"
	"github.com/Syfaro/telegram-bot-api"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"strings"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	telegramBot()
}

func telegramBot() {
	bot, err := tgbotapi.NewBotAPI("643861723:AAHOqxU2GCQ1bqMdqycM1QPCGZEK1ekaH8s")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			arr := strings.Split(update.Message.Text, " ")
			if len(arr) > 1 {
				service.ParseRequest(arr[1:])
			}
			switch arr[0] {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I'm ready to serve you Master. Send me /your_command and i'll do it for you. But remember: Be careful what you wish for!")
				bot.Send(msg)
			case "/send_push":
			case "/count":
			case "/find":
			case "/delete": //DELETE
			case "/update": //JSON
			case "/get_user": //GET
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Try again")
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}
