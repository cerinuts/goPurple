/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

import "code.cerinuts.io/libs/goBlue/network"

//EmoteList is a list of emotes
type EmoteList struct {
	Emoticons []Emoticon `json:"emoticons"`
}

//EmoteMap is a map of emotes, mapped by code
type EmoteMap struct {
	Emoticons map[string]Emoticon `json:"emoticons"`
}

//Emoticon contains information about an emote
type Emoticon struct {
	Code        string `json:"code"`
	EmoticonSet int    `json:"emoticon_set"`
	ID          int    `json:"id"`
}

//GetEmoticons fetches a list of all! emotes on twitch. This can take a long time and transfer a lot of data (~a few MB).
//You really should cache this list, it isn't cached here on purpose to provide freedom to the dev
func (tk *TwitchKraken) GetEmoticons(emotesets []string) (resp *EmoteList, jsoerr *network.JSONError, err error) {
	resp = new(EmoteList)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	if len(emotesets) > 0 {
		sets := ""
		for _, set := range emotesets {
			sets += set + ","
		}
		jsoerr, err = jac.Request(BaseURL+"/chat/emoticon_images?emotesets="+sets[:len(sets)-1], hMap, &resp)
	} else {
		jsoerr, err = jac.Request(BaseURL+"/chat/emoticon_images", hMap, &resp)
	}
	return
}

//GetEmoticonMap fetches a list of all! emotes on twitch as a map, mapped by code. This can take a long time and transfer a lot of data (~a few MB).
//You really should cache this list, it isn't cached here on purpose to provide freedom to the dev
func (tk *TwitchKraken) GetEmoticonMap(emotesets []string) (resp *EmoteMap, jsoerr *network.JSONError, err error) {
	resp = new(EmoteMap)
	resp.Emoticons = make(map[string]Emoticon)
	list, jsoerr, err := tk.GetEmoticons(emotesets)
	if jsoerr != nil || err != nil {
		return
	}

	for _, emote := range list.Emoticons {
		resp.Emoticons[emote.Code] = emote
	}

	return
}
