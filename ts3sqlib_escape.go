package ts3sqlib

import (
	"bytes"
)

type pair struct {
	a []byte
	b []byte
}

var (
	escapechars = [...]pair{
		pair{[]byte{92}, []byte{92, 92}},
		pair{[]byte{47}, []byte{92, 47}},
		pair{[]byte{32}, []byte{92, 115}},
		pair{[]byte{124}, []byte{92, 112}},
		pair{[]byte{7}, []byte{92, 97}},
		pair{[]byte{8}, []byte{92, 98}},
		pair{[]byte{12}, []byte{92, 102}},
		pair{[]byte{10}, []byte{92, 110}},
		pair{[]byte{13}, []byte{92, 114}},
		pair{[]byte{9}, []byte{92, 116}},
		pair{[]byte{11}, []byte{92, 118}},
	}
)

//Escape escapes a given string as described in the ts3 server query manual.
func Escape(msg string) string {
	//TODO escape the msg
	bytemsg := []byte(msg)

	for i := range escapechars {
		bytemsg = bytes.Replace(bytemsg, escapechars[i].a, escapechars[i].b, -1)
	}

	return string(bytemsg)
}

//Unescape unescapes a given escaped string as described in the ts3 server
//query manual.
func Unescape(escaped_msg string) string {
	//TODO unescape a escaped message
	bytemsg := []byte(escaped_msg)

	for i := range escapechars {
		bytemsg = bytes.Replace(bytemsg, escapechars[i].b, escapechars[i].a, -1)
	}

	return string(bytemsg)
}
