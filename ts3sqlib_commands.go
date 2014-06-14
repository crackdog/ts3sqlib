package ts3sqlib

import (
	"fmt"
)

func (c *SqConn) Login(username, password string) (err error) {
	msg := fmt.Sprintf("login %s %s\n", Escape(username), Escape(password))

	_, err = c.Send(msg)

	return
}
