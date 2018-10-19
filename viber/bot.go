package viber

import (
	"context"
	"fmt"
	"gamelinkBot/iface"
	"github.com/mileusna/viber"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type (
	Bot struct {
		v *viber.Viber
		c chan iface.RequesterResponder
	}
	RoundTrip struct {
		r                                   iface.Reactor
		chatId, userName, request, response string
	}
)

func NewBot() iface.Reactor {
	return &Bot{}
}

func (b Bot) RequesterResponderWithContext(ctx context.Context) (<-chan iface.RequesterResponder, error) {
	if ctx.Err() != nil {
		log.Debug("context is closed already")
		return nil, ctx.Err()
	}
	rrchan := make(chan iface.RequesterResponder)
	b.c = rrchan
	go func() {
		b.v = &viber.Viber{
			AppKey: "48a1d03bf2a7d4d6-a2dad05a4ffdae0c-851813c42a88b053",
			Sender: viber.Sender{
				Name: "AKBotTest",
			},
			Message:   b.myMsgReceivedFunc, // your function for handling messages
			Delivered: b.myDeliveredFunc,
		}
		//http.Handle("/", b.v)
		err := http.ListenAndServe("localhost:8088", b.v)
		if err != nil {
			log.Fatal(err)
			close(rrchan)
			return
		}
		return
	}()
	return rrchan, nil
}

func (b Bot) myMsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	log.Warn("receive func")
	b.v = v // ----- почему тут b.v nil и приходмится еще раз присваивать?
	switch m.(type) {
	case *viber.TextMessage:
		msg := m.(*viber.TextMessage).Text
		b.c <- &RoundTrip{b, u.ID, u.Name, msg, ""}
	}
}

func (b Bot) myDeliveredFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "delivered to user ID", userID)
}

func (b Bot) Respond(r iface.Response) error {
	if r.Response() == "" {
		return nil
	}
	fmt.Println("b.v", b.v)
	fmt.Println("chat id", r.ChatId())
	fmt.Println("bresponsw", r.Response())
	b.v.SendTextMessage(r.ChatId(), r.Response())
	return nil
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
func (rt RoundTrip) ChatId() string {
	return rt.chatId
}

//Response - return response string
func (rt RoundTrip) Response() string {
	return rt.response
}

func (rt RoundTrip) Respond(message string) error {
	log.WithField("message", message).Debug("telegram.RoundTrip.Respond call")
	rt.response = message
	fmt.Println("bbbbbooooootttt", rt.r)
	return rt.r.Respond(rt)
}
