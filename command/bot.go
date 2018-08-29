package command

import (
	"github.com/Syfaro/telegram-bot-api"
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
		RequesterResponder() chan<- RequesterResponder
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

func NewBot() Reactor {
	return &Bot{}
}

func (b Bot) RequesterResponder() chan<- RequesterResponder {
	rrchan := make(chan<- RequesterResponder)
	go func(chanel *chan<- RequesterResponder) {
		config := tgbotapi.NewUpdate(0)
		config.Timeout = 60
		updates, err := b.bot.GetUpdatesChan(config)
		if err != nil {
			return
		}
		for update := range updates {
			*chanel <- &RoundTrip{b, update.Message.Chat.ID,
				update.Message.From.UserName, update.Message.Text, ""}
		}
	}(&rrchan)
	return rrchan
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
