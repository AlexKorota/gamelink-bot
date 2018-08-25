package service

import (
	"gamelinkBot/common"
	"gamelinkBot/prot"
	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/context"
)

func Sender(ctx context.Context) {
	ctxData := ctx.Value("ContextStruct").(common.ContextStruct)
	switch ctxData.Request.Command {
	case "/start":
		msg := tgbotapi.NewMessage(ctxData.ChatID, "I'm ready to serve you Master. Send me /your_command and i'll do it for you. But remember: Be careful what you wish for!")
		ctxData.Bot.Send(msg)
	case "/send_push":
	case "/count":
		resp, err := ctxData.Client.Count(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
		ResponseHandler(resp.String(), err, ctxData)
	case "/find":
		resp, err := ctxData.Client.Find(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
		ResponseHandler(resp.String(), err, ctxData)
	case "/delete":
		resp, err := ctxData.Client.Delete(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
		ResponseHandler(resp.String(), err, ctxData)
	case "/update":
		resp, err := ctxData.Client.Update(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
		ResponseHandler(resp.String(), err, ctxData)
	case "/get_user":
		resp, err := ctxData.Client.Delete(ctx, &prot.MultiCriteriaRequest{Params: ctxData.Request.Params})
		ResponseHandler(resp.String(), err, ctxData)
	default:
		msg := tgbotapi.NewMessage(ctxData.ChatID, "Invalid command. Try again")
		ctxData.Bot.Send(msg)
	}
}
