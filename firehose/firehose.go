/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package firehose is used for twitch's firehose
package firehose

import (
	"encoding/json"
	"strings"

	"code.cerinuts.io/libs/goBlue/log"
	"code.cerinuts.io/libs/goBlue/network"
)

//AppName is the name of the application
const AppName string = "goPurple/firehose"

//VersionMajor 0 means in development, >1 ensures compatibility with each minor version, but breakes with new major version
const VersionMajor string = "0"

//VersionMinor introduces changes that require a new version number. If the major version is 0, they are likely to break compatibility
const VersionMinor string = "1"

//VersionBuild is the type of this release. s(table), b(eta), d(evelopment), n(ightly)
const VersionBuild string = "s"

//FullVersion contains the full name and version of this package in a printable string
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

const firehoseURL = "https://tmi.twitch.tv/firehose?oauth_token="

//Firehose offers an interface to twitch tmi's firehose. Use with caution.
type Firehose struct {
	client   *network.EventsourceClient
	msgQueue chan Message
}

type internalMessage struct {
	Command string `json:"command"`
	Room    string `json:"room"`
	Nick    string `json:"nick"`
	Target  string `json:"target"`
	Body    string `json:"body"`
	Tags    string `json:"tags"`
}

//Message is the Firehose version of an IRC message
type Message struct {
	Command string            `json:"command"`
	Room    string            `json:"room"`
	Nick    string            `json:"nick"`
	Target  string            `json:"target"`
	Body    string            `json:"body"`
	Tags    map[string]string `json:"tags"`
	Raw     string
}

//Connect to Firehose with given OAuth token
func (f *Firehose) Connect(token string) (messageQueue chan Message, err error) {
	f.client = new(network.EventsourceClient)
	_, err = f.client.Subscribe(firehoseURL + token)
	if err != nil {
		return nil, err
	}
	f.msgQueue = make(chan Message)
	go f.handle()

	return f.msgQueue, nil
}

//Disconnect from firehose
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

func parse(ev *network.Event) (msg *Message, err error) {
	imsg := new(internalMessage)
	msg = new(Message)
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
