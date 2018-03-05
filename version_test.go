package p9p

import (
	"net"
	"testing"
)

func testConn(t *testing.T) (client, server *Bio) {
	t.Helper()
	fd, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	wait := make(chan interface{})
	go func() {
		bio, err := Accept(fd)
		if err != nil {
			wait <- err
			return
		}
		wait <- bio
		close(wait)
	}()
	bio0, err := Dial("tcp", fd.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	rt := <-wait
	switch rt.(type) {
	case error:
		t.Fatal(err)
		return nil, nil
	}
	return bio0, rt.(*Bio)
}

func TestConn(t *testing.T) {
	testConn(t)
}
func TestVersion(t *testing.T) {
	c, s := testConn(t)
	cv, err := c.Version()
	if err != nil {
		t.Fatalf("client version error: %s", err)
	}
	sv, err := s.Version()
	if err != nil {
		t.Fatalf("server version error: %s", err)
	}
	if cv != sv {
		t.Fatalf("client and server differ: %q vs %q\n", cv, sv)
	}

}
