package command

import (
	"context"
	"gamelinkBot/config"
	"github.com/sirupsen/logrus"
)

func Example() {
	reactor, err := NewBot(config.TBotToken)
	if err != nil {
		logrus.Fatal(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	requests, err := reactor.RequesterResponderWithContext(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	for req := range requests {
		cmd, err := SharedParser().TryParse(req)
		if err != nil {
			req.Respond(err.Error())
			continue
		}
		cmd.Execute(ctx)
	}
	logrus.Info("exiting...")
}
