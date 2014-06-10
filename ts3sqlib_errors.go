package ts3sqlib

import (
	"fmt"
	"strings"
)

//Error contains additional error information.
type Error struct {
	Id       int
	Msg      string
	ExtraMsg string
}

//Error returns the error in a string representation.
func (err Error) Error() string {
	s := fmt.Sprintf("error id=%d msg=%s", err.Id, err.Msg)
	if err.ExtraMsg != "" {
		s += fmt.Sprintf(" extra_msg=%s", err.ExtraMsg)
	}
	return s
}

//NewError creates a new Error from an id, message and an extra_message.
func NewError(id int, msg, extramsg string) Error {
	return Error{id, msg, extramsg}
}

//isError tests if a given string is a ts3 server query error
func isError(line string) bool {
	return strings.Contains(line, "error") &&
		strings.Contains(line, "id=") && strings.Contains(line, "msg=")
}

//toError converts a given string into a ts3 server query error.
func toError(line string) (err Error) {
	if !isError(line) {
		err = Error{666, "line is not an error!", ""}
		return
	}
	err.Id = 0
	err.Msg = "msg"
	err.ExtraMsg = "extra_msg"
	return
}
