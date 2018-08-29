package command

import (
	"gamelinkBot/config"
	"github.com/sirupsen/logrus"
)

func Example() {
	reactor, err := NewBot(config.TBotToken)
	if err != nil {
		logrus.Fatal(err)
	}
	rts := reactor.RequesterResponder()
	for rt := range rts {
		logrus.Info(rt.Request())
		logrus.Info(rt.UserName())
		rt.Respond("пошел в жопу!")
	}
}
