package main

import (
	_ "gamelinkBot/admincmd"
	"gamelinkBot/config"
	_ "gamelinkBot/generalcmd"
	_ "gamelinkBot/mongo"
	"gamelinkBot/parser"
	_ "gamelinkBot/rpc"
	"gamelinkBot/telegram"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
)

func init() {
	config.LoadEnvironment()
	if config.IsDevelopmentEnv() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	l, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	log.SetOutput(l)
	reactor, err := telegram.NewBot(config.TBotToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	requests, err := reactor.RequesterResponderWithContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for req := range requests {
		log.Info("user: " + req.UserName() + " command: " + req.Request())
		cmd, err := parser.SharedParser().TryParse(req)
		if err != nil {
			req.Respond(err.Error())
			continue
		}
		cmd.Execute(ctx)
	}
	log.Info("exiting...")
}
