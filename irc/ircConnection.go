package irc

import (
	"github.com/ceriath/goBlue/archium"
	"github.com/ceriath/goBlue/log"
	"github.com/ceriath/goBlue/network"
	"strings"
	"time"
)

var waitingChannel = make(chan struct{})
var archiumCore = archium.ArchiumCore

var ArchiumPrefix = "twitch.irc."
var ArchiumDataIdentifier = "Message"

type IrcConnection struct {
	client                      *network.Client
	oauth, Username, Host, Port string
	JoinedChannels              []string
	currentReconnectAttempts    int
}

func (ircConn *IrcConnection) Connect(ip, port string) error {
	cli := new(network.Client)
	ircConn.Host = ip
	ircConn.Port = port
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
	ircConn.Username = nick
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
			ircConn.Reconnect()
			//			waitingChannel <- struct{}{}
			return
		}
		result := Parse(line)
		if result != nil {
			ev := archium.CreateEvent()
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
	ircConn.JoinedChannels = append(ircConn.JoinedChannels, channel)
	ircConn.client.Sendln("JOIN #" + channel)
}

func (ircConn *IrcConnection) Leave(channel string) {
	for i, c := range ircConn.JoinedChannels {
		if c == channel {
			ircConn.JoinedChannels = append(ircConn.JoinedChannels[:i], ircConn.JoinedChannels[i+1:]...)
			break
		}
	}
	ircConn.client.Sendln("PART #" + channel)
}

func (ircConn *IrcConnection) Quit() {
	ircConn.client.Sendln("QUIT")
	ircConn.client.Close()
}

func (ircConn *IrcConnection) Reconnect() {
	//	ircConn.client.Close()
	if ircConn.currentReconnectAttempts >= 11 {
		log.F("11 attempts to reconnect failed, giving up.")
		return
	}
	ircConn.currentReconnectAttempts++
	log.I("Trying to recover, attempt:", ircConn.currentReconnectAttempts)
	time.Sleep(time.Duration(ircConn.currentReconnectAttempts*10) * time.Second)
	ircConn.Connect(ircConn.Host, ircConn.Port)
	ircConn.Init(ircConn.oauth, ircConn.Username)
}
