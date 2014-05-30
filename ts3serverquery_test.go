package ts3sqlib

import (
	"testing"
)

func TestEscape(t *testing.T) {
	var a, b, x string
	for i := range escapeSeq {
		a = "test" + string(escapeSeq[i].ascii) + "test"
		b = "test" + string(escapeSeq[i].escape) + "test"

		x = escape(a)
		if x != b {
			t.Errorf("escape('%s') = '%s', want '%s'", a, x, b)
		}

		x = unescape(b)
		if x != a && i > 0 {
			t.Errorf("unescape('%s') = '%s', want '%s', index=%d", b, x, a, i)
		}
	}
}
