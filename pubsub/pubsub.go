/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package pubsub is used for twitch's pubsub. This implementation is not tested at all.
package pubsub

import (
	"encoding/json"
	"time"

	"code.cerinuts.io/libs/goBlue/archium"
	"code.cerinuts.io/libs/goBlue/log"
	"code.cerinuts.io/libs/goBlue/util"
	"golang.org/x/net/websocket"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goPurple/pubsub", "0", "1", "d"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//ArchiumPrefix for twitch pubsub messages
var ArchiumPrefix = "twitch.pubsub."

//ArchiumDataIdentifier which contains the actual payload
var ArchiumDataIdentifier = "Message"

//Pubsub is a simple container for the pubsub
type Pubsub struct {
	URL string
	ws  *websocket.Conn
}

type psListen struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Data  struct {
		Topics    []string `json:"topics"`
		AuthToken string   `json:"auth_token"`
	} `json:"data"`
}

//PsResult is the message sent by pubsub
type PsResult struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Error string `json:"error"`
	Data  struct {
		Topic      string                 `json:"topic"`
		RawMessage string                 `json:"message"`
		Message    map[string]interface{} `json:"-"`
	} `json:"data"`
}

//Connect connects to a pubsub
func (ps *Pubsub) Connect() {
	ws, err := websocket.Dial(ps.URL, "", "http://localhost/")
	if err != nil {
		log.F(err)
	}
	ps.ws = ws
	go func() {
		for {
			ps.ping()
			time.Sleep(3 * time.Minute)
		}
	}()
}

//Recv waits for incoming messages
func (ps *Pubsub) Recv() {
	for {
		msg, err := ps.recv()
		if err != nil {
			log.E(err)
			if err.Error() == "EOF" {
				ps.Connect()
			}
			continue
		}
		go ps.parse(msg)
	}
}

//Listen subscribes to a list of topics
func (ps *Pubsub) Listen(topics []string, token string) {
	ls := new(psListen)
	ls.Type = "LISTEN"
	ls.Nonce = util.GetRandomAlphanumericString(10)
	ls.Data.Topics = topics
	ls.Data.AuthToken = token

	msg, err := json.Marshal(ls)
	if err != nil {
		log.F(err)
	}
	log.D("SENT", msg)
	ps.ws.Write(msg)

}

func (ps *Pubsub) ping() {
	ping := make(map[string]string)
	ping["type"] = "PING"
	data, err := json.Marshal(ping)
	if err != nil {
		return
	}
	ps.ws.Write(data)
}

func (ps *Pubsub) recv() ([]byte, error) {
	var msg = make([]byte, 2048)
	var n int
	n, err := ps.ws.Read(msg)
	if err != nil {
		log.E(err)
		return nil, err
	}
	return msg[:n], nil
}

func (ps *Pubsub) parse(msg []byte) {
	psR := new(PsResult)
	err := json.Unmarshal(msg, psR)
	if err != nil {
		log.E(err)
		return
	}
	switch psR.Type {
	case "PONG":
		return
	case "RECONNECT":
		ps.reconnect()
	case "RESPONSE":
		log.D(string(msg))
	default:
		if err := json.Unmarshal(msg, &psR); err != nil {
			log.E(err)
		}
		if err := json.Unmarshal([]byte(psR.Data.RawMessage), &psR.Data.Message); err != nil {
			log.E(err)
		}
		log.D(psR)
		ev := archium.CreateEvent()
		ev.EventType = ArchiumPrefix + psR.Data.Topic
		ev.EventSource = "Twitch PubSub"
		ev.Data[ArchiumDataIdentifier] = psR
		archium.ArchiumCore.FireEvent(*ev)

	}

}

//Close closes the connection
func (ps *Pubsub) Close() {
	ps.ws.Close()
}

func (ps *Pubsub) reconnect() {
	ps.Close()
	ps.Connect()
}
