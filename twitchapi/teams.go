package twitchapi

import ()

type Team struct {
	ID           int    `json:"_id"`
	Background   string `json:"background"`
	Banner       string `json:"banner"`
	Created_at   string `json:"created_at"`
	Display_name string `json:"display_name"`
	Info         string `json:"info"`
	Logo         string `json:"logo"`
	Name         string `json:"name"`
	Updated_at   string `json:"updated_at"`
}

type Teams struct {
	Teams []Team `json:"teams"`
}
