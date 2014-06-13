package ts3sqlib

import (
	"testing"
)

func TestToError(t *testing.T) {
	line1 := "error id=2 msg=hallo extra_msg=welt"
	line2 := "error id=404 msg=notfound"
	line3 := "test"

	err1 := NewError(2, "hallo", "welt")
	err2 := NewError(404, "notfound", "")
	err3 := NewError(666, "line is not an error!", "")

	if x := toError(line1); x != err1 {
		t.Errorf("toError(%s) = '%s', want '%s'", line1, x.Error(), err1.Error())
	}

	if x := toError(line2); x != err2 {
		t.Errorf("toError(%s) = '%s', want '%s'", line2, x.Error(), err2.Error())
	}

	if x := toError(line3); x != err3 {
		t.Errorf("toError(%s) = '%s', want '%s'", line2, x.Error(), err2.Error())
	}
}

func TestEqual(t *testing.T) {
	line1 := "error id=0 msg=ok"

	if x := MsgEndError.equals(toError(line1)); !x {
		t.Errorf("equals error")
	}
}
