package firehose

import (
	"encoding/json"
	"strings"

	"gitlab.ceriath.net/libs/goBlue/log"
	"gitlab.ceriath.net/libs/goBlue/network"
)

type Firehose struct {
	client   *network.EventsourceClient
	msgQueue chan FirehoseMessage
}

type internalFirehoseMessage struct {
	Command string `json:"command"`
	Room    string `json:"room"`
	Nick    string `json:"nick"`
	Target  string `json:"target"`
	Body    string `json:"body"`
	Tags    string `json:"tags"`
}

type FirehoseMessage struct {
	Command string            `json:"command"`
	Room    string            `json:"room"`
	Nick    string            `json:"nick"`
	Target  string            `json:"target"`
	Body    string            `json:"body"`
	Tags    map[string]string `json:"tags"`
	Raw     string
}

func (f *Firehose) Connect(token string) (messageQueue chan FirehoseMessage, err error) {
	f.client = new(network.EventsourceClient)
	_, err = f.client.Subscribe("https://tmi.twitch.tv/firehose?oauth_token=" + token)
	if err != nil {
		return nil, err
	}
	f.msgQueue = make(chan FirehoseMessage)
	go f.handle()

	return f.msgQueue, nil
}

func (f *Firehose) Disconnect() {
	f.client.Close()
}

func (f *Firehose) handle() {
	for {
		ev := <-f.client.Stream.EventQueue
		go func() {
			msg, err := parse(&ev)
			if err != nil {
				log.E(err)
				return
			}
			f.msgQueue <- *msg
		}()
	}
}

func parse(ev *network.Event) (msg *FirehoseMessage, err error) {
	imsg := new(internalFirehoseMessage)
	msg = new(FirehoseMessage)
	err = json.Unmarshal([]byte(ev.Payload), &imsg)
	if err != nil {
		return nil, err
	}
	msg.Body = imsg.Body
	msg.Command = imsg.Command
	msg.Nick = imsg.Nick
	msg.Room = imsg.Room
	msg.Target = imsg.Target
	msg.Raw = ev.Payload
	msg.Tags = make(map[string]string)
	for _, tag := range strings.Split(imsg.Tags, ";") {
		if len(tag) < 1 {
			return
		}
		tagSplit := strings.Split(tag, "=")
		msg.Tags[tagSplit[0]] = tagSplit[1]
	}

	return msg, err
}
