package irc

import (
	"regexp"
)

type IrcMessage struct {
	Tags                       map[string]string
	Raw, Channel, Msg, Command string
}

type IrcParser struct {
}

var regexMsg = "^(?:@([^ ]+) )?(?:[:](\\S+) )?(\\S+)(?: #?(?:)(.+?))?(?: [:](.+))?$"
var regexTags = "([^=;]+)=([^;]*)"

var RAW = 0
var TAGS = 1
var SERVER = 2
var COMMAND = 3
var CHANNEL = 4
var MESSAGE = 5

func (ircParser *IrcParser) Parse(line string) *IrcMessage {
	println(line)
	parsed := regexp.MustCompile(regexMsg).FindAllStringSubmatch(line, -1)
	if parsed != nil {
		result := new(IrcMessage)
		pLine := parsed[0]
		result.Raw = pLine[RAW]
		result.Channel = pLine[CHANNEL]
		result.Msg = pLine[MESSAGE]
		result.Command = pLine[COMMAND]
		val := regexp.MustCompile(regexTags).FindAllStringSubmatch(pLine[TAGS], -1)
		if val != nil {
			result.Tags = make(map[string]string)
			for i, _ := range val {
				for range val[i] {
					result.Tags[val[i][1]] = val[i][2]
				}
			}

		}
		//		for k, v := range result.Tags {
		//			println(k, v)
		//		}

		return result

	}
	println("Couldnt be parsed")
	return nil
}


