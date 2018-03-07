package p9p

import "testing"

func TestString(t *testing.T) {
	c, s := testConn(t)
		c.writestring("hello")
		c.Flush()
	s.readstring()
	if s.String() != "hello" {
		t.Fatalf("have %s want %q err %s", s.String(), "hello", s.err)
	}
}
