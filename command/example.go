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
		in := inputData{req.Request(), user{req.UserName()}}
		bc, err := in.MakeExecutable()
		if err != nil {
			req.Respond(err.Error())
		}

		// TODO: реализовать концепцию

		//logrus.Info(req.Request())
		//cmd, err := NewCommand(req)
		//if err != nil {
		//	req.Respond("Ху нью!")
		//}
		//cmd.ExecuteWithContext(ctx)
	}
	logrus.Info("exiting...")
}
