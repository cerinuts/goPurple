package twitchapi

import (
	"encoding/json"
	"time"
)

type Video struct {
	ID            json.Number `json:"_id"`
	BroadcastID   json.Number `json:"broadcast_id"`
	BroadcastType string      `json:"broadcast_type"`
	Channel       struct {
		ID          string `json:"_id"`
		DisplayName string `json:"display_name"`
		Name        string `json:"name"`
	} `json:"channel"`
	CreatedAt       time.Time `json:"created_at"`
	Description     string    `json:"description"`
	DescriptionHTML string    `json:"description_html"`
	Fps             struct {
		Chunked float64 `json:"chunked"`
		High    float64 `json:"high"`
		Low     float64 `json:"low"`
		Medium  float64 `json:"medium"`
		Mobile  float64 `json:"mobile"`
	} `json:"fps"`
	Game     string `json:"game"`
	Language string `json:"language"`
	Length   int    `json:"length"`
	Preview  struct {
		Large    string `json:"large"`
		Medium   string `json:"medium"`
		Small    string `json:"small"`
		Template string `json:"template"`
	} `json:"preview"`
	PublishedAt time.Time `json:"published_at"`
	Resolutions struct {
		Chunked string `json:"chunked"`
		High    string `json:"high"`
		Low     string `json:"low"`
		Medium  string `json:"medium"`
		Mobile  string `json:"mobile"`
	} `json:"resolutions"`
	Status     string `json:"status"`
	TagList    string `json:"tag_list"`
	Thumbnails struct {
		Large []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"large"`
		Medium []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"medium"`
		Small []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"small"`
		Template []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"template"`
	} `json:"thumbnails"`
	Title      string      `json:"title"`
	URL        string      `json:"url"`
	Viewable   string      `json:"viewable"`
	ViewableAt interface{} `json:"viewable_at"`
	Views      int         `json:"views"`
}

type Videos struct {
	Total  int     `json:"_total"`
	Videos []Video `json:"videos"`
}

