package ts3sqlib

import (
	"fmt"
	"strconv"
	"strings"
)

//Client hold all information for a client from the clientlist command.
type Client struct {
	Cid                     int    `json:"-"`
	Clid                    int    `json:"clid"`
	ClientDatabaseID        int    `json:"-"`
	ClientNickname          string `json:"client_nickname"`
	ClientType              int    `json:"client_type"`
	ConnectionConnectedTime int    `json:"connection_connected_time"`
}

//NewClient creates a Client datastructure from a map of strings
func NewClient(cmap map[string]string) Client {
	var newC Client

	newC.Cid, _ = strconv.Atoi(cmap["cid"])
	newC.Clid, _ = strconv.Atoi(cmap["clid"])
	newC.ClientDatabaseID, _ = strconv.Atoi(cmap["client_database_id"])
	newC.ClientNickname = cmap["client_nickname"]
	newC.ClientType, _ = strconv.Atoi(cmap["client_type"])

	newC.ConnectionConnectedTime = 0

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
		parts[i], err = MsgToMap(lines[i])
		if err != nil {
			return
		}
	}

	return
}

//MsgToMap converts a given ts3 serverquery answer into a map of
//key-value-pairs seperated by a '='.
func MsgToMap(msg string) (part map[string]string, err error) {
	part = make(map[string]string)

	msg = strings.Replace(msg, "\n", "", -1)
	pairs := strings.Split(msg, " ")

	for j := range pairs {
		pair := strings.Split(pairs[j], "=")
		if len(pair) != 2 {
			continue
		}
		part[pair[0]] = Unescape(pair[1])
	}

	return
}

//SendToMap combines a Send and a MsgToMap.
func (c *SqConn) SendToMap(msg string) (pairs map[string]string, err error) {
	answer, err := c.Send(msg)
	if err != nil {
		return
	}

	pairs, err = MsgToMap(answer)

	return
}

//SendToMaps combines a Send and a MsgToMaps.
func (c *SqConn) SendToMaps(msg string) (parts []map[string]string, err error) {
	answer, err := c.Send(msg)
	if err != nil {
		return
	}

	parts, err = MsgToMaps(answer)

	return
}

func (c *SqConn) GetConnectionTimeForCL(clientlist []Client) (clients []Client, err error) {
	for i := range clientlist {
		msg := fmt.Sprint("clientinfo clid=", clients[i].Clid)
		answer, err := c.SendToMap(msg)
		if err != nil {
			break
		}
		fmt.Println(msg)
		s := answer["connection_connected_time"]
		clients[i].ConnectionConnectedTime, _ = strconv.Atoi(s)
	}
	return clientlist, err
}
