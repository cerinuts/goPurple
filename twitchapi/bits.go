/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

import (
	"time"

	"code.cerinuts.io/libs/goBlue/network"
)

//Cheeremotes struct
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

//GetGlobalCheeremotes returns all cheeremotes that are accesible from every channel
func (tk *TwitchKraken) GetGlobalCheeremotes() (resp *Cheeremotes, jsoerr *network.JSONError, err error) {
	resp = new(Cheeremotes)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/bits/actions", hMap, &resp)
	return
}

//GetCheeremotesForChannel returns channel specific cheeremotes
func (tk *TwitchKraken) GetCheeremotesForChannel(channelID string) (resp *Cheeremotes, jsoerr *network.JSONError, err error) {
	resp = new(Cheeremotes)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request(BaseURL+"/bits/actions?channel="+channelID, hMap, &resp)
	return
}
