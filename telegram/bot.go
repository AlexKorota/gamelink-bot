package telegram

import (
	"context"
	"gamelinkBot/iface"
	"github.com/Syfaro/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type (
	//Bot - struct that contains bot
	Bot struct {
		bot *tgbotapi.BotAPI
	}
	//RoundTrip - struct for round trip params
	RoundTrip struct {
		r                           iface.Reactor
		chatId                      int64
		userName, request, response string
	}
)

//NewBot - create new Reactor
func NewBot(token string) (iface.Reactor, error) {
	log.WithField("token", token).Debug("telegram.NewBot call")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Debug("error creating telegram bot api object")
		return nil, err
	}
	return &Bot{bot}, nil
}

//RequesterResponderWithContext - listen for updates from bot, then create RoundTrip and path it to channel
func (b Bot) RequesterResponderWithContext(ctx context.Context) (<-chan iface.RequesterResponder, error) {
	log.Debug("telegram.Bot.RequesterResponderWithContext call")
	if ctx.Err() != nil {
		log.Debug("context is closed already")
		return nil, ctx.Err()
	}
	rrchan := make(chan iface.RequesterResponder)
	go func(chanel chan<- iface.RequesterResponder, ctx context.Context) {
		log.Debug("telegram.Bot.RequesterResponderWithContext.goroutine call")
		if ctx.Err() != nil {
			close(rrchan)
			return
		}
		config := tgbotapi.NewUpdate(0)
		config.Timeout = 60
		updates, err := b.bot.GetUpdatesChan(config)
		if err != nil {
			log.Error("error getting update chan from telegram api obtained")
			close(rrchan)
			return
		}
		log.Debug("chan for getting original updates from telegram api obtained")
		for {
			select {
			case update := <-updates:
				log.Debug("new update arrived")
				if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
					chanel <- &RoundTrip{b, update.Message.Chat.ID,
						update.Message.From.UserName, update.Message.Text, ""}
				}
			case <-ctx.Done():
				log.Debug("context was closed")
				close(rrchan)
				return
			}
		}
		log.Debug("exiting telegram.Bot.RequesterResponderWithContext.goroutine")
	}(rrchan, ctx)
	return rrchan, nil
}

//Respond - send msg to bot
func (b Bot) Respond(r iface.Response) error {
	log.Debug("telegram.Bot.Respond call")
	if r.Response() == "" {
		return nil
	}
	_, err := b.bot.Send(tgbotapi.NewMessage(r.ChatId(), r.Response()))
	return err

}

//Request - return request string
func (rt RoundTrip) Request() string {
	log.Debug("telegram.RoundTrip.Request call")
	return rt.request
}

//UserName - return user name who send msg to bot
func (rt RoundTrip) UserName() string {
	log.Debug("telegram.RoundTrip.UserName call")
	return rt.userName
}

//ChatId - return chat id
func (rt RoundTrip) ChatId() int64 {
	log.Debug("telegram.RoundTrip.ChatId call")
	return rt.chatId
}

//Response - return response string
func (rt RoundTrip) Response() string {
	log.Debug("telegram.RoundTrip.Response call")
	return rt.response
}

func (rt RoundTrip) Respond(message string) error {
	log.WithField("message", message).Debug("telegram.RoundTrip.Respond call")
	rt.response = message
	return rt.r.Respond(rt)
}
