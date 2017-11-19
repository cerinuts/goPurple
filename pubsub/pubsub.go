package pubsub

import (
	"encoding/json"
	"gitlab.ceriath.net/libs/goBlue/archium"
	"gitlab.ceriath.net/libs/goBlue/log"
	"gitlab.ceriath.net/libs/goBlue/util"
	"golang.org/x/net/websocket"
	"time"
)

var ArchiumPrefix = "twitch.pubsub."
var ArchiumDataIdentifier = "Message"

type Pubsub struct {
	Url string
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

type PsResult struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
	Error string `json:"error"`
	Data  struct {
		Topic   string                 `json:"topic"`
		RawMessage string `json:"message"`
		Message map[string]interface{} `json:"-"`
	} `json:"data"`
}

func (ps *Pubsub) Connect() {
	ws, err := websocket.Dial(ps.Url, "", "http://localhost/")
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

func (ps *Pubsub) Wait() {
	for {
		err, msg := ps.recv()
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
	data, _ := json.Marshal(ping)
	ps.ws.Write(data)
}

func (ps *Pubsub) recv() (error, []byte) {
	var msg = make([]byte, 2048)
	var n int
	n, err := ps.ws.Read(msg)
	if err != nil {
		log.E(err)
		return err, nil
	}
	return nil, msg[:n]
}

func (ps *Pubsub) parse(msg []byte) {
	psR := new(PsResult)
	json.Unmarshal(msg, psR)
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

func (ps *Pubsub) Close() {
	ps.ws.Close()
}

func (ps *Pubsub) reconnect() {
	ps.Close()
	ps.Connect()
}
