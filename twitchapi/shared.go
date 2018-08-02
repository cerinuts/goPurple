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

//AppName is the name of the application
const AppName string = "goPurple/twitchapi"

//VersionMajor 0 means in development, >1 ensures compatibility with each minor version, but breakes with new major version
const VersionMajor string = "0"

//VersionMinor introduces changes that require a new version number. If the major version is 0, they are likely to break compatibility
const VersionMinor string = "2"

//VersionBuild is the type of this release. s(table), b(eta), d(evelopment), n(ightly)
const VersionBuild string = "s"

//FullVersion contains the full name and version of this package in a printable string
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//BaseURL of the twitch kraken API
const BaseURL = "https://api.twitch.tv/kraken"

//APIVersionHeader is the common header to set the selected api version (v5)
const APIVersionHeader = "application/vnd.twitchtv.v5+json"
