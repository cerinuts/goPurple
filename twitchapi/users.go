package twitchapi

import (
	"encoding/json"
	"gitlab.ceriath.net/libs/goBlue/network"
	"strings"
	"time"
)

type User struct {
	ID           json.Number `json:"_id"`
	Bio          string      `json:"bio"`
	Created_at   time.Time   `json:"created_at"`
	Display_name string      `json:"display_name"`
	Logo         string      `json:"logo"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Updated_at   time.Time   `json:"updated_at"`
}

type Users struct {
	Total int    `json:"_total"`
	Users []User `json:"users"`
}

func (tk *TwitchKraken) GetUsersByNames(names []string) (resp *Users, jsoerr *network.JsonError, err error) {
	resp = new(Users)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID

	logins := ""
	for _, val := range names {
		logins = logins + val + ","
	}
	logins = strings.TrimRight(logins, ",")
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/users?login="+logins, hMap, &resp)
	return
}

func (tk *TwitchKraken) GetUserByName(name string) (resp *Users, jsoerr *network.JsonError, err error) {
	resp = new(Users)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/users?login="+name, hMap, &resp)
	return
}

type Following struct {
	CreatedAt     time.Time `json:"created_at"`
	Channel       Channel   `json:"channel"`
	Notifications bool      `json:"notifications"`
}

func (tk *TwitchKraken) IsUserFollowingChannel(userId, channelId string) (resp *Following, jsoerr *network.JsonError, err error) {
	resp = new(Following)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/users/"+userId+"/follows/channels/"+channelId, hMap, &resp)
	return
}

type Badge struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type ChatterInfo struct {
	ID            string  `json:"_id"`
	Login         string  `json:"login"`
	DisplayName   string  `json:"display_name"`
	Color         string  `json:"color"`
	IsVerifiedBot bool    `json:"is_verified_bot"`
	IsKnownBot    bool    `json:"is_known_bot"`
	Badges        []Badge `json:"badges"`
}

func (htk *HiddenKraken) GetChatInformationForUser(userId string) (resp *ChatterInfo, jsoerr *network.JsonError, err error) {
	resp = new(ChatterInfo)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = htk.Tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/users/"+userId+"/chat", hMap, &resp)
	return
}

func (htk *HiddenKraken) GetChatInformationForUserByChannel(userId, channelId string) (resp *ChatterInfo, jsoerr *network.JsonError, err error) {
	resp = new(ChatterInfo)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = htk.Tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/users/"+userId+"/chat/channels/"+channelId, hMap, &resp)
	return
}
