package irc

import (
	"github.com/ceriath/goBlue/network"
)

var waitingChannel = make(chan struct{})

type IrcConnection struct {
	client *network.Client
}

func (ircConn *IrcConnection) Connect() {
	cli := new(network.Client)
	ircConn.client = cli
	ircConn.client.Connect("irc.twitch.tv", "6667")
}

func (ircConn *IrcConnection) Init(oauth, nick string) {
	ircConn.Sendln("PASS " + oauth)
	ircConn.Sendln("NICK " + nick)
	ircConn.Sendln("CAP REQ :twitch.tv/tags")
	ircConn.Sendln("CAP REQ :twitch.tv/commands")
	go ircConn.start()
}

func (ircConn *IrcConnection) start(){
	for {
		line, err  := ircConn.client.Recv()
		if err != nil {
			println("Error occured", err.Error())
			waitingChannel <- struct{}{}
			return
		}
		new(IrcParser).Parse(line)
	}
	waitingChannel <- struct{}{}
}

func (ircConn *IrcConnection) Sendln(line string){
	ircConn.client.Sendln(line)
}

func (ircConn *IrcConnection) Wait(){
	<- waitingChannel
}

func (ircConn *IrcConnection) Send(line, channel string){
	ircConn.client.Sendln("PRIVMSG #" + channel + " :" + line)
}

func (ircConn *IrcConnection) Join(channel string){
	ircConn.client.Sendln("JOIN #" + channel)
}

func (ircConn *IrcConnection) Leave(channel string){
	ircConn.client.Sendln("PART #" + channel)
}

func (ircConn *IrcConnection) Quit(){
	ircConn.client.Sendln("QUIT")
	ircConn.client.Close()
}

