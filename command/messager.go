package command

import (
	"github.com/Syfaro/telegram-bot-api"
)

type (
	Sender interface {
		Send(message string) error
	}

	Chat struct {
		bot    *tgbotapi.BotAPI
		chatId int64
	}
)

func NewChat(bot *tgbotapi.BotAPI, chatId int64) Sender {
	return &Chat{bot, chatId}
}

func (c Chat) Send(message string) (e error) {
	_, e = c.bot.Send(tgbotapi.NewMessage(c.chatId, message))
	return
}
