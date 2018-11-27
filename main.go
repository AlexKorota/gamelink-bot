package main

import (
	_ "gamelink-bot/admincmd"
	"gamelink-bot/config"
	_ "gamelink-bot/generalcmd"
	"gamelink-bot/parser"
	_ "gamelink-bot/permission"
	_ "gamelink-bot/rpc"
	"gamelink-bot/telegram"
	"gamelink-bot/version"
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
	log.Printf(
		"Starting the service: commit: %s, build time: %s, release: %s",
		version.Commit, version.BuildTime, version.Release,
	)
	reactor, err := telegram.NewBot(config.TBotToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("reactor initialized")
	ctx, _ := context.WithCancel(context.Background())
	requests, err := reactor.RequesterResponderWithContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("reactor's event emitter started")
	for req := range requests {
		log.WithFields(log.Fields{"from": req.UserName(), "request": req.Request()}).Warn("new request arrived")
		cmd, err := parser.SharedParser().TryParse(req)
		if err != nil {
			log.WithError(err).Warn("error while parsing request")
			req.Respond(err.Error())
			continue
		}
		cmd.Execute(ctx)
	}
	log.Info("exiting...")
}
