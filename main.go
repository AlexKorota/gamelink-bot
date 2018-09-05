package main

import (
	"gamelinkBot/command"
	"gamelinkBot/config"
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
	reactor, err := command.NewBot(config.TBotToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	requests, err := reactor.RequesterResponderWithContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for req := range requests {
		cmd, err := command.SharedParser().TryParse(req)
		if err != nil {
			req.Respond(err.Error())
			continue
		}
		cmd.Execute(ctx)
	}
	log.Info("exiting...")
}
