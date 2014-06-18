package ts3sqlib

import (
	"strings"
)

type Stringmap map[string]string

func (c *SqConn) ClientlistToMaps() (clients []Stringmap, err error) {
	answer, err := c.Send("clientlist\n")
	if err != nil {
		return
	}

	tmpclients := strings.Split(answer, "|")
	clients = make([]Stringmap, len(tmpclients))

	for i := range tmpclients {
		clients[i] = make(map[string]string)

		tmpclients[i] = strings.Replace(tmpclients[i], "\n", "", -1)
		pairs := strings.Split(tmpclients[i], " ")

		for j := range pairs {
			pair := strings.Split(pairs[j], "=")
			if len(pair) != 2 {
				continue
			}
			clients[i][pair[0]] = Unescape(pair[1])
		}
	}

	return
}
