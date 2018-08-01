/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package gql is used for twitch's GraphQL
package gql

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"code.cerinuts.io/libs/goBlue/log"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goPurple/gql", "0", "2", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

const gqlURL = "https://gql.twitch.tv/gql"

//Client is a client for twitch's gql
type Client struct {
	ClientID string
	OAuth    string
}

//CreateVideoBookmarkInputResponse is the response sent by gql for creating video bookmarks
type CreateVideoBookmarkInputResponse struct {
	Data struct {
		CreateVideoBookmark struct {
			Error         interface{} `json:"error"`
			VideoBookmark struct {
				PositionSeconds int `json:"positionSeconds"`
			} `json:"videoBookmark"`
		} `json:"createVideoBookmark"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int `json:"durationMilliseconds"`
	} `json:"extensions"`
}

//CreateVideoBookmarkInput creates a video bookmark (marker)
func (gqlc *Client) CreateVideoBookmarkInput(description, streamID string) (response *CreateVideoBookmarkInputResponse, err error) {

	body := `{
		"query": "mutation($input: CreateVideoBookmarkInput!) {\n\tcreateVideoBookmark(input: $input) {\n\t\terror {\n\t\t\tcode\n\t\t}\n\t\tvideoBookmark {\n\t\t\tpositionSeconds\n\t\t}\n\t}\n}",
		"variables": {
			"input": {
				"broadcastID": "` + streamID + `",
				"description": "` + description + `",
				"medium": "chat",
				"platform": "web"
			}
		}
	}`

	log.D(body)

	header := make(map[string]string)
	header["Client-ID"] = gqlc.ClientID
	header["Authorization"] = "OAuth " + gqlc.OAuth
	header["Content-Type"] = "application/json"

	response = new(CreateVideoBookmarkInputResponse)

	req, err := http.NewRequest(http.MethodPost, gqlURL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.E(err)
		return nil, err
	}

	err = runRequest(req, header, response)
	return
}

//runRequest actually runs the request prepared by functions above
func runRequest(req *http.Request, header map[string]string, response interface{}) error {
	cli := &http.Client{
		Timeout: time.Second * 10,
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	res, getErr := cli.Do(req)
	if getErr != nil {
		log.E(getErr)
		return getErr
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return nil
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.E(readErr)
		return readErr
	}

	if res.StatusCode == 200 {
		return json.Unmarshal(body, &response)
	}
	return errors.New(strconv.Itoa(res.StatusCode))

}
