package service

import (
	"fmt"
	"gamelinkBot/common"
	"gamelinkBot/prot"
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
)

func Sender(ctx context.Context) {
	ctxData := ctx.Value("ContextStruct").(common.ContextStruct)
	var r fmt.Stringer
	var e error
	var msg tgbotapi.Chattable
	switch ctxData.Request.Command {
	case "/start":
		msg = tgbotapi.NewMessage(ctxData.ChatID, "I'm ready to serve you Master. Send me /your_command and i'll do it for you. But remember: Be careful what you wish for!")
	case "/send_push":
	case "/count":
		r, e = ctxData.Client.Count(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
	case "/find":
		r, e = ctxData.Client.Find(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
	case "/delete":
		r, e = ctxData.Client.Delete(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
	case "/update":
		r, e = ctxData.Client.Update(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
	case "/get_user":
		r, e = ctxData.Client.Delete(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
	default:
		msg = tgbotapi.NewMessage(ctxData.ChatID, "Invalid command. Try again")

	}
	if msg == nil {
		if err != nil {
			msg = tgbotapi.NewMessage(ctxData.ChatID, e.Error())
		} else if r != nil {
			msg = tgbotapi.NewMessage(ctxData.ChatID, r.String())
		}
	}
	ctxData.Bot.Send(msg)
}
