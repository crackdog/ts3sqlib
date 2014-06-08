package ts3sqlib

import (
	"fmt"
	"strings"
)

type Error struct {
	Id       int
	Msg      string
	ExtraMsg string
}

func (err Error) Error() string {
	s := fmt.Sprintf("error id=%d msg=%s", err.Id, err.Msg)
	if err.ExtraMsg != "" {
		s += fmt.Sprintf(" extra_msg=%s", err.ExtraMsg)
	}
	return s
}

func NewError(id int, msg, extramsg string) Error {
	return Error{id, msg, extramsg}
}

func isError(line string) bool {
	return strings.Contains(line, "error") &&
		strings.Contains(line, "id=") && strings.Contains(line, "msg=")
}

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
