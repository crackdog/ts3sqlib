package ts3sqlib

type Clientlist []Client
type Channellist []Channel

type Channel struct {
	Name    string            `json:"channel_name"`
	Data    map[string]string `json:"-"`
	Clients []Client          `json:"clients"`
}
