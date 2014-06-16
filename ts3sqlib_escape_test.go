package ts3sqlib

import (
	"testing"
)

func TestEscape(t *testing.T) {
	a := "hallo welt"
	ae := "hallo\\swelt"

	if x := Escape(a); x != ae {
		t.Errorf("Escape(%s) = '%s', want '%s'", a, x, ae)
	}
}

func TestUnescape(t *testing.T) {
	ae := "hallo\\swelt"
	a := "hallo welt"

	if x := Unescape(ae); x != a {
		t.Errorf("Unescape(%s) = '%s', want '%s'", ae, x, a)
	}
}
