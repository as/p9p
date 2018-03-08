package p9p

import (
	"context"
	"log"
	"net"
	"testing"
	"time"
)

const realPlan9 = "1.1.1.1:808"

func testConn(t *testing.T) (client, server *Conn) {
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
	return bio0, rt.(*Conn)
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

func ckHasPlan9(t *testing.T) {
	t.Helper()
	if realPlan9 == "" {
		t.Skip("no real plan9 defined")
	}
	result := make(chan error)
	ctx, fn := context.WithTimeout(context.Background(), time.Second)
	defer fn()
	go func() {
		dialer := net.Dialer{}
		conn, err := dialer.DialContext(ctx, "tcp", realPlan9)
		result <- err
		if err == nil {
			conn.Close()
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-result:
		if err == nil {
			return
		}
	}

	t.Skip("plan9 isn't responding here:", realPlan9)

}

func TestVersionPlan9(t *testing.T) {
	ckHasPlan9(t)
	conn, err := Dial("tcp", realPlan9)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(conn.Version())
}
