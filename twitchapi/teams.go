/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

//Team contains information about a team
type Team struct {
	ID          int    `json:"_id"`
	Background  string `json:"background"`
	Banner      string `json:"banner"`
	CreatedAt   string `json:"created_at"`
	DisplayName string `json:"display_name"`
	Info        string `json:"info"`
	Logo        string `json:"logo"`
	Name        string `json:"name"`
	UpdatedAt   string `json:"updated_at"`
}

//Teams is a list of teams
type Teams struct {
	Teams []Team `json:"teams"`
}
