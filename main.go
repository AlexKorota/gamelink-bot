package main

import (
	"context"
	_ "gamelinkBot/admincmd"
	"gamelinkBot/config"
	_ "gamelinkBot/generalcmd"
	"gamelinkBot/parser"
	_ "gamelinkBot/permission"
	_ "gamelinkBot/rpc"
	"gamelinkBot/viber"
	log "github.com/sirupsen/logrus"
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

	log.Warn("starting...")
	var err error
	reactor := viber.NewBot() // ---------------------- viber bot
	//reactor := fb.NewBot() //-------------------------- fb bot
	//reactor, err := telegram.NewBot(config.TBotToken) // telegram bot
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
