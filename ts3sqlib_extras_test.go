package ts3sqlib

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMsgToMaps(t *testing.T) {
	testmsg := "client=test client2=test2 | client=test3 client2=test4\n"
	testmap := make([]map[string]string, 2)
	testmap[0] = make(map[string]string)
	testmap[1] = make(map[string]string)
	testmap[0]["client"] = "test"
	testmap[0]["client2"] = "test2"
	testmap[1]["client"] = "test3"
	testmap[1]["client2"] = "test4"

	xs, err := MsgToMaps(testmsg)

	if err != nil {
		t.Errorf("MsgToMaps(%s) gives error: '%s'", err.Error())
	} else {
		for i, x := range xs {
			if i < len(testmap) {
				if !reflect.DeepEqual(x, testmap[i]) {
					t.Errorf("MsgToMaps(%s) = '%s', want'%s'",
						fmt.Sprint(testmsg),
						fmt.Sprint(x),
						fmt.Sprint(testmap))
				}
			} else {
				t.Errorf("MsgToMaps dimension error")
			}
		}
	}
}
