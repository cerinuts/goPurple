/*
Copyright (c) 2018 ceriath
This Package is part of the "goPurple"-Library
It is licensed under the MIT License
*/

//Package twitchapi is used for twitch's API
package twitchapi

//TwitchKraken exposes the twitch kraken API (v5)
type TwitchKraken struct {
	ClientID string
}

//HiddenKraken exposes hidden, non-supported functions of the twitch kraken api (v5)
//use with caution, those might break more likely
type HiddenKraken struct {
	Tk *TwitchKraken
}

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goPurple/twitchapi", "0", "2", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//BaseURL of the twitch kraken API
const BaseURL = "https://api.twitch.tv/kraken"

//APIVersionHeader is the common header to set the selected api version (v5)
const APIVersionHeader = "application/vnd.twitchtv.v5+json"
