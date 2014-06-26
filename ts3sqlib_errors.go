package ts3sqlib

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	//MsgEndError is the normal Error at the end of each message.
	MsgEndError = NewError(0, "ok", "")
	//ClosedError is the Error of a closed connection.
	ClosedError = NewError(-1, "connection closed", "")
	//InvalidLoginError is the Error for an invalid loginname or password.
	InvalidLoginError = NewError(520, "invalid loginname or password", "")
	//PermissionError is the Error for a lack of permissions.
	PermissionError = NewError(27, "insufficient client permissions", "")
)

//Error contains additional error information.
type Error struct {
	ID       int
	Msg      string
	ExtraMsg string
}

//Error returns the error in a string representation.
func (err Error) Error() string {
	s := fmt.Sprintf("error id=%d msg=%s", err.ID, err.Msg)
	if err.ExtraMsg != "" {
		s += fmt.Sprintf(" extra_msg=%s", err.ExtraMsg)
	}
	return s
}

//NewError creates a new Error from an id, message and an extra_message.
func NewError(id int, msg, extramsg string) Error {
	return Error{id, msg, extramsg}
}

//Equals compares an Error with another Error
func (err Error) Equals(compareErr error) bool {
	if err2, ok := compareErr.(Error); ok {
		return err.ID == err2.ID
	}
	return false
}

//isError tests if a given string is a ts3 server query error
func isError(line string) bool {
	return strings.HasPrefix(line, "error") &&
		strings.Contains(line, "id=") && strings.Contains(line, "msg=")
}

//toError converts a given string into a ts3 server query error.
func toError(line string) (err Error) {
	if !isError(line) {
		err = NewError(666, "line is not an error!", "")
		return
	}

	parts := strings.Split(line, " ")
	for i := range parts {
		if strings.Contains(parts[i], "=") {
			key, value, err2 := splitAtEqual(parts[i])
			if err2 {
				continue
			}

			switch {
			case strings.Contains(key, "id"):
				var err3 error
				err.ID, err3 = strconv.Atoi(value)
				if err3 != nil {
					err.ID = 999
				}
			case strings.Contains(key, "extra_msg"):
				err.ExtraMsg = Unescape(value)
			case strings.Contains(key, "msg"):
				err.Msg = Unescape(value)
			}
		}
	}
	return
}

func splitAtEqual(s string) (key string, value string, err bool) {
	if strings.Contains(s, "=") {
		tmp := strings.Split(s, "=")
		//fmt.Println(" -> ", s, " => ", tmp)

		key = tmp[0]
		value = tmp[1]

		err = false
	} else {
		err = true
	}
	return
}
