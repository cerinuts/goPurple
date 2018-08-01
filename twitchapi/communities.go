/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

import (
	"encoding/json"
)

//Community contains information about a community
type Community struct {
	ID              json.Number `json:"_id"`
	AvatarImageURL  string      `json:"avatar_image_url"`
	CoverImageURL   string      `json:"cover_image_url"`
	Description     string      `json:"description"`
	DescriptionHTML string      `json:"description_html"`
	Language        string      `json:"language"`
	Name            string      `json:"name"`
	OwnerID         string      `json:"owner_id"`
	Rules           string      `json:"rules"`
	RulesHTML       string      `json:"rules_html"`
	Summary         string      `json:"summary"`
}

//Communities is a list of communities
type Communities struct {
	Communities []Community `json:"communities"`
}
