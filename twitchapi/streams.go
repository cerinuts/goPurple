package twitchapi

import (
	"encoding/json"
	"github.com/ceriath/goBlue/network"
	"time"
)

type Stream struct {
	Stream struct {
		ID          json.Number `json:"_id"`
		Game        string      `json:"game"`
		Viewers     json.Number         `json:"viewers"`
		VideoHeight json.Number         `json:"video_height"`
		AverageFps  json.Number     `json:"average_fps"`
		Delay       json.Number         `json:"delay"`
		CreatedAt   time.Time   `json:"created_at"`
		IsPlaylist  bool        `json:"is_playlist"`
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

	if optType == "" {
		optType = "live"
	}

	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/streams/"+channelId+"?stream_type="+optType, hMap, &resp)
	return
}
