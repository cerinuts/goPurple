package irc

import (
	"github.com/ceriath/goBlue/archium"
	"github.com/ceriath/goBlue/log"
	"github.com/ceriath/goBlue/network"
	"strings"
)

var waitingChannel = make(chan struct{})
var archiumCore = archium.ArchiumCore

var ArchiumPrefix = "twitch.irc."
var ArchiumDataIdentifier = "Message"

type IrcConnection struct {
	client          *network.Client
	oauth, username, ip, port string
}

func (ircConn *IrcConnection) Connect(ip, port string) error{
	cli := new(network.Client)
	ircConn.ip = ip
	ircConn.port = port
	ircConn.client = cli
	err := ircConn.client.Connect(ip, port)
	return err
}

func (ircConn *IrcConnection) Init(oauth, nick string) {
	til := new(TwitchIRCListener)
	til.ArchiumDataIdentifier = ArchiumDataIdentifier
	til.ArchiumPrefix = ArchiumPrefix
	til.IrcConn = ircConn
	archiumCore.Register(til)
	ircConn.oauth = oauth
	ircConn.username = nick
	ircConn.Sendln("PASS " + oauth)
	ircConn.Sendln("NICK " + nick)
	ircConn.Sendln("CAP REQ :twitch.tv/tags")
	ircConn.Sendln("CAP REQ :twitch.tv/commands")
	go ircConn.start()
}

func (ircConn *IrcConnection) start() {
	for {
		line, err := ircConn.client.Recv()
		if err != nil {
			println("Error occured", err.Error())
			waitingChannel <- struct{}{}
			return
		}
		result := Parse(line)
		if result != nil {
			ev := archium.CreateEvent(1)
			ev.EventType = ArchiumPrefix + result.Channel + "." + strings.ToLower(result.Command)
			ev.EventSource = result.Channel
			ev.Data[ArchiumDataIdentifier] = result
			archiumCore.FireEvent(*ev)
		}
	}
	waitingChannel <- struct{}{}
}

func (ircConn *IrcConnection) Sendln(line string) {
	log.D("SENT", line)
	ircConn.client.Sendln(line)
}

func (ircConn *IrcConnection) Wait() {
	<-waitingChannel
}

func (ircConn *IrcConnection) Send(line, channel string) {
	log.D("SENT", line)
	ircConn.client.Sendln("PRIVMSG #" + channel + " :" + line)
}

func (ircConn *IrcConnection) Join(channel string) {
	ircConn.client.Sendln("JOIN #" + channel)
}

func (ircConn *IrcConnection) Leave(channel string) {
	ircConn.client.Sendln("PART #" + channel)
}

func (ircConn *IrcConnection) Quit() {
	ircConn.client.Sendln("QUIT")
	ircConn.client.Close()
}

func (ircConn *IrcConnection) Reconnect() {
	ircConn.client.Close()
	ircConn.Connect(ircConn.ip, ircConn.port)
	ircConn.Init(ircConn.oauth, ircConn.username)
	go ircConn.start()
}
