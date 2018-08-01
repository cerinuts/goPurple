/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

import (
	"encoding/json"
	"strings"
	"time"

	"code.cerinuts.io/libs/goBlue/network"
)

//User contains information about a user
type User struct {
	ID          json.Number `json:"_id"`
	Bio         string      `json:"bio"`
	CreatedAt   time.Time   `json:"created_at"`
	DisplayName string      `json:"display_name"`
	Logo        string      `json:"logo"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

//Users is a list of users
type Users struct {
	Total int    `json:"_total"`
	Users []User `json:"users"`
}

//GetUsersByNames returns a list of users searched by their names
func (tk *TwitchKraken) GetUsersByNames(names []string) (resp *Users, jsoerr *network.JSONError, err error) {
	resp = new(Users)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID

	logins := ""
	for _, val := range names {
		logins = logins + val + ","
	}
	logins = strings.TrimRight(logins, ",")
	jsoerr, err = jac.Request(BaseURL+"/users?login="+logins, hMap, &resp)
	return
}

//GetUserByName returns a single user searched by name
func (tk *TwitchKraken) GetUserByName(name string) (resp *Users, jsoerr *network.JSONError, err error) {
	resp = new(Users)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/users?login="+name, hMap, &resp)
	return
}

//Following contains information about a follow relation of a user
type Following struct {
	CreatedAt     time.Time `json:"created_at"`
	Channel       Channel   `json:"channel"`
	Notifications bool      `json:"notifications"`
}

//IsUserFollowingChannel returns a relation if a specific user is following a channel
func (tk *TwitchKraken) IsUserFollowingChannel(userID, channelID string) (resp *Following, jsoerr *network.JSONError, err error) {
	resp = new(Following)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/users/"+userID+"/follows/channels/"+channelID, hMap, &resp)
	return
}

//Badge contains a badge
type Badge struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

//ChatterInfo contains chat related information about a viewer in a specific chat
type ChatterInfo struct {
	ID            string  `json:"_id"`
	Login         string  `json:"login"`
	DisplayName   string  `json:"display_name"`
	Color         string  `json:"color"`
	IsVerifiedBot bool    `json:"is_verified_bot"`
	IsKnownBot    bool    `json:"is_known_bot"`
	Badges        []Badge `json:"badges"`
}

//GetChatInformationForUser returns information about a user related to their own chat
func (htk *HiddenKraken) GetChatInformationForUser(userID string) (resp *ChatterInfo, jsoerr *network.JSONError, err error) {
	resp = new(ChatterInfo)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = htk.Tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/users/"+userID+"/chat", hMap, &resp)
	return
}

//GetChatInformationForUserByChannel returns information about a user related to any chat
func (htk *HiddenKraken) GetChatInformationForUserByChannel(userID, channelID string) (resp *ChatterInfo, jsoerr *network.JSONError, err error) {
	resp = new(ChatterInfo)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = htk.Tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/users/"+userID+"/chat/channels/"+channelID, hMap, &resp)
	return
}
