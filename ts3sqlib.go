package ts3sqlib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	stdPort string = ":10011"
)

//SqConn contains the connection to a ts3 server.
type SqConn struct {
	conn       net.Conn
	logger     *log.Logger
	sendMutex  *sync.Mutex
	receiving  bool
	recvNotify chan string
	recvChan   chan string
	closed     chan bool

	WelcomeMsg string
}

//Dial creates a new SqConn and connects to the ts3 server with the given
//address and returns a pointer to it and an error.
//If the logger is nil, then the standard logger is used.
func Dial(address string, logger *log.Logger) (conn *SqConn, err error) {
	conn = nil

	if !strings.Contains(address, ":") {
		address += stdPort
	}

	connection, err := net.Dial("tcp", address)
	if err != nil {
		return
	}

	/*if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}*/

	conn = &SqConn{
		conn:       connection,
		logger:     logger,
		sendMutex:  &sync.Mutex{},
		receiving:  true,
		recvNotify: make(chan string),
		recvChan:   make(chan string),
		closed:     make(chan bool),
		WelcomeMsg: "",
	}

	go conn.recv() //goroutine that splits the incoming messages into notify
	//and normal messages.

	if !strings.Contains(<-conn.recvChan, "TS3") {
		conn = nil
		err = fmt.Errorf("no connection to a ts3-server-query")
		return
	}

	conn.WelcomeMsg = <-conn.recvChan

	return
}

//TODO: remove this function...
func (c *SqConn) RecvTest() {
	c.logger = log.New(os.Stderr, "", log.LstdFlags)
	_, _ = c.Send("use 1\n")
	/*c.logger.Println(answer)
	if err != nil {
		c.logger.Println(err)
	}*/

	_, _ = c.Send("quit\n")

	_, _ = c.Send("help login\n")

	c.Close()
}

func (c *SqConn) recv() {
	var err error
	line := ""
	scan := bufio.NewScanner(c.conn)

	for c.receiving {
		//read line
		if !scan.Scan() {
			if err = scan.Err(); err != nil && c.logger != nil {
				c.logger.Print(err)
			}
			c.receiving = false
			c.closed <- true
			break
		}

		line = scan.Text()
		line = strings.Replace(line, string([]byte{13}), "", -1)

		if err = scan.Err(); err != nil {
			if c.logger != nil {
				c.logger.Print(err)
			}
			continue
		}

		//decide if its notify or not and put it to the correct channel
		if strings.HasPrefix(line, "notify") {
			c.recvNotify <- line
		} else {
			c.recvChan <- line
		}
	}
}

//RecvNotify returns a notify message or blocks until one arrives.
func (c *SqConn) RecvNotify() (answer string, err error) {
	answer = ""
	err = nil

	answer = <-c.recvNotify
	return
}

//Close closes the connection to the ts3 server.
func (c *SqConn) Close() {
	c.receiving = false
	c.conn.Close()
}

//Send sends a message to the server and returns the answer and an error.
func (c *SqConn) Send(msg string) (answer string, err error) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()

	if !c.receiving {
		answer = ""
		err = fmt.Errorf("connection closed")
		return
	}

	if c.logger != nil {
		c.logger.Print("send: ", msg)
	}

	answer = ""

	_, err = c.conn.Write([]byte(msg)) //If the return value is smaller than
	if err != nil {                    //the length of msg, it's an error.
		return
	}

	line := ""
	err = nil

	for !isError(line) && c.receiving {
		select {
		case line = <-c.recvChan:
			if isError(line) {
				err = toError(line)
				break
			}

			answer += line + "\n"

		case <-c.closed:
			break
		}
	}

	if err == nil {
		err = NewError(-1, "connection closed", "")
	}

	if MsgEndError.equals(err) {
		err = nil
	}

	//logging
	if c.logger != nil {
		if answer != "" && answer != "\n" {
			c.logger.Println(answer)
		}
		if err != nil {
			c.logger.Println(err)
		}
	}

	return
}
