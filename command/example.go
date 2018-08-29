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
	requests, _ := reactor.RequesterResponderWithContext(ctx)
	for req := range requests {
		//cmd, err := NewCommand(req)
		//if err != nil {
		//	req.Respond("Ху нью!")
		//}
		//cmd.ExecuteWithContext(ctx)
		req.Respond("xyz")
	}
}
