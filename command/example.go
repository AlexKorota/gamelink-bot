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
	ctx, done := context.WithCancel(context.Background())
	requests, err := reactor.RequesterResponderWithContext(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	for req := range requests {
		logrus.Info(req.Request())
		//cmd, err := NewCommand(req)
		//if err != nil {
		//	req.Respond("Ху нью!")
		//}
		//cmd.ExecuteWithContext(ctx)
		if req.Request() == "/done" {
			done()
		} else {
			req.Respond("sfidbigbigbs")
		}
	}
	logrus.Info("exiting...")
}
