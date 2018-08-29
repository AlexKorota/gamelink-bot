package command

import (
	"context"
	"github.com/Syfaro/telegram-bot-api"
	"reflect"
)

type (
	Requester interface {
		Request() string
		UserName() string
	}

	Responder interface {
		Respond(message string) error
	}

	Response interface {
		Response() string
		ChatId() int64
	}

	RequesterResponder interface {
		Requester
		Responder
	}

	Reactor interface {
		RequesterResponderWithContext(ctx context.Context) (<-chan RequesterResponder, error)
		Respond(r Response)
	}

	Bot struct {
		bot *tgbotapi.BotAPI
	}

	RoundTrip struct {
		r                           Reactor
		chatId                      int64
		userName, request, response string
	}
)

func NewBot(token string) (Reactor, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{bot}, nil
}

func (b Bot) RequesterResponderWithContext(ctx context.Context) (<-chan RequesterResponder, error) {
	rrchan := make(chan RequesterResponder)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
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

func (b Bot) Respond(r Response) {
	b.bot.Send(tgbotapi.NewMessage(r.ChatId(), r.Response()))
}

func (rt RoundTrip) Request() string {
	return rt.request
}

func (rt RoundTrip) UserName() string {
	return rt.userName
}

func (rt RoundTrip) ChatId() int64 {
	return rt.chatId
}

func (rt RoundTrip) Response() string {
	return rt.response
}

func (rt RoundTrip) Respond(message string) (e error) {
	rt.response = message
	rt.r.Respond(rt)
	return
}
