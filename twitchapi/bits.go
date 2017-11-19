package twitchapi

import (
	"gitlab.ceriath.net/libs/goBlue/network"
	"time"
)

type Cheeremotes struct {
	Actions []struct {
		Prefix string   `json:"prefix"`
		Scales []string `json:"scales"`
		Tiers  []struct {
			MinBits int    `json:"min_bits"`
			ID      string `json:"id"`
			Color   string `json:"color"`
			Images  struct {
				Dark struct {
					Animated struct {
						Num1 string `json:"1"`
						Num2 string `json:"2"`
						Num3 string `json:"3"`
						Num4 string `json:"4"`
						One5 string `json:"1.5"`
					} `json:"animated"`
					Static struct {
						Num1 string `json:"1"`
						Num2 string `json:"2"`
						Num3 string `json:"3"`
						Num4 string `json:"4"`
						One5 string `json:"1.5"`
					} `json:"static"`
				} `json:"dark"`
				Light struct {
					Animated struct {
						Num1 string `json:"1"`
						Num2 string `json:"2"`
						Num3 string `json:"3"`
						Num4 string `json:"4"`
						One5 string `json:"1.5"`
					} `json:"animated"`
					Static struct {
						Num1 string `json:"1"`
						Num2 string `json:"2"`
						Num3 string `json:"3"`
						Num4 string `json:"4"`
						One5 string `json:"1.5"`
					} `json:"static"`
				} `json:"light"`
			} `json:"images"`
		} `json:"tiers"`
		Backgrounds []string  `json:"backgrounds"`
		States      []string  `json:"states"`
		Type        string    `json:"type"`
		UpdatedAt   time.Time `json:"updated_at"`
		Priority    int       `json:"priority"`
	} `json:"actions"`
}

func (tk *TwitchKraken) GetGlobalCheeremotes() (resp *Cheeremotes, jsoerr *network.JsonError, err error) {
	resp = new(Cheeremotes)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/bits/actions", hMap, &resp)
	return
}

func (tk *TwitchKraken) GetCheeremotesForChannel(channelId string) (resp *Cheeremotes, jsoerr *network.JsonError, err error) {
	resp = new(Cheeremotes)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/bits/actions?channel="+channelId, hMap, &resp)
	return
}
