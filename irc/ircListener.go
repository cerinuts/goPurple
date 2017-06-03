package irc

import (
	"github.com/ceriath/goBlue/archium"
)



type TwitchIRCListener struct{
	ArchiumDataIdentifier, ArchiumPrefix string
}

func (til *TwitchIRCListener) Trigger(ae archium.ArchiumEvent){
	val := ae.Data[til.ArchiumDataIdentifier]
	println(val)
}
	
func (til *TwitchIRCListener) GetType() string{
	return til.ArchiumPrefix + "*"
}

