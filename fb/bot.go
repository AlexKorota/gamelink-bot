package fb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gamelinkBot/iface"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type (
	Bot struct {
		r *mux.Router
		c chan iface.RequesterResponder
	}

	RoundTrip struct {
		r                                   iface.Reactor
		chatId, userName, request, response string
	}
)

type Callback struct {
	Object string `json:"object,omitempty"`
	Entry  []struct {
		ID        string      `json:"id,omitempty"`
		Time      int         `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

type Messaging struct {
	Sender    User    `json:"sender,omitempty"`
	Recipient User    `json:"recipient,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Message struct {
	MID        string `json:"mid,omitempty"`
	Text       string `json:"text,omitempty"`
	QuickReply *struct {
		Payload string `json:"payload,omitempty"`
	} `json:"quick_reply,omitempty"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
	Attachment  *Attachment   `json:"attachment,omitempty"`
}

type Attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Response struct {
	Recipient User    `json:"recipient,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type Payload struct {
	URL string `json:"url,omitempty"`
}

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
		fmt.Println("new goruitine")
		b.r = mux.NewRouter()
		b.r.HandleFunc("/webhook", b.verificationEndpoint).Methods("GET")
		b.r.HandleFunc("/webhook", b.messagesEndpoint).Methods("POST")
		err := http.ListenAndServe("localhost:8088", b.r)
		if err != nil {
			log.Fatal(err)
			close(rrchan)
			return
		}
		return
	}()
	return rrchan, nil
}

func (b Bot) verificationEndpoint(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	fmt.Println("mode", mode)
	fmt.Println("token", token)
	if mode != "" && token == "qqqwww" {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		fmt.Println(mode)
		fmt.Println("t", os.Getenv("VERIFY_TOKEN"))
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token!!!!"))
	}
}

func (b Bot) messagesEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("message endpoint")
	var callback Callback
	json.NewDecoder(r.Body).Decode(&callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				b.c <- &RoundTrip{b, event.Sender.ID, event.Sender.ID, event.Message.Text, ""}
			}
		}
		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
	}
}

func (b Bot) Respond(r iface.Response) error {
	fmt.Println("resp", r)
	if r.Response() == "" {
		return nil
	}
	client := &http.Client{}
	response := Response{
		Recipient: User{
			ID: r.ChatId(),
		},
		Message: Message{
			Text: r.Response(),
		},
	}
	fmt.Println("resp", response)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=EAADgsjzTwy4BABdmGjqQHrHcmVkrq9QQpZCvRGGQXl0tjxn1wDmv5z6i7CuZAqNxv6u4MdL50INZBPSCR9c4XwEdKM4ZCMk3JZAREM8S9por325fZAEZANz59iclrPhDZBytaeaV46gf67OKp6dIJiQSQqNt994kCpGlI5NSUM7m7QZDZD"
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
	return rt.r.Respond(rt)
}
