package irc

import (
	"github.com/ceriath/goBlue/archium"
	"github.com/ceriath/goBlue/log"
	"github.com/ceriath/goBlue/network"
	"strings"
	"time"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goPurple/irc", "0", "1", "b"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

var waitingChannel = make(chan int)
var archiumCore = archium.ArchiumCore

var ArchiumPrefix = "twitch.irc."
var ArchiumDataIdentifier = "Message"

type IrcConnection struct {
	client                      *network.Client
	oauth, Username, Host, Port string
	JoinedChannels              map[string]struct{}
	currentReconnectAttempts    int
	closed                      bool
	ModOnly                     bool
	privmsgLimiter              *network.Ratelimiter
	joinLimiter                 *network.Ratelimiter
}

func (ircConn *IrcConnection) Connect(ip, port string) error {
	cli := new(network.Client)
	ircConn.Host = ip
	ircConn.Port = port
	ircConn.client = cli
	err := ircConn.client.Connect(ip, port)
	if err == nil {
		ircConn.closed = false
	}
	return err
}

func (ircConn *IrcConnection) Init(oauth, nick string) {
	til := new(TwitchIRCListener)
	til.ArchiumDataIdentifier = ArchiumDataIdentifier
	til.ArchiumPrefix = ArchiumPrefix
	til.IrcConn = ircConn
	archiumCore.Register(til)
	ircConn.JoinedChannels = make(map[string]struct{})
	ircConn.oauth = oauth
	ircConn.Username = nick
	ircConn.privmsgLimiter = new(network.Ratelimiter)
	if ircConn.ModOnly {
		ircConn.privmsgLimiter.Init("twitchirc", 100, 30*time.Second)
	} else {
		ircConn.privmsgLimiter.Init("twitchirc", 1, 1500*time.Millisecond)
	}
	ircConn.joinLimiter = new(network.Ratelimiter)
	ircConn.joinLimiter.Init("twitchirc-join", 50, 15*time.Second)
	ircConn.Sendln("PASS " + oauth)
	ircConn.Sendln("NICK " + nick)
	ircConn.Sendln("CAP REQ :twitch.tv/tags twitch.tv/commands")


	ircConn.currentReconnectAttempts = 0
	time.Sleep(3 * time.Second)
	go ircConn.start()
}

func (ircConn *IrcConnection) start() {
	for {
		line, err := ircConn.client.Recv()
		if err != nil {
			if ircConn.closed {
				break
			}
			log.E("Error occured", err.Error())
			ircConn.Reconnect()
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
	waitingChannel <- 1
}

func (ircConn *IrcConnection) Sendln(line string) {
	go ircConn.sendlnInternal(line)
}

func (ircConn *IrcConnection) Wait() {
	ircConn.privmsgLimiter.WaitForQueue()
	<-waitingChannel
}

func (ircConn *IrcConnection) WaitForQueue() {
	ircConn.privmsgLimiter.WaitForQueue()
}

func (ircConn *IrcConnection) Send(line, channel string) {
	go ircConn.sendInternal("PRIVMSG #" + channel + " :" + line)
}

func (ircConn *IrcConnection) BlockingSend(line, channel string) {
	ircConn.sendInternal("PRIVMSG #" + channel + " :" + line)
}

func (ircConn *IrcConnection) Join(channel string) {
	ircConn.JoinedChannels[channel] = struct{}{}
	go ircConn.joinInternal("JOIN #" + channel)
}

func (ircConn *IrcConnection) Leave(channel string) {
	delete(ircConn.JoinedChannels, channel)
	go ircConn.sendlnInternal("PART #" + channel)
}

func (ircConn *IrcConnection) Quit() {
	ircConn.closed = true
	ircConn.sendlnInternal("QUIT")
	ircConn.client.Close()
}

func (ircConn *IrcConnection) sendlnInternal(line string) {
	ircConn.client.Sendln(line)
	log.D("SENT", line)
}

func (ircConn *IrcConnection) sendInternal(line string) {
	<-ircConn.privmsgLimiter.Request(false)
	ircConn.client.Sendln(line)
	log.D("SENT", line)
}

func (ircConn *IrcConnection) joinInternal(line string) {
	<-ircConn.joinLimiter.Request(false)
	ircConn.client.Sendln(line)
	log.D("SENT", line)
}

func (ircConn *IrcConnection) Reconnect() {
	//	ircConn.client.Close()
	if ircConn.currentReconnectAttempts >= 11 {
		log.F("11 attempts to reconnect failed, giving up.")
		return
	}
	log.I("Trying to recover, attempt:", ircConn.currentReconnectAttempts)
	time.Sleep(time.Duration(ircConn.currentReconnectAttempts*10) * time.Second)
	tmpJoinedChannels := ircConn.JoinedChannels
	ircConn.currentReconnectAttempts++
	ircConn.Connect(ircConn.Host, ircConn.Port)
	ircConn.Init(ircConn.oauth, ircConn.Username)
	for k, _ := range  tmpJoinedChannels{
		ircConn.Join(k)
	}
}
