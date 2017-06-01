package irc

import (
	"github.com/ceriath/goBlue/network"
)

type IrcConnection struct {
	client *network.Client
}

func (ircConn *IrcConnection) Connect() {
	cli := new(network.Client)
	ircConn.client = cli
	ircConn.client.Connect("irc.twitch.tv", "6667")
}

func (ircConn *IrcConnection) Init(oauth, nick string) {
	ircConn.Send("PASS " + oauth)
	ircConn.Send("NICK " + nick)
	ircConn.Send("CAP REQ :twitch.tv/tags")
	ircConn.Send("CAP REQ :twitch.tv/commands")
	go ircConn.start()
}

func (ircConn *IrcConnection) start(){
	for true {
		line, err  := ircConn.client.Recv()
		if err != nil {
			println("Error occured", err.Error())
			return
		}
		new(IrcParser).Parse(line)
	}
}

func (ircConn *IrcConnection) Send(line string){
	ircConn.client.Sendln(line)
}
