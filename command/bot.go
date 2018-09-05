package command

import (
	"context"
	"github.com/Syfaro/telegram-bot-api"
	"reflect"
)

type (
	//Requester - interface for getting request info
	Requester interface {
		Request() string
		UserName() string
	}
	//Responder - interface for respond
	Responder interface {
		Respond(message string) error
	}
	//Response - interface for getting response info
	Response interface {
		Response() string
		ChatId() int64
	}
	//RequesterResponder - interface for RequesterResponder
	RequesterResponder interface {
		Requester
		Responder
	}
	//Reactor - interface for react to request
	Reactor interface {
		RequesterResponderWithContext(ctx context.Context) (<-chan RequesterResponder, error)
		Respond(r Response) error
	}
	//Bot - struct that contains bot
	Bot struct {
		bot *tgbotapi.BotAPI
	}
	//RoundTrip - struct for round trip params
	RoundTrip struct {
		r                           Reactor
		chatId                      int64
		userName, request, response string
	}
)

//NewBot - create new Reactor
func NewBot(token string) (Reactor, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{bot}, nil
}

//RequesterResponderWithContext - listen for updates from bot, then create RoundTrip and path it to channel
func (b Bot) RequesterResponderWithContext(ctx context.Context) (<-chan RequesterResponder, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	rrchan := make(chan RequesterResponder)
	go func(chanel chan<- RequesterResponder, ctx context.Context) {
		if ctx.Err() != nil {
			close(rrchan)
			return
		}
		config := tgbotapi.NewUpdate(0)
		config.Timeout = 60
		updates, err := b.bot.GetUpdatesChan(config)
		if err != nil {
			close(rrchan)
			return
		}
		for {
			select {
			case update := <-updates:
				if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
					chanel <- &RoundTrip{b, update.Message.Chat.ID,
						update.Message.From.UserName, update.Message.Text, ""}
				}
			case <-ctx.Done():
				close(rrchan)
				return
			}
		}
	}(rrchan, ctx)
	return rrchan, nil
}

//Respond - send msg to bot
func (b Bot) Respond(r Response) error {
	if r.Response() == "" {
		return nil
	}
	_, err := b.bot.Send(tgbotapi.NewMessage(r.ChatId(), r.Response()))
	return err

}

//Request - return request string
func (rt RoundTrip) Request() string {
	return rt.request
}

//UserName - return user name who send msg to bot
func (rt RoundTrip) UserName() string {
	return rt.userName
}

//ChatId - return chat id
func (rt RoundTrip) ChatId() int64 {
	return rt.chatId
}

//Response - return response string
func (rt RoundTrip) Response() string {
	return rt.response
}

func (rt RoundTrip) Respond(message string) error {
	rt.response = message
	return rt.r.Respond(rt)
}
