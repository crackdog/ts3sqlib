package ts3sqlib

import "testing"

func TestEscape(t *testing.T) {
	const a, b = "Hello World !", "Hello\\sWorld\\s!"
	if x := escape(a); x != b {
		t.Errorf("escape('%s') = '%s', want '%s'", a, x, b)
	}
}
