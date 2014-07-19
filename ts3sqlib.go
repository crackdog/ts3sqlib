//Package ts3sqlib provides a simple interface for the TeamSpeak3 ServerQuery.
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

var (
	//StdoutLogger is the standard logger that prints to stdout.
	StdoutLogger = log.New(os.Stdout, "", log.LstdFlags)
	//StderrLogger is the standard logger that prints to stderr.
	StderrLogger = log.New(os.Stderr, "", log.LstdFlags)
)

//SqConn contains the connection to a ts3 server.
type SqConn struct {
	conn           net.Conn
	logger         *log.Logger
	sendMutex      *sync.Mutex
	receiving      bool
	recvNotify     chan string
	notifyChannels []chan string
	recvChan       chan string
	closed         chan bool
	//WelcomeMsg contains the welcome message from the ts3 server
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

	conn = &SqConn{
		conn:           connection,
		logger:         logger,
		sendMutex:      &sync.Mutex{},
		receiving:      true,
		recvNotify:     make(chan string),
		notifyChannels: make([]chan string, 0),
		recvChan:       make(chan string),
		closed:         make(chan bool),
		WelcomeMsg:     "",
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
			//send to notifyChannels
			for _, nc := range c.notifyChannels {
				nc <- line
			}
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

func (c *SqConn) Notify() <-chan string {
	s := make(chan string)
	//adding s to notify channels
	c.notifyChannels = append(c.notifyChannels, s)

	return s
}

//Close closes the connection to the ts3 server.
func (c *SqConn) Close() {
	c.receiving = false
	c.conn.Close()
}

//IsClosed checks if the connection is closed
func (c *SqConn) IsClosed() bool {
	return !c.receiving
}

//Send sends a message to the server and returns the answer and an error.
func (c *SqConn) Send(msg string) (answer string, err error) {
	if c == nil {
		fmt.Errorf("nil pointer error")
		return
	}

	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()

	//msg must have a newline at the end
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

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

	if MsgEndError.Equals(err) {
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
