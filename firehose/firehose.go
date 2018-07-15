package firehose

import (
	"gitlab.ceriath.net/libs/goBlue/network"
)

type Firehose struct {
	client *network.EventsourceClient
}

func (f *Firehose) Connect(token string) (events chan network.Event, err error) {
	f.client = new(network.EventsourceClient)
	stream, err := f.client.Subscribe("https://tmi.twitch.tv/firehose?oauth_token=" + token)
	if err != nil {
		return nil, err
	}

	return stream.EventQueue, nil
}

func (f *Firehose) Disconnect() {
	f.client.Close()
}
