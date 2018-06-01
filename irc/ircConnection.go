package irc

import (
	"strings"
	"time"

	"gitlab.ceriath.net/libs/goBlue/archium"
	"gitlab.ceriath.net/libs/goBlue/log"
	"gitlab.ceriath.net/libs/goBlue/network"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goPurple/irc", "0", "2", "b"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

var waitingChannel = make(chan int)
var archiumCore = archium.ArchiumCore
var runningKeepalive = false

var ArchiumPrefix = "twitch.irc."
var ArchiumDataIdentifier = "Message"

type IrcConnection struct {
	client                      *network.Client
	oauth, Username, Host, Port string
	JoinedChannels              map[string]struct{}
	currentReconnectAttempts    int
	closed                      bool
	ModOnly                     bool
	privmsgLimiter              *network.Tokenbucket
	joinLimiter                 *network.Tokenbucket
	runningReconnect            bool
	openQueries                 int
	lastActivity                time.Time
}

func (ircConn *IrcConnection) Connect(ip, port string) error {
	cli := new(network.Client)
	ircConn.Host = ip
	ircConn.Port = port
	ircConn.client = cli
	ircConn.openQueries = 0
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
	if ircConn.ModOnly {
		ircConn.privmsgLimiter = network.NewTokenbucket(30*time.Second, 100)
	} else {
		ircConn.privmsgLimiter = network.NewTokenbucket(30*time.Second, 20)
	}
	ircConn.joinLimiter = network.NewTokenbucket(15*time.Second, 50)
	ircConn.lastActivity = time.Now()
	if !strings.HasPrefix(nick, "justinfan") {
		ircConn.Sendln("PASS " + oauth)
	}
	ircConn.Sendln("NICK " + nick)
	ircConn.Sendln("CAP REQ :twitch.tv/tags twitch.tv/commands")

	time.Sleep(3 * time.Second)

	(*(til.IrcConn)).runningReconnect = false
	go ircConn.start()
	go keepalive(ircConn)
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

func keepalive(ircConn *IrcConnection) {
	if runningKeepalive {
		return
	}
	runningKeepalive = true
	for {
		if time.Since(ircConn.lastActivity) >= 10*time.Second {
			ircConn.Sendln("PING")
			time.Sleep(5 * time.Second)
			if time.Since(ircConn.lastActivity) >= 15*time.Second {
				ircConn.Reconnect()
				return
			}
		}
	}
}

func (ircConn *IrcConnection) Sendln(line string) {
	go ircConn.sendlnInternal(line)
}

func (ircConn *IrcConnection) Wait() {
	for ircConn.openQueries > 0 {
	}
	<-waitingChannel
}

func (ircConn *IrcConnection) WaitForQueue() {
	for ircConn.openQueries > 0 {
	}
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
	ircConn.openQueries++
	err := ircConn.privmsgLimiter.Wait()
	if err != nil {
		ircConn.sendInternal(line)
		return
	}
	ircConn.openQueries--

	ircConn.client.Sendln(line)
	log.D("SENT", line)
}

func (ircConn *IrcConnection) joinInternal(line string) {
	err := ircConn.joinLimiter.Wait()
	if err != nil {
		ircConn.joinInternal(line)
		return
	}
	ircConn.client.Sendln(line)
	log.D("SENT", line)
}

func (ircConn *IrcConnection) Reconnect() {
	if ircConn.runningReconnect {
		return
	}
	ircConn.runningReconnect = true
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
	for k, _ := range tmpJoinedChannels {
		ircConn.Join(k)
	}
}
