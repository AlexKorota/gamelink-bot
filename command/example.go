package command

import "github.com/sirupsen/logrus"

func BotExample() {
	reactor := NewBot()
	rts := reactor.RequesterResponder()
	for rt := range rts {
		logrus.Info(rt.Request())
		logrus.Info(rt.UserName())
		rt.Respond("пошел в жопу!")
	}
}
