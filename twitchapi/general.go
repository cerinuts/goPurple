package twitchapi

import (

)

type TwitchKraken struct{
	ClientID string
}

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goPurple/twitchapi", "0", "1", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

