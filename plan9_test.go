package p9p

import (
	"context"
	"net"
	"testing"
	"time"
)

const realPlan9 = "1.1.1.1:808"

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

// The server should run the equivalent of
//	aux/listen1 -t -v tcp!*!808 exportfs -r /
//
func TestPlan9Version(t *testing.T) {
	ckHasPlan9(t)
	conn, err := Dial("tcp", realPlan9)
	if err != nil {
		t.Fatal(err)
	}
	max, version, err := conn.Ver()
	max = max
	version = version
	if err != nil {
		t.Fatal(err)
	}
	err = conn.Attach(1, 2, "/root", "xxxxxxxxx")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlan9Error(t *testing.T) {
	ckHasPlan9(t)
	conn, err := Dial("tcp", realPlan9)
	if err != nil {
		t.Fatal(err)
	}
	iounit, version, err := conn.Ver()
	t.Logf("got iounit=%d version=%q\n", iounit, version)
	if err != nil {
		t.Fatalf("error: %s\n", err)
	}
	err = conn.Error("because")
	err = conn.Error("because")
	err = conn.Error("because")
	err = conn.Error("because")
	t.Fatalf("error %s", err)

}
