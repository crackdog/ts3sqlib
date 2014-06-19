package ts3sqlib

import (
	"strings"
)

func (c *SqConn) ClientlistToMaps() (clients []map[string]string, err error) {
	answer, err := c.Send("clientlist\n")
	if err != nil {
		return
	}

	clients, err = MsgToMaps(answer)
	return
}

//MsgToMaps converts a given ts3 serverquery answer into a slice of maps,
//with key-value-pairs with a '='.
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
				//TODO: maybe add to map...
				continue
			}
			parts[i][pair[0]] = Unescape(pair[1])
		}
	}

	return
}
