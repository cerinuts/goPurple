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

//Stream contains information about a stream
type Stream struct {
	Stream struct {
		ID          json.Number `json:"_id"`
		Game        string      `json:"game"`
		Viewers     json.Number `json:"viewers"`
		VideoHeight json.Number `json:"video_height"`
		AverageFps  json.Number `json:"average_fps"`
		Delay       json.Number `json:"delay"`
		CreatedAt   time.Time   `json:"created_at"`
		IsPlaylist  bool        `json:"is_playlist"`
		StreamType  string      `json:"stream_type"`
		Preview     struct {
			Small    string `json:"small"`
			Medium   string `json:"medium"`
			Large    string `json:"large"`
			Template string `json:"template"`
		} `json:"preview"`
		Chan Channel `json:"channel"`
	} `json:"stream"`
}

//GetStream gets information about a stream
func (tk *TwitchKraken) GetStream(channelID, optType string) (resp *Stream, jsoerr *network.JSONError, err error) {
	resp = new(Stream)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID

	if optType == "" {
		optType = "live"
	}

	jsoerr, err = jac.Request(BaseURL+"/streams/"+channelID+"?stream_type="+optType, hMap, &resp)
	return
}
