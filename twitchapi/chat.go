package twitchapi

import "gitlab.ceriath.net/libs/goBlue/network"

type EmoteList struct {
	Emoticons []Emoticon `json:"emoticons"`
}

type EmoteMap struct {
	Emoticons map[string]Emoticon `json:"emoticons"`
}

type Emoticon struct {
	Code        string `json:"code"`
	EmoticonSet int    `json:"emoticon_set"`
	ID          int    `json:"id"`
}

func (tk *TwitchKraken) GetEmoticons(emotesets []string) (resp *EmoteList, jsoerr *network.JsonError, err error) {
	resp = new(EmoteList)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	if len(emotesets) > 0 {
		sets := ""
		for _, set := range emotesets {
			sets += set + ","
		}
		jsoerr, err = jac.Request("https://api.twitch.tv/kraken/chat/emoticon_images?emotesets="+sets[:len(sets)-1], hMap, &resp)
	} else {
		jsoerr, err = jac.Request("https://api.twitch.tv/kraken/chat/emoticon_images", hMap, &resp)
	}
	return
}

func (tk *TwitchKraken) GetEmoticonMap(emotesets []string) (resp *EmoteMap, jsoerr *network.JsonError, err error) {
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
