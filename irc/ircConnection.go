/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package irc is used for twitch's irc
package irc

import (
	"strings"
	"time"

	"code.cerinuts.io/libs/goBlue/archium"
	"code.cerinuts.io/libs/goBlue/log"
	"code.cerinuts.io/libs/goBlue/network"
	"code.cerinuts.io/libs/goBlue/util"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goPurple/irc", "0", "4", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

var waitingChannel = make(chan int)
var archiumCore = archium.ArchiumCore
var runningKeepalive = false

//ArchiumPrefix is the prefix used for twitch irc messages on archium
var ArchiumPrefix = "twitch.irc."

//ArchiumDataIdentifier is the identifier to find the actual irc message in the archium message
var ArchiumDataIdentifier = "Message"

//Connection holds everything required for an Irc connection to twitch
type Connection struct {
	client                      *network.Client
	oauth, Username, Host, Port string
	JoinedChannels              []string
	currentReconnectAttempts    int
	openQueries                 int
	closed                      bool
	ModOnly                     bool
	runningReconnect            bool
	privmsgLimiter              *network.Tokenbucket
	joinLimiter                 *network.Tokenbucket
	lastActivity                time.Time
}

//Connect opens a connection to the irc server
func (ircConn *Connection) Connect(ip, port string) error {
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

//Init initalizes the connection and logs in
func (ircConn *Connection) Init(oauth, nick string) {
	til := new(TwitchIRCListener)
	til.ArchiumDataIdentifier = ArchiumDataIdentifier
	til.ArchiumPrefix = ArchiumPrefix
	til.IrcConn = ircConn
	archiumCore.Register(til)
	ircConn.JoinedChannels = make([]string, 1)
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

func (ircConn *Connection) start() {
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

func keepalive(ircConn *Connection) {
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
		time.Sleep(5 * time.Second)
	}
}

//Sendln sends a raw string to irc. Pls consider using Send instead expect you know what you are doing.
func (ircConn *Connection) Sendln(line string) {
	go ircConn.sendlnInternal(line)
}

//Wait waits until the queue is done and a message to waitingChannel is sent, e.g. the connection is terminated
func (ircConn *Connection) Wait() {
	for ircConn.openQueries > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	<-waitingChannel
}

//WaitForQueue waits until all queued up irc messages are sent within ratelimits
func (ircConn *Connection) WaitForQueue() {
	for ircConn.openQueries > 0 {
		time.Sleep(1 * time.Millisecond)
	}
}

//Send sends a message to a channel
func (ircConn *Connection) Send(line, channel string) {
	go ircConn.sendInternal("PRIVMSG #" + channel + " :" + line)
}

//BlockingSend sends a message to a channel and blocks until its actually sent within ratelimits
func (ircConn *Connection) BlockingSend(line, channel string) {
	ircConn.sendInternal("PRIVMSG #" + channel + " :" + line)
}

//Join joins a channel
func (ircConn *Connection) Join(channel string) {
	ircConn.JoinedChannels = append(ircConn.JoinedChannels, channel)
	go ircConn.joinInternal("JOIN #" + channel)
}

//Leave leaves a channel
func (ircConn *Connection) Leave(channel string) {
	util.RemoveFromStringSlice(ircConn.JoinedChannels, channel)
	go ircConn.sendlnInternal("PART #" + channel)
}

//Quit closes the irc connection
func (ircConn *Connection) Quit() {
	ircConn.closed = true
	ircConn.sendlnInternal("QUIT")
	ircConn.client.Close()
}

func (ircConn *Connection) sendlnInternal(line string) {
	ircConn.client.Sendln(line)
	log.D("SENT", line)
}

func (ircConn *Connection) sendInternal(line string) {
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

func (ircConn *Connection) joinInternal(line string) {
	err := ircConn.joinLimiter.Wait()
	if err != nil {
		ircConn.joinInternal(line)
		return
	}
	ircConn.client.Sendln(line)
	log.D("SENT", line)
}

//Reconnect reconnects to the irc server.
func (ircConn *Connection) Reconnect() {
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
	for _, k := range tmpJoinedChannels {
		ircConn.Join(k)
	}
}
