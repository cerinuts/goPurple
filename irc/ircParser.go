/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package irc is used for twitch's irc
package irc

import (
	"strings"
)

//Message contains an irc message
type Message struct {
	Tags                               map[string]string
	Raw, Channel, Msg, Command, Source string
}

const globalChannel = "global"

//Parse parses an irc message
func Parse(line string) *Message {
	msg := new(Message)
	msg.Raw = line
	if line[0] == '@' {
		line = line[1:]
		msgSplit := strings.SplitN(line, " ", 2)
		parse(msgSplit[1], msg)
		tagSplit := strings.Split(msgSplit[0], ";")
		msg.Tags = make(map[string]string)
		for v := range tagSplit {
			tag := strings.Split(tagSplit[v], "=")
			msg.Tags[tag[0]] = tag[1]
		}
	} else if line[0] == ':' {
		parse(line, msg)
	} else if strings.HasPrefix(line, "PING") {
		msg.Command = "PING"
		msg.Channel = globalChannel
		msg.Source = strings.SplitN(line, ":", 2)[1]
	} else if strings.HasPrefix(line, "PONG") {
		msg.Command = "PONG"
		msg.Channel = globalChannel
		msg.Source = strings.SplitN(line, ":", 2)[1]
	}
	return msg
}

func parse(line string, msg *Message) *Message {
	split := strings.SplitN(line[1:], " ", 3)
	msg.Source = split[0]
	msg.Command = split[1]
	switch msg.Command {
	case "001":
		fallthrough
	case "002":
		fallthrough
	case "003":
		fallthrough
	case "004":
		fallthrough
	case "372":
		fallthrough
	case "375":
		fallthrough
	case "376":
		fallthrough
	case "421":
		msg.Channel = globalChannel
		msg.Msg = strings.SplitN(split[2], " :", 2)[1]
	case "353":
		s := strings.SplitN(strings.SplitN(split[2], "= #", 2)[1], " :", 2)
		msg.Channel = s[0]
		msg.Msg = s[1]
	case "366":
		s := strings.SplitN(split[2], " ", 3)
		msg.Channel = s[1][1:]
		msg.Msg = s[2][1:]
	case "ROOMSTATE":
		fallthrough
	case "JOIN":
		fallthrough
	case "PART":
		msg.Channel = split[2][1:]
	case "CLEARCHAT":
		fallthrough
	case "NOTICE":
		fallthrough
	case "USERNOTICE":
		fallthrough
	case "USERSTATE":
		fallthrough
	case "PRIVMSG":
		s := strings.SplitN(split[2], " ", 2)
		msg.Channel = s[0][1:]
		if len(s) > 1 {
			msg.Msg = s[1][1:]
		}

	case "CAP":
		if strings.HasPrefix(split[2], "* ACK") {
			msg.Channel = globalChannel
			msg.Command = "CAPACK"
			msg.Msg = strings.SplitN(split[2], ":", 2)[1]
		}
	case "HOSTTARGET":
		fallthrough
	case "MODE":
		s := strings.SplitN(split[2], " ", 2)
		msg.Channel = s[0][1:]
		msg.Msg = s[1]
	}

	return msg
}
