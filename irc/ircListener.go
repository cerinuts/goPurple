package irc

import (
	"time"

	"gitlab.ceriath.net/libs/goBlue/archium"
	"gitlab.ceriath.net/libs/goBlue/log"
)

type TwitchIRCListener struct {
	ArchiumDataIdentifier, ArchiumPrefix string
	IrcConn                              *IrcConnection
}

func (til *TwitchIRCListener) Trigger(ae archium.ArchiumEvent) {
	msg := ae.Data[til.ArchiumDataIdentifier].(*IrcMessage)
	til.IrcConn.lastActivity = time.Now()
	if msg.Command == "PING" {
		(*(til.IrcConn)).Sendln("PONG " + msg.Msg)
	}
	if msg.Command == "001" {
		(*(til.IrcConn)).currentReconnectAttempts = 0
	}
	if msg.Command == "RECONNECT" {
		(*(til.IrcConn)).Reconnect()
		log.I("Forced Reconnect...")
		//Might aswell be handled by general reconnect recovery. I'll leave this here until confirmed what happens exactly.
	}
}

func (til *TwitchIRCListener) GetTypes() []string {
	return []string{til.ArchiumPrefix + "*"}
}
