/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"code.cerinuts.io/libs/goBlue/network"
	"code.cerinuts.io/libs/goBlue/util"
)

var waitForToken chan string

//GetOauthToken starts an oauth grant process by starting a callback server on the callbackURL which will receive the token
//and returns the url which the user is required to open
func (tk *TwitchKraken) GetOauthToken(forceAuth bool, callbackURL, scopes string, callback func(token string)) (url string) {
	waitForToken = make(chan string, 1)
	state := util.GetRandomAlphanumericString(10)

	url = BaseURL + "/oauth2/authorize?" +
		"response_type=token" +
		"&client_id=" + tk.ClientID +
		"&redirect_uri=" + callbackURL + "/callback" +
		"&scope=" + scopes +
		"&state=" + state
	if forceAuth {
		url += "&force_verify=true"
	}
	srv := startTokenServer(callbackURL)
	go func(srv *http.Server) {
		result := <-waitForToken
		srv.Shutdown(context.Background())
		callback(result)
	}(srv)

	return url
}

func handleResponse(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(getHTTPResponse(r.Host)))
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query()["token"][0]
	token = strings.TrimSpace(token)
	if token == "" {
		w.Write([]byte(getHTTPResponseError()))
	} else {
		w.Write([]byte(getHTTPResponseSuccess()))
	}
	waitForToken <- token
}

func startTokenServer(callbackURL string) *http.Server {
	srv := &http.Server{Addr: strings.TrimPrefix(strings.TrimPrefix(callbackURL, "http://"), "https://")}

	http.HandleFunc("/token", handleToken)
	http.HandleFunc("/callback", handleResponse)

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func getHTTPResponse(callbackURL string) string {
	return `<!DOCTYPE html>
			<html><head>
			<title>OAuth Token</title>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
			<script type="text/javascript">
			
			function initiate() {	
				if((new URL(location)).searchParams.get("error") != null){
					document.location.replace("http://` + callbackURL + `/token?token=");
				}
			
				var hash = document.location.hash.substr(1);
				var hashParts = hash.split("&");	
				var access_token = null;
				
				for (var i=0;i<hashParts.length;++i) {		
					var parameterParts = hashParts[i].split("=");
					var key = parameterParts[0];		
					var value = parameterParts[1];
					
					if (key == "access_token") {			
					access_token = value;		
					}	
				}
				
				
				document.getElementById("javascript").className = "";	

				if(access_token != null){
					document.location.replace("http://` + callbackURL + `/token?token="+access_token);
				}
			}
				
				
				</script><style type="text/css">
				body { font-family: Consolas, sans-serif; text-align: center; background-color: #FFF; max-width: 500px; margin: auto; }
				input { font-family: Consolas, sans-serif; width: 300px; font-size: 1em; }
				noscript { color: red;  }.hide { display: none; }</style></head>
				<body onload="initiate()">
					<h1>OAuth token</h1>
					<noscript>
					    <p>This page requires <strong>JavaScript</strong> to get the token.</p>
					</noscript>				
					<p id="javascript" class="hide">You should be redirected..</p></body>
				</html>`
}

func getHTTPResponseSuccess() string {
	return `<!DOCTYPE html>
			<html><head>
			<title>OAuth Token</title>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
			
			<style type="text/css">
			body { font-family: Consolas, sans-serif; text-align: center; background-color: #FFF; max-width: 500px; margin: auto; }
			input { font-family: Consolas, sans-serif; width: 300px; font-size: 1em; }
			noscript { color: red;  }.hide { display: none; }</style></head>
				<body>
					<h1>OAuth token</h1>			
					<p>Success! You may close this page now.</p>
				</body>
			</html>`
}

func getHTTPResponseError() string {
	return `<!DOCTYPE html>
			<html><head>
			<title>OAuth Token</title>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
				
			<style type="text/css">
			body { font-family: Consolas, sans-serif; text-align: center; background-color: #FFF; max-width: 500px; margin: auto; }
			input { font-family: Consolas, sans-serif; width: 300px; font-size: 1em; }</style></head>
				<body>
					<h1>OAuth token</h1>						
					<p>Unfortunately something went wrong. Please restart the application</p>
				</body>
			</html>`
}

//RevokeToken revokes a previously granted oauth token
func (tk *TwitchKraken) RevokeToken(oauth string) {
	jac := new(network.JSONAPIClient)
	jac.Post(BaseURL+"/oauth2/revoke?client_id="+tk.ClientID+"&token="+oauth, nil, "", nil)
}

//TokenValidation is the twitch api response for token validation
type TokenValidation struct {
	Identified bool `json:"identified"`
	Token      struct {
		Valid         bool `json:"valid"`
		Authorization struct {
			Scopes    []string  `json:"scopes"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"authorization"`
		UserName string `json:"user_name"`
		ClientID string `json:"client_id"`
	} `json:"token"`
	Links struct {
		Channel  string `json:"channel"`
		Chat     string `json:"chat"`
		Teams    string `json:"teams"`
		User     string `json:"user"`
		Users    string `json:"users"`
		Streams  string `json:"streams"`
		Ingests  string `json:"ingests"`
		Channels string `json:"channels"`
	} `json:"_links"`
}

//ValidateToken sends an oauth token to twitch api for validation
func (tk *TwitchKraken) ValidateToken(oauth string) (resp *TokenValidation, jsoerr *network.JSONError, err error) {
	resp = new(TokenValidation)
	jac := new(network.JSONAPIClient)
	hMap := make(map[string]string)
	hMap["Accept"] = APIVersionHeader
	hMap["Client-ID"] = tk.ClientID
	hMap["Authorization"] = "OAuth " + oauth
	jsoerr, err = jac.Request(BaseURL, hMap, &resp)
	return
}
