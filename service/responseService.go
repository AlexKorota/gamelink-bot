package service

import (
	"gamelinkBot/common"
	"github.com/Syfaro/telegram-bot-api"
)

func ResponseHandler(respMsg string, err error, ctx common.ContextStruct) {
	if err != nil {
		ctx.Bot.Send(tgbotapi.NewMessage(ctx.ChatID, err.Error()))
	}
	ctx.Bot.Send(tgbotapi.NewMessage(ctx.ChatID, respMsg))
}
