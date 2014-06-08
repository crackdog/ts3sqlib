package ts3sqlib

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	timeout = 5 * time.Second
)

type SqConn struct {
	conn      net.Conn
	logger    *log.Logger
	sendMutex *sync.Mutex
	timeout   time.Duration
}

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
		conn:      connection,
		logger:    logger,
		sendMutex: &sync.Mutex{},
		timeout:   1000 * time.Millisecond,
	}
	return
}

func (c *SqConn) Close() {
	c.conn.Close()
}

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
		line, err = bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}
		answer += line
	}
	err = toError(line)

	return
}
