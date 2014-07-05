package ts3sqlib

import (
	"strconv"
	"strings"
)

//Client hold all information for a client from the clientlist command.
type Client struct {
	Cid              int    `json:"-"`
	Clid             int    `json:"clid"`
	ClientDatabaseId int    `json:"-"`
	ClientNickname   string `json:"client_nickname"`
	ClientType       int    `json:"-"`
}

//NewClient creates a Client datastructure from a map of strings
func NewClient(cmap map[string]string) Client {
	var newC Client

	newC.Cid, _ = strconv.Atoi(cmap["cid"])
	newC.Clid, _ = strconv.Atoi(cmap["clid"])
	newC.ClientDatabaseId, _ = strconv.Atoi(cmap["client_database_id"])
	newC.ClientNickname = cmap["client_nickname"]
	newC.ClientType, _ = strconv.Atoi(cmap["client_type"])

	return newC
}

//ClientmapsToClients converts an array of string maps to an array of Client's.
func ClientmapsToClients(clientmaps []map[string]string) (clients []Client, err error) {
	clients = make([]Client, len(clientmaps))

	for i, clientmap := range clientmaps {
		clients[i] = NewClient(clientmap)
	}

	return
}

//ClientlistToClients gets the clientlist from the ts3 server and returns it as
//a slice of Client's.
//The params are described in the TS3 ServerQuery Manual.
func (c *SqConn) ClientlistToClients(params string) (clients []Client, err error) {
	clientmaps, err := c.ClientlistToMaps(params)
	if err != nil {
		return
	}

	clients, err = ClientmapsToClients(clientmaps)
	return
}

//ClientlistToMaps gets the clientlist from the ts3 server and returns it as
//a slice of maps.
//The params are described in the TS3 ServerQuery Manual.
func (c *SqConn) ClientlistToMaps(params string) (clients []map[string]string, err error) {
	answer, err := c.Send("clientlist " + params + "\n")
	if err != nil {
		return
	}

	clients, err = MsgToMaps(answer)
	return
}

//MsgToMaps converts a given ts3 serverquery answer into a slice of maps,
//with key-value-pairs seperated by a '='.
func MsgToMaps(msg string) (parts []map[string]string, err error) {
	lines := strings.Split(msg, "|")
	parts = make([]map[string]string, len(lines))

	for i := range lines {
		parts[i] = make(map[string]string)

		lines[i] = strings.Replace(lines[i], "\n", "", -1)
		pairs := strings.Split(lines[i], " ")

		for j := range pairs {
			pair := strings.Split(pairs[j], "=")
			if len(pair) != 2 {
				continue
			}
			parts[i][pair[0]] = Unescape(pair[1])
		}
	}

	return
}
