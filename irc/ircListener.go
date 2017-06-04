package irc

import (
	"github.com/ceriath/goBlue/archium"
	"github.com/ceriath/goBlue/log"
)



type TwitchIRCListener struct{
	ArchiumDataIdentifier, ArchiumPrefix string
	IrcConn *IrcConnection
}

func (til *TwitchIRCListener) Trigger(ae archium.ArchiumEvent){
	msg := ae.Data[til.ArchiumDataIdentifier].(*IrcMessage)
	if(msg.Command == "PING"){
		(*(til.IrcConn)).Sendln("PONG " + msg.Msg)
		log.I("PONG" + msg.Msg)
	}
	if(msg.Command == "RECONNECT"){
		(*(til.IrcConn)).Reconnect()
		log.I("Forced Reconnect...")
	}
}
	
func (til *TwitchIRCListener) GetTypes() []string{
	return []string{til.ArchiumPrefix + "*"}
}

