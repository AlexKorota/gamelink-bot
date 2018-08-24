package main

import (
	"gamelinkBot/config"
	"gamelinkBot/prot"
	"gamelinkBot/service"
	"github.com/Syfaro/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"reflect"
	"strings"
	"time"
)

type requestStruct struct {
	params  []*prot.OneCriteriaStruct
	command string
}

type contextStruct struct {
	request requestStruct
	chatID  int64
	bot     *tgbotapi.BotAPI
	client  prot.AdminServiceClient
}

func init() {
	config.LoadEnvironment()
	if config.IsDevelopmentEnv() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	conn, err := grpc.Dial(config.DialAddress, grpc.WithInsecure())
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
	bot, err := tgbotapi.NewBotAPI(config.TBotToken)
	if err != nil {
		log.Panic(err)
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
				"/start":     0,
				"/send_push": 1,
				"/count":     2,
				"/find":      3,
				"/delete":    4,
				"/update":    5,
				"/get_user":  6,
			}
			arr := strings.Split(strings.Trim(update.Message.Text, " "), " ")
			if _, ok := commands[arr[0]]; !ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command. Try again")
				bot.Send(msg)
				continue
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
			ctxStruct := contextStruct{request: rq, chatID: update.Message.Chat.ID, bot: bot, client: c}
			ctxV := context.WithValue(ctxC, "contextStruct", ctxStruct)
			ctxT, _ := context.WithTimeout(ctxV, time.Second*5)
			go sender(ctxT)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}

func sender(ctx context.Context) {
	ctxData := ctx.Value("contextStruct").(contextStruct)
	switch ctxData.request.command {
	case "/start":
		msg := tgbotapi.NewMessage(ctxData.chatID, "I'm ready to serve you Master. Send me /your_command and i'll do it for you. But remember: Be careful what you wish for!")
		ctxData.bot.Send(msg)
	case "/send_push":
	case "/count":
		resp, err := ctxData.client.Count(ctx, &prot.MultiCriteriaRequest{Params: ctxData.request.params})
		if err != nil {
			ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, err.Error()))
		}
		ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, resp.String()))
	case "/find":
		resp, err := ctxData.client.Find(ctx, &prot.MultiCriteriaRequest{Params: ctxData.request.params})
		if err != nil {
			ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, err.Error()))
		}
		ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, resp.String()))
	case "/delete":
		resp, err := ctxData.client.Delete(ctx, &prot.MultiCriteriaRequest{Params: ctxData.request.params})
		if err != nil {
			ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, err.Error()))
		}
		ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, resp.String()))
	case "/update":
		resp, err := ctxData.client.Update(ctx, &prot.MultiCriteriaRequest{Params: ctxData.request.params})
		if err != nil {
			ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, err.Error()))
		}
		ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, resp.String()))
	case "/get_user":
		resp, err := ctxData.client.Delete(ctx, &prot.MultiCriteriaRequest{Params: ctxData.request.params})
		if err != nil {
			ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, err.Error()))
		}
		ctxData.bot.Send(tgbotapi.NewMessage(ctxData.chatID, resp.String()))
	default:
		msg := tgbotapi.NewMessage(ctxData.chatID, "Invalid command. Try again")
		ctxData.bot.Send(msg)
	}
}
