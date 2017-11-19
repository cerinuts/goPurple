package twitchapi

import (
	"encoding/json"
	"gitlab.ceriath.net/libs/goBlue/network"
	"time"
)

type ChannelAuthenticated struct {
	Mature                          bool        `json:"mature"`
	Status                          string      `json:"status"`
	Broadcaster_language            string      `json:"broadcaster_language"`
	Display_name                    string      `json:"display_name"`
	Game                            string      `json:"game"`
	Language                        string      `json:"language"`
	ID                              json.Number `json:"_id"`
	Name                            string      `json:"name"`
	Created_at                      time.Time   `json:"created_at"`
	Updated_at                      time.Time   `json:"updated_at"`
	Partner                         bool        `json:"partner"`
	Logo                            string      `json:"logo"`
	Video_banner                    string      `json:"video_banner"`
	Profile_banner                  string      `json:"profile_banner"`
	Profile_banner_background_color string      `json:"profile_banner_background_color"`
	Url                             string      `json:"url"`
	Views                           int         `json:"views"`
	Followers                       int         `json:"followers"`
	Broadcaster_type                string      `json:"broadcaster_type"`
	Description                     string      `json:"description"`
	Stream_key                      string      `json:"stream_key"`
	Email                           string      `json:"email"`
}

func (tk *TwitchKraken) GetChannelAuthenticated(oauth string) (resp *ChannelAuthenticated, jsoerr *network.JsonError, err error) {
	resp = new(ChannelAuthenticated)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channel", hMap, &resp)
	return
}

type Channel struct {
	Mature                          bool        `json:"mature"`
	Status                          string      `json:"status"`
	Broadcaster_language            string      `json:"broadcaster_language"`
	Display_name                    string      `json:"display_name"`
	Game                            string      `json:"game"`
	Language                        string      `json:"language"`
	ID                              json.Number `json:"_id"`
	Name                            string      `json:"name"`
	Created_at                      time.Time   `json:"created_at"`
	Updated_at                      time.Time   `json:"updated_at"`
	Partner                         bool        `json:"partner"`
	Logo                            string      `json:"logo"`
	Video_banner                    string      `json:"video_banner"`
	Profile_banner                  string      `json:"profile_banner"`
	Profile_banner_background_color string      `json:"profile_banner_background_color"`
	Url                             string      `json:"url"`
	Views                           int         `json:"views"`
	Followers                       int         `json:"followers"`
	Description                     string      `json:"description"`
	Broadcaster_type                string      `json:"broadcaster_type"`
}

func (tk *TwitchKraken) GetChannel(channelId string) (resp *Channel, jsoerr *network.JsonError, err error) {
	resp = new(Channel)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId, hMap, &resp)
	return
}

type ChannelUpdate struct {
	Data struct {
		Status               string `json:"status,omitempty"`
		Game                 string `json:"game,omitempty"`
		Delay                string `json:"delay,omitempty"`
		Channel_feed_enabled bool   `json:"channel_feed_enabled,omitempty"`
	} `json:"channel"`
}

func (tk *TwitchKraken) UpdateChannel(oauth, channelId string, data *ChannelUpdate) (resp *Channel, jsoerr *network.JsonError, err error) {
	resp = new(Channel)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	hMap["Content-Type"] = "application/json"
	jsoerr, err = jac.Put("https://api.twitch.tv/kraken/channels/"+channelId, hMap, &data, &resp)
	return
}

type EditorList struct {
	Editors []User `json:"users"`
}

func (tk *TwitchKraken) GetChannelEditors(oauth, channelId string) (resp *EditorList, jsoerr *network.JsonError, err error) {
	resp = new(EditorList)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/editors", hMap, &resp)
	return
}

type FollowsList struct {
	Cursor  string `json:"_cursor"`
	Total   int    `json:"_total"`
	Follows []struct {
		Created_at    string `json:"created_at"`
		Notifications bool   `json:"notifications"`
		User          User   `json:"user"`
	}
}

func (tk *TwitchKraken) GetChannelFollows(channelId, optLimit, optCursor, optOffset, optDirection string) (resp *FollowsList, jsoerr *network.JsonError, err error) {
	resp = new(FollowsList)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	if optLimit == "" {
		optLimit = "25"
	}

	if optOffset == "" {
		optOffset = "0"
	}

	if optDirection == "" {
		optDirection = "desc"
	}

	if optCursor == "" {
		jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/follows?limit="+optLimit+"&offset="+optOffset+"&direction="+optDirection, hMap, &resp)
	} else {
		jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/follows?limit="+optLimit+"&offset="+optOffset+
			"&direction="+optDirection+"&cursor="+optCursor, hMap, &resp)
	}

	return
}

func (tk *TwitchKraken) GetChannelTeams(channelId string) (resp *Teams, jsoerr *network.JsonError, err error) {
	resp = new(Teams)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/teams", hMap, &resp)
	return
}

type Subscriber struct {
	ID            json.Number `json:"_id"`
	Created_at    time.Time      `json:"created_at"`
	Sub_plan      string      `json:"sub_plan"`
	Sub_plan_name string      `json:"sub_plan_name"`
	Subscriber    User        `json:"user"`
}

type Subscribers struct {
	Total       int          `json:"_total"`
	Subscribers []Subscriber `json:"subscriptions"`
}

func (tk *TwitchKraken) GetChannelSubscribers(oauth, channelId, optLimit, optOffset, optDirection string) (resp *Subscribers, jsoerr *network.JsonError, err error) {
	resp = new(Subscribers)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth

	if optLimit == "" {
		optLimit = "25"
	}

	if optOffset == "" {
		optOffset = "0"
	}

	if optDirection == "" {
		optDirection = "asc"
	}

	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/subscriptions?limit="+optLimit+"&offset="+optOffset+"&direction="+optDirection, hMap, &resp)
	return
}

func (tk *TwitchKraken) IsSubscribedToChannel(oauth, channelId, userId string) (resp *Subscriber, jsoerr *network.JsonError, err error) {
	resp = new(Subscriber)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/subscriptions/"+userId, hMap, &resp)
	return
}

func (tk *TwitchKraken) GetChannelVideos(channelId, optLimit, optOffset, optBroadcastType, optLanguage, optSort string) (resp *Videos, jsoerr *network.JsonError, err error) {
	resp = new(Videos)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID

	if optLimit == "" {
		optLimit = "10"
	}

	if optOffset == "" {
		optOffset = "0"
	}

	if optSort == "" {
		optSort = "time"
	}

	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/videos?limit="+optLimit+"&offset="+optOffset+"&broadcast_type="+optBroadcastType+"&language="+optLanguage+"&sort="+optSort, hMap, &resp)
	return
}

type CommercialLength struct {
	Length int `json:"Length"`
}

type CommercialResponse struct {
	Length     int    `json:"Length"`
	Message    string `json:"Message"`
	RetryAfter int    `json:"RetryAfter"`
}

func (tk *TwitchKraken) RunCommercial(oauth, channelId string, length *CommercialLength) (resp *CommercialResponse, jsoerr *network.JsonError, err error) {
	resp = new(CommercialResponse)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Content-Type"] = "application/json"
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Post("https://api.twitch.tv/kraken/channels/"+channelId+"/videos", hMap, length, &resp)
	return
}

func (tk *TwitchKraken) ResetStreamkey(oauth, channelId string) (resp *ChannelAuthenticated, jsoerr *network.JsonError, err error) {
	resp = new(ChannelAuthenticated)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Delete("https://api.twitch.tv/kraken/channels/"+channelId+"/stream_key", hMap, &resp)
	return
}

func (tk *TwitchKraken) GetChannelCommunities(channelId string) (resp *Communities, jsoerr *network.JsonError, err error) {
	resp = new(Communities)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/channels/"+channelId+"/communities", hMap, &resp)
	return
}

type CommunityIds struct {
	Community_ids []string `json:"community_ids"`
}

func (tk *TwitchKraken) AddChannelToCommunities(oauth, channelId string, communityIds CommunityIds) (jsoerr *network.JsonError, err error) {
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	hMap["Content-Type"] = "application/json"
	jsoerr, err = jac.Put("https://api.twitch.tv/kraken/channels/"+channelId+"/communities", hMap, &communityIds, nil)
	return
}

func (tk *TwitchKraken) RemoveChannelFromAllCommunites(oauth, channelId string) (jsoerr *network.JsonError, err error) {
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Delete("https://api.twitch.tv/kraken/channels/"+channelId+"/community", hMap, nil)
	return
}
