package common

import (
	"gamelinkBot/prot"
	"github.com/Syfaro/telegram-bot-api"
)

type RequestStruct struct {
	Params  []*prot.OneCriteriaStruct
	Command string
}

type ContextStruct struct {
	Request RequestStruct
	ChatID  int64
	Bot     *tgbotapi.BotAPI
	Client  prot.AdminServiceClient
}
