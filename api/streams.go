package api

import (
"github.com/ceriath/goBlue/network"
"time"
)

type Stream struct {
	Stream struct {
		ID          int64     `json:"_id"`
		Game        string    `json:"game"`
		Viewers     int       `json:"viewers"`
		VideoHeight int       `json:"video_height"`
		AverageFps  int       `json:"average_fps"`
		Delay       int       `json:"delay"`
		CreatedAt   time.Time `json:"created_at"`
		IsPlaylist  bool      `json:"is_playlist"`
		Preview     struct {
			Small    string `json:"small"`
			Medium   string `json:"medium"`
			Large    string `json:"large"`
			Template string `json:"template"`
		} `json:"preview"`
		Chan Channel `json:"channel"`
	} `json:"stream"`
}



func (tk *TwitchKraken) GetStream(channelId, optType string) (resp *Stream, jsoerr *network.JsonError, err error) {
	resp = new(Stream)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	
	if optType == ""{
		optType = "live"
	}
	
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/streams/"+channelId+"?stream_type="+optType, hMap, &resp)
	return
}