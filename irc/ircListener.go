/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package irc is used for twitch's irc
package irc

import (
	"time"

	"code.cerinuts.io/libs/goBlue/archium"
	"code.cerinuts.io/libs/goBlue/log"
)

//TwitchIRCListener is a listener that matches general irc messages on twitch like PING or RECONNECT
type TwitchIRCListener struct {
	ArchiumDataIdentifier, ArchiumPrefix string
	IrcConn                              *Connection
}

//Trigger handles PING and RECONNECT
func (til *TwitchIRCListener) Trigger(ae archium.Event) {
	msg := ae.Data[til.ArchiumDataIdentifier].(*Message)
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

//GetTypes returns the twitch irc prefix for archium
func (til *TwitchIRCListener) GetTypes() []string {
	return []string{til.ArchiumPrefix + "*"}
}
