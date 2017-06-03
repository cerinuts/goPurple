package irc

import (
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	//connect
	in1 := ":tmi.twitch.tv 001 <user> :Welcome, GLHF!"
	in2 := ":tmi.twitch.tv 002 <user> :Your host is tmi.twitch.tv"
	in3 := ":tmi.twitch.tv 003 <user> :This server is rather new"
	in4 := ":tmi.twitch.tv 004 <user> :-"
	in5 := ":tmi.twitch.tv 375 <user> :-"
	in6 := ":tmi.twitch.tv 372 <user> :You are in a maze of twisty passages."
	in7 := ":tmi.twitch.tv 376 <user> :>"
	//unknown command
	in8 := ":tmi.twitch.tv 421 <user> WHO :Unknown command"
	//join
	in9 := ":<user>!<user>@<user>.tmi.twitch.tv JOIN #<channel>"
	in10 := ":<user>.tmi.twitch.tv 353 <user> = #<channel> :<user>"
	in11 := ":<user>.tmi.twitch.tv 366 <user> #<channel> :End of /NAMES list"
	//part
	in12 := ":<user>!<user>@<user>.tmi.twitch.tv PART #<channel>"
	//privmsg (notag)
	in13 := ":<user>!<user>@<user>.tmi.twitch.tv PRIVMSG #<channel> :This is a sample message"
	//cap ack
	in14 := ":tmi.twitch.tv CAP * ACK :twitch.tv/membership"
	in15 := ":tmi.twitch.tv CAP * ACK :twitch.tv/tags"
	in16 := ":tmi.twitch.tv CAP * ACK :twitch.tv/commands"
	//someone joins
	in18 := ":<user>!<user>@<user>.tmi.twitch.tv JOIN #<channel>"
	//mod/unmod
	in19 := ":jtv MODE #<channel> +o <user>"
	in20 := ":jtv MODE #<channel> -o <user>"
	//names
	in21 := ":<user>.tmi.twitch.tv 353 <user> = #<channel> :<user> <user2> <user3>"
	in22 := ":<user>.tmi.twitch.tv 353 <user> = #<channel> :<user4> <user5> ... <userN>"
	in23 := ":<user>.tmi.twitch.tv 366 <user> #<channel> :End of /NAMES list"
	//someone parts
	in24 := ":<user>!<user>@<user>.tmi.twitch.tv PART #<channel>"
	//clearchat
	in25 := "@ban-reason=Follow\\sthe\\srules :tmi.twitch.tv CLEARCHAT #dallas :ronni"
	//globaluserstate
	in26 := "@color=<color>;display-name=<display-name>;emote-sets=<emote-sets>;turbo=<turbo>;user-id=<user-id>;user-type=<user-type> :tmi.twitch.tv GLOBALUSERSTATE"
	//privmsg with tags
	in27 := "@badges=<badges>;bits=<bits>;color=<color>;display-name=<display-name>;emotes=<emotes>;id=<id>;mod=<mod>;room-id=<room-id>;subscriber=<subscriber>;turbo=<turbo>;user-id=<user-id>;user-type=<user-type> :<user>!<user>@<user>.tmi.twitch.tv PRIVMSG #<channel> :<message>"
	//bits
	in28 := "@badges=staff/1,bits/1000;bits=100;color=;display-name=dallas;emotes=;id=b34ccfc7-4977-403a-8a94-33c6bac34fb8;mod=0;room-id=1337;subscriber=0;turbo=1;user-id=1337;user-type=staff :ronni!ronni@ronni.tmi.twitch.tv PRIVMSG #dallas :cheer100"
	//roomstate
	in29 := "@broadcaster-lang=<broadcaster-lang>;r9k=<r9k>;slow=<slow>;subs-only=<subs-only> :tmi.twitch.tv ROOMSTATE #<channel>"
	in30 := "@slow=10 :tmi.twitch.tv ROOMSTATE #dallas"
	//usernotice(sub)
	in31 := "@badges=<badges>;color=<color>;display-name=<display-name>;emotes=<emotes>;mod=<mod>;msg-id=<msg-id>;msg-param-months=<msg-param-months>;msg-param-sub-plan=<msg-param-sub-plan>;msg-param-sub-plan-name=<msg-param-sub-plan-name>;room-id=<room-id>;subscriber=<subscriber>;system-msg=<system-msg>;login=<user>;turbo=<turbo>;user-id=<user-id>;user-type=<user-type> :tmi.twitch.tv USERNOTICE #<channel> :<message>"
	//userstate
	in32 := "@color=<color>;display-name=<display-name>;emote-sets=<emotes>;mod=<mod>;subscriber=<subscriber>;turbo=<turbo>;user-type=<user-type> :tmi.twitch.tv USERSTATE #<channel>"
	//clearchat no tag
	in33 := ":tmi.twitch.tv CLEARCHAT #<channel> :<user>"
	//host
	in34 := ":tmi.twitch.tv HOSTTARGET #hosting_channel <channel> [<number-of-viewers>]"
	//unhost
	in35 := ":tmi.twitch.tv HOSTTARGET #hosting_channel :- [<number-of-viewers>]"
	//notice
	in36 := "@msg-id=slow_off :tmi.twitch.tv NOTICE #dallas :This room is no longer in slow mode."
	//Reconnect
	in37 := ":tmi.twitch.tv RECONNECT"
	//Roomstate no tag
	in38 := ":tmi.twitch.tv ROOMSTATE #<channel>"
	//usernotice no tag
	in39 := ":tmi.twitch.tv USERNOTICE #<channel> :message"

	println("---1---")
	parse(in1)
	println("---2---")
	parse(in2)
	println("---3---")
	parse(in3)
	println("---4---")
	parse(in4)
	println("---5---")
	parse(in5)
	println("---6---")
	parse(in6)
	println("---7---")
	parse(in7)
	println("---8---")
	parse(in8)
	println("---9---")
	parse(in9)
	println("---10---")
	parse(in10)
	println("---11---")
	parse(in11)
	println("---12---")
	parse(in12)
	println("---13---")
	parse(in13)
	println("---14---")
	parse(in14)
	println("---15---")
	parse(in15)
	println("---16---")
	parse(in16)
	println("---17 SKIP---")
	//	parse(in17)
	println("---18---")
	parse(in18)
	println("---19---")
	parse(in19)
	println("---20---")
	parse(in20)
	println("---21---")
	parse(in21)
	println("---22---")
	parse(in22)
	println("---23---")
	parse(in23)
	println("---24---")
	parse(in24)
	println("---25---")
	parse(in25)
	println("---26---")
	parse(in26)
	println("---27---")
	parse(in27)
	println("---28---")
	parse(in28)
	println("---29---")
	parse(in29)
	println("---30---")
	parse(in30)
	println("---31---")
	parse(in31)
	println("---32---")
	parse(in32)
	println("---33---")
	parse(in33)
	println("---34---")
	parse(in34)
	println("---35---")
	parse(in35)
	println("---36---")
	parse(in36)
	println("---37---")
	parse(in37)
	println("---38---")
	parse(in38)
	println("---39---")
	parse(in39)
	println("---end---")
}

func parse(line string) {
	//orig := ""^(?:@([^ ]+) )?(?:[:](\\S+) )?(\\S+ )?#?(\\S+)?(?: [:](.+))?$""
	//  re := "^(?:@([^ ]+) )?(?:[:](\\S+) )?(\\S+)(?: (?!:)(.+?))?(?: [:](.+))?$"
	//	re := "^(?:@([^ ]+) )?(?:[:](\\S+) )?(\\S+)(?: ([.*[^:]])(.+?))?(?: [:](.+))?$"
//	re := "^(?:@([^ ]+) )?(?:[:](\\S+) )?(\\S+)(?: #?([\\S]+?))?(?: [:](.+))?$"
	re := "^(?:@([^ ]+) )?(?:[:](\\S+) )?(\\S+)(?: #?(?:)(.+?))?(?: [:](.+))?$"
	val := regexp.MustCompile(re).FindAllStringSubmatch(line, -1)
	if val != nil {
		//			result.Tags = make(map[string]string)

		for i, _ := range val {
			println(val[i][0])
			println("1:", val[i][1], "2:", val[i][2], "3:", val[i][3], "4:", val[i][4], "5:", val[i][5])
		}

	} else {
		println("err")
	}
}
