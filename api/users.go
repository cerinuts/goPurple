package api

import (
	"strings"
"github.com/ceriath/goBlue/network"
)

type User struct {
	ID           string `json:"_id"`
	Bio          string `json:"bio"`
	Created_at   string `json:"created_at"`
	Display_name string `json:"display_name"`
	Logo         string `json:"logo"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Updated_at   string `json:"updated_at"`
}

type Users struct{
	Total int `json:"_total"`
	Users []User `json:"users"`
}

func (tk *TwitchKraken) GetByNames(names []string) (resp *Users, jsoerr *network.JsonError, err error) {
	resp = new(Users)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	
	logins := ""
	for _, val := range names{
		logins = logins + val + ","
	}
	logins = strings.TrimRight(logins, ",")
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/users?login="+logins, hMap, &resp)
	return
}

func (tk *TwitchKraken) GetByName(name string) (resp *Users, jsoerr *network.JsonError, err error) {
	resp = new(Users)
	jac := new(network.JsonApiClient)
	hMap := make(map[string]string)
	hMap["Accept"] = "application/vnd.twitchtv.v5+json"
	hMap["Client-ID"] = tk.ClientID
	jsoerr, err = jac.Request("https://api.twitch.tv/kraken/users?login="+name, hMap, &resp)
	return
}