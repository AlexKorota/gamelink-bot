package main

import (
	"gamelinkBot/bot"
	_ "gamelinkBot/command"
	"gamelinkBot/config"
	"gamelinkBot/parser"
	_ "gamelinkBot/permission"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
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
	reactor, err := bot.NewBot(config.TBotToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	requests, err := reactor.RequesterResponderWithContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for req := range requests {
		cmd, err := parser.SharedParser().TryParse(req)
		if err != nil {
			req.Respond(err.Error())
			continue
		}
		cmd.Execute(ctx)
	}
	log.Info("exiting...")
}
