package ts3sqlib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	timeout = 5 * time.Second //the timeout of the connection
)

//SqConn contains the connection to a ts3 server.
type SqConn struct {
	conn       net.Conn
	logger     *log.Logger
	sendMutex  *sync.Mutex
	timeout    time.Duration
	receiving  bool
	recvNotify chan string
	recvChan   chan string
}

//Dial creates a new SqConn and connects to the ts3 server with the given
//address and returns a pointer to it and an error.
func Dial(address string, logger *log.Logger) (conn *SqConn, err error) {
	conn = nil

	if !strings.Contains(address, ":") {
		address += ":9987"
	}

	connection, err := net.Dial("tcp", address)
	if err != nil {
		return
	}

	conn = &SqConn{
		conn:       connection,
		logger:     logger,
		sendMutex:  &sync.Mutex{},
		timeout:    1000 * time.Millisecond,
		receiving:  true,
		recvNotify: make(chan string),
		recvChan:   make(chan string),
	}

	go conn.recv() //goroutine that splits the incoming messages into notify
	//and normal messages.

	return
}

func (c *SqConn) recv() {
	line := ""
	var err error
	for c.receiving {
		//read line
		line, err = bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			panic(err)
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

	if c == nil {
		err = fmt.Errorf("no SqConn")
		return
	}

	answer = <-c.recvNotify
	return
}

//Close closes the connection to the ts3 server.
func (c *SqConn) Close() {
	c.conn.Close()
}

//Send sends a message to the server and returns the answer.
func (c *SqConn) Send(msg string) (answer string, err error) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()

	answer = ""

	_, err = c.conn.Write([]byte(msg)) //If the return value is smaller than
	if err != nil {                    //the length of msg, it's an error.
		return
	}

	//wait for answer...
	err = c.conn.SetReadDeadline(time.Now().Add(timeout))

	line := ""

	for !isError(line) {
		/*line, err = bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}*/
		answer += line
		line = <-c.recvChan
	}
	err = toError(line)

	return
}
