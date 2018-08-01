/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

import (
	"encoding/json"
	"time"

	"code.cerinuts.io/libs/goBlue/network"
)

//ChannelAuthenticated contains all available information about a channel, including private information
type ChannelAuthenticated struct {
	ID                           json.Number `json:"_id"`
	Mature                       bool        `json:"mature"`
	Partner                      bool        `json:"partner"`
	Status                       string      `json:"status"`
	BroadcasterLanguage          string      `json:"broadcaster_language"`
	DisplayName                  string      `json:"display_name"`
	Game                         string      `json:"game"`
	Language                     string      `json:"language"`
	Name                         string      `json:"name"`
	Logo                         string      `json:"logo"`
	VideoBanner                  string      `json:"video_banner"`
	ProfileBanner                string      `json:"profile_banner"`
	ProfileBannerBackgroundColor string      `json:"profile_banner_background_color"`
	URL                          string      `json:"url"`
	BroadcasterType              string      `json:"broadcaster_type"`
	Description                  string      `json:"description"`
	StreamKey                    string      `json:"stream_key"`
	Email                        string      `json:"email"`
	CreatedAt                    time.Time   `json:"created_at"`
	UpdatedAt                    time.Time   `json:"updated_at"`
	Views                        int         `json:"views"`
	Followers                    int         `json:"followers"`
}

//GetChannelAuthenticated returns authenticated channel information
func (tk *TwitchKraken) GetChannelAuthenticated(oauth string) (resp *ChannelAuthenticated, jsoerr *network.JSONError, err error) {
	resp = new(ChannelAuthenticated)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Request(BaseURL+"/channel", hMap, &resp)
	return resp, jsoerr, err
}

//Channel contains public information about a channel
type Channel struct {
	ID                           json.Number `json:"_id"`
	Mature                       bool        `json:"mature"`
	Partner                      bool        `json:"partner"`
	Status                       string      `json:"status"`
	BroadcasterLanguage          string      `json:"broadcaster_language"`
	DisplayName                  string      `json:"display_name"`
	Game                         string      `json:"game"`
	Language                     string      `json:"language"`
	Name                         string      `json:"name"`
	Logo                         string      `json:"logo"`
	VideoBanner                  string      `json:"video_banner"`
	ProfileBanner                string      `json:"profile_banner"`
	ProfileBannerBackgroundColor string      `json:"profile_banner_background_color"`
	URL                          string      `json:"url"`
	Description                  string      `json:"description"`
	BroadcasterType              string      `json:"broadcaster_type"`
	CreatedAt                    time.Time   `json:"created_at"`
	UpdatedAt                    time.Time   `json:"updated_at"`
	Views                        int         `json:"views"`
	Followers                    int         `json:"followers"`
}

//GetChannel returns public information about a channel
func (tk *TwitchKraken) GetChannel(channelID string) (resp *Channel, jsoerr *network.JSONError, err error) {
	resp = new(Channel)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID, hMap, &resp)
	return
}

//ChannelUpdate contains the Information you can set through the API
type ChannelUpdate struct {
	Data struct {
		Status             string `json:"status,omitempty"`
		Game               string `json:"game,omitempty"`
		Delay              string `json:"delay,omitempty"`
		ChannelFeedEnabled bool   `json:"channel_feed_enabled,omitempty"`
	} `json:"channel"`
}

//UpdateChannel posts the data object to the API to update the selected channel
func (tk *TwitchKraken) UpdateChannel(oauth, channelID string, data *ChannelUpdate) (resp *Channel, jsoerr *network.JSONError, err error) {
	resp = new(Channel)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	hMap["Content-Type"] = "application/json"
	jsoerr, err = jac.Put(BaseURL+"/channels/"+channelID, hMap, &data, &resp)
	return
}

//EditorList contains a list of users that have the editor permission
type EditorList struct {
	Editors []User `json:"users"`
}

//GetChannelEditors returns a list of all editors for a channel
func (tk *TwitchKraken) GetChannelEditors(oauth, channelID string) (resp *EditorList, jsoerr *network.JSONError, err error) {
	resp = new(EditorList)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/editors", hMap, &resp)
	return
}

//FollowsList contains information about whom a channel follows
type FollowsList struct {
	Cursor  string `json:"_cursor"`
	Total   int    `json:"_total"`
	Follows []struct {
		CreatedAt     string `json:"created_at"`
		Notifications bool   `json:"notifications"`
		User          User   `json:"user"`
	}
}

//GetChannelFollows fetches a list of all channels a channel follows
func (tk *TwitchKraken) GetChannelFollows(channelID, optLimit, optCursor, optOffset, optDirection string) (resp *FollowsList, jsoerr *network.JSONError, err error) {
	resp = new(FollowsList)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
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
		jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/follows?limit="+optLimit+"&offset="+optOffset+"&direction="+optDirection, hMap, &resp)
	} else {
		jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/follows?limit="+optLimit+"&offset="+optOffset+
			"&direction="+optDirection+"&cursor="+optCursor, hMap, &resp)
	}

	return
}

//GetChannelTeams returns a list of teams the channel is part of
func (tk *TwitchKraken) GetChannelTeams(channelID string) (resp *Teams, jsoerr *network.JSONError, err error) {
	resp = new(Teams)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/teams", hMap, &resp)
	return
}

//Subscriber contains information about a subscriber on a channel
type Subscriber struct {
	ID          json.Number `json:"_id"`
	CreatedAt   time.Time   `json:"created_at"`
	SubPlan     string      `json:"sub_plan"`
	SubPlanName string      `json:"sub_plan_name"`
	Subscriber  User        `json:"user"`
}

//Subscribers is a list of Subscribers
type Subscribers struct {
	Total       int          `json:"_total"`
	Subscribers []Subscriber `json:"subscriptions"`
}

//GetChannelSubscribers returns a list of all channel subscribers
func (tk *TwitchKraken) GetChannelSubscribers(oauth, channelID, optLimit, optOffset, optDirection string) (resp *Subscribers, jsoerr *network.JSONError, err error) {
	resp = new(Subscribers)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
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

	jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/subscriptions?limit="+optLimit+"&offset="+optOffset+"&direction="+optDirection, hMap, &resp)
	return
}

//IsSubscribedToChannel checks if a specific user is subscribed to a channel
func (tk *TwitchKraken) IsSubscribedToChannel(oauth, channelID, userID string) (resp *Subscriber, jsoerr *network.JSONError, err error) {
	resp = new(Subscriber)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/subscriptions/"+userID, hMap, &resp)
	return
}

//GetChannelVideos returns a list of all video on the channel
func (tk *TwitchKraken) GetChannelVideos(channelID, optLimit, optOffset, optBroadcastType, optLanguage, optSort string) (resp *Videos, jsoerr *network.JSONError, err error) {
	resp = new(Videos)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
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

	jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/videos?limit="+optLimit+"&offset="+optOffset+"&broadcast_type="+optBroadcastType+"&language="+optLanguage+"&sort="+optSort, hMap, &resp)
	return
}

//CommercialLength is..well, the length of a commercial
type CommercialLength struct {
	Length int `json:"Length"`
}

//CommercialResponse is a response for a requested commercial
type CommercialResponse struct {
	Length     int    `json:"Length"`
	Message    string `json:"Message"`
	RetryAfter int    `json:"RetryAfter"`
}

//RunCommercial tries to run a commercial
func (tk *TwitchKraken) RunCommercial(oauth, channelID string, length *CommercialLength) (resp *CommercialResponse, jsoerr *network.JSONError, err error) {
	resp = new(CommercialResponse)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Content-Type"] = "application/json"
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Post(BaseURL+"/channels/"+channelID+"/videos", hMap, length, &resp)
	return
}

//ResetStreamkey resets the streamkey
func (tk *TwitchKraken) ResetStreamkey(oauth, channelID string) (resp *ChannelAuthenticated, jsoerr *network.JSONError, err error) {
	resp = new(ChannelAuthenticated)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Delete(BaseURL+"/channels/"+channelID+"/stream_key", hMap, &resp)
	return
}

//GetChannelCommunities returns a list of all communities a channel is part of
func (tk *TwitchKraken) GetChannelCommunities(channelID string) (resp *Communities, jsoerr *network.JSONError, err error) {
	resp = new(Communities)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/channels/"+channelID+"/communities", hMap, &resp)
	return
}

//CommunityIDs a list of community ids
type CommunityIDs struct {
	CommunityIDs []string `json:"community_ids"`
}

//AddChannelToCommunities adds a channel to one or more communities
func (tk *TwitchKraken) AddChannelToCommunities(oauth, channelID string, communityIDs CommunityIDs) (jsoerr *network.JSONError, err error) {
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	hMap["Content-Type"] = "application/json"
	jsoerr, err = jac.Put(BaseURL+"/channels/"+channelID+"/communities", hMap, &communityIDs, nil)
	return
}

//RemoveChannelFromAllCommunites remove a channel from all! communities it is part of
func (tk *TwitchKraken) RemoveChannelFromAllCommunites(oauth, channelID string) (jsoerr *network.JSONError, err error) {
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Delete(BaseURL+"/channels/"+channelID+"/community", hMap, nil)
	return
}
