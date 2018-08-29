package command

import "github.com/Syfaro/telegram-bot-api"

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
	return nil
}

func (b Bot) Respond(r Response) {

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
	rt.r.Respond(rt)
	return
}
