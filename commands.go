package ts3sqlib

import (
	"fmt"
)

//Quit closes the TS3 Server Query connection
func (c *SqConn) Quit() (err error) {
	_, err = c.Send("quit\n")

	return
}

//Login authenticates with the TS3 Server with a given username and password.
func (c *SqConn) Login(username, password string) (err error) {
	msg := fmt.Sprintf("login %s %s\n", Escape(username), Escape(password))

	_, err = c.Send(msg)

	return
}

//Use selects a virtual server specified with sid.
func (c *SqConn) Use(sid int) (err error) {
	msg := fmt.Sprintf("use %d\n", sid)

	_, err = c.Send(msg)

	return
}

//Servernotifyregister registers for a specified category of events.
//Categories are server, channel, textchannel or textprivate.
func (c *SqConn) Servernotifyregister(event string) (err error) {
	msg := fmt.Sprintf("servernotifyregister event=%s\n", event)

	_, err = c.Send(msg)

	return
}

//Sendtextmessage send a text message to a specified target.
func (c *SqConn) Sendtextmessage(targetmode, target int, msg string) (err error) {
	msg = Escape(msg)
	msg = fmt.Sprintf("sendtextmessage targetmode=%d target=%d msg=%s",
		targetmode, target, msg)

	_, err = c.Send(msg)

	return
}
