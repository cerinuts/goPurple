package gql

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gitlab.ceriath.net/libs/goBlue/log"
)

const GQLURL = "https://gql.twitch.tv/gql"

type GQLClient struct {
	ClientId string
	OAuth    string
}

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

func (gqlc *GQLClient) CreateVideoBookmarkInput(description, streamId string) (response *CreateVideoBookmarkInputResponse, err error) {

	body := `{
		"query": "mutation($input: CreateVideoBookmarkInput!) {\n\tcreateVideoBookmark(input: $input) {\n\t\terror {\n\t\t\tcode\n\t\t}\n\t\tvideoBookmark {\n\t\t\tpositionSeconds\n\t\t}\n\t}\n}",
		"variables": {
			"input": {
				"broadcastID": "` + streamId + `",
				"description": "` + description + `",
				"medium": "chat",
				"platform": "web"
			}
		}
	}`

	log.D(body)

	header := make(map[string]string)
	header["Client-ID"] = gqlc.ClientId
	header["Authorization"] = "OAuth " + gqlc.OAuth
	header["Content-Type"] = "application/json"

	response = new(CreateVideoBookmarkInputResponse)

	req, err := http.NewRequest(http.MethodPost, GQLURL, bytes.NewBuffer([]byte(body)))
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
	} else {
		return errors.New(strconv.Itoa(res.StatusCode))
	}
}
