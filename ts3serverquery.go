//ts3sq provides a library for the ts3 server query interface.
package ts3sqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	InvalidLoginError = NewError(520, "invalid loginname or password", "")
)

//Error contains an Error Message from the server
type Error struct {
	Id       int    //Error Id
	Msg      string //Error Message
	ExtraMsg string //An additional Error Message
}

//Error returns an Error String
func (err Error) Error() string {
	if err.ExtraMsg == "" {
		return fmt.Sprintf("error id=%d msg=%s", err.Id, err.Msg)
	} else {
		return fmt.Sprintf("error id=%d msg=%s extra_msg=%s", err.Id, err.Msg, err.ExtraMsg)
	}
}

//NewError creates a new Error.
func NewError(id int, msg, extraMsg string) Error {
	return Error{id, msg, extraMsg}
}

//Ts3sqs contains a connection to a ts3server.
type Ts3sqs struct {
	serverconn net.Conn
	log        *log.Logger
	WelcomeMsg string
}

//New create a new ts3serverquery connection.
func New(address string, logger *log.Logger) (t *Ts3sqs, err error) {
	c, err := net.Dial("tcp", address)
	if err == nil {
		t := new(Ts3sqs)
		t.serverconn = c
		if logger != nil {
			t.log = logger
		} else {
			t.log = log.New(os.Stdout, "", log.LstdFlags)
		}
		t.log.Print("connected to ts3server")
		t.WelcomeMsg, err = t.WaitForMessageLine()
		if err != nil {
			c.Close()
			return nil, err
		}
		s, err := t.WaitForMessageLine()
		if err != nil {
			c.Close()
			return nil, err
		}
		t.WelcomeMsg += s
		fmt.Println(t.WelcomeMsg) //logging...
		return t, nil
	} else {
		return nil, err
	}
}

//Close closes the connection to the ts3 server.
func (s *Ts3sqs) Close() {
	s.serverconn.Close()
}

func (s *Ts3sqs) send(msg string) error {
	length, err := s.serverconn.Write([]byte(msg))
	if err == nil && length < len(msg) {
		return fmt.Errorf("only %d of %d bytes were sended.", length, len(msg))
	} else {
		return err
	}
}

func r2s(runes ...byte) []byte {
	return []byte(runes)
}

type escapeSequence struct {
	ascii  []byte
	escape []byte
}

var escapeSeq [10]escapeSequence = [...]escapeSequence{
	{r2s(92), r2s(92, 92)},
	{r2s(47), r2s(92, 47)},
	{r2s(32), r2s(92, 115)},
	{r2s(124), r2s(92, 112)},
	{r2s(7), r2s(92, 97)},
	{r2s(8), r2s(92, 98)},
	{r2s(12), r2s(92, 102)},
	{r2s(10), r2s(92, 110)},
	{r2s(13), r2s(92, 114)},
	{r2s(9), r2s(92, 116)},
}

func escape(s string) string {
	for i := range escapeSeq {
		s = string(bytes.Replace([]byte(s), escapeSeq[i].ascii, escapeSeq[i].escape, -1))
	}
	return s
}

func unescape(s string) string {
	//for i := range escapeSeq {
	for i := 1; i < len(escapeSeq); i++ {
		s = string(bytes.Replace([]byte(s), escapeSeq[i].escape, escapeSeq[i].ascii, -1))
	}
	//s = string(bytes.Replace([]byte(s), r2s(92, 92), []byte("\\"), -1))
	return s
	/*s, err := strconv.Unquote(s)
	if err != nil {
		fmt.Println(err)
	}
	return s*/
}

func (s *Ts3sqs) WaitForMessageLine() (string, error) {
	//WaitForMessageLine reads a line from the server connection.
	return bufio.NewReader(s.serverconn).ReadString('\n')
}

func (s *Ts3sqs) getError() error {
	//getError returns the error message from the server.
	msg, err := s.WaitForMessageLine()
	if err != nil {
		return err
	} else {
		if strings.Contains(msg, "error id=0 msg=ok") {
			return nil
		} else {
			//return fmt.Errorf("msg error: '%s'", unescape(strings.TrimSpace(msg)))
			return NewError(0, unescape(strings.TrimSpace(msg)), "")
		}
	}
}

func (s *Ts3sqs) sendWithGettingError(msg string) error {
	err := s.send(msg)
	if err != nil {
		return err
	} else {
		return s.getError()
	}
}

func (s *Ts3sqs) Login(username, password string) error {
	//logging in...
	username = escape(username)
	password = escape(password)
	msg := fmt.Sprintf("login client_login_name=%s client_login_password=%s\n", username, password)
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Logout() error {
	//logging out
	return s.sendWithGettingError("logout\n")
}

func (s *Ts3sqs) Clientlist() (string, error) {
	//Clientlist sends a clientlist request to the ts3 server.
	panic("Clientlist() is not implemented yet")
}

func (s *Ts3sqs) Use(server_id int) error {
	//Use sends a request to use a server.
	msg := fmt.Sprintf("use sid=%d\n", server_id)
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Servernotifyregister(event string) error {
	//Servernotifyregister sends a notify request for a given event.
	msg := fmt.Sprintf("servernotifyregister event=%s", escape(event))
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Servernotifyunregister(event string) error {
	//Servernotifyunregister sends a unnotify request for a given event.
	msg := fmt.Sprintf("servernotitfyunregister event=%s", escape(event))
	return s.sendWithGettingError(msg)
}

func (s *Ts3sqs) Sendtextmessage(targetmode, target int, raw_msg string) error {
	msg := fmt.Sprintf("sendtextmessage targetmode=%d target=%d msg=%s",
		targetmode, target, escape(raw_msg))
	return s.sendWithGettingError(msg)
}
