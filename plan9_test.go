package p9p

import (
	"context"
	"io"
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

func TestPlan9Version(t *testing.T) {
	testPlan9Version(t)
}

// The server should run the equivalent of
//	aux/listen1 -t -v tcp!*!808 exportfs -r /
func testPlan9Version(t *testing.T) *Conn {
	t.Helper()
	ckHasPlan9(t)
	conn, err := Dial("tcp", realPlan9)
	if err != nil {
		t.Fatal(err)
	}
	iounit, version, err := conn.Ver()
	t.Logf("iounit, version: %v, %v\n", iounit, version)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func TestPlan9WalkWrite(t *testing.T) {
	conn := testPlan9Version(t)
	qid, err := conn.Attach(1, NoFid, "none", "/tmp")
	if err != nil {
		t.Fatal(err)
	}
	qid, iounit, err := conn.Create(1, "testglenda4.txt", 0, 0)
	t.Logf("qid, iounit, err %v %v %v\n", qid, iounit, err)
	if err != nil {
		t.Fatal(err)
	}
	p := make([]byte, iounit)
	n, err := conn.WriteFid(2, 0, p)
	t.Logf("n, err, p %v %v %v\n", n, err, p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlan9WalkReadNdb(t *testing.T) {
	conn := testPlan9Version(t)
	qid, err := conn.Attach(1, NoFid, "none", "/")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("qid is %#v\n", qid)
	q, err := conn.Walk(1, 2, "/", "lib", "ndb", "local")
	t.Logf("quids, err  %v %v\n", q, err)
	if err != nil {
		t.Fatal(err)
	}
	qid, iounit, err := conn.Open(2, 0)
	t.Logf("qid, iounit, err %v %v %v\n", qid, iounit, err)
	if err != nil {
		t.Fatal(err)
	}
	p := make([]byte, iounit)
	n, err := conn.ReadFid(2, 0, p)
	t.Logf("n, err, p %v %v %v\n", n, err, p)
	if err != io.EOF && err != nil {
		t.Fatal(err)
	}
}

func TestPlan9Create(t *testing.T) {
	conn := testPlan9Version(t)
	qid, err := conn.Attach(1, NoFid, "none", "/")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("qid is %#v\n", qid)
	qid, iounit, err := conn.Create(1, "glend4.txt", 0, 0)
	t.Logf("qid, iounit, err %v %v %v\n", qid, iounit, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlan9Flush(t *testing.T) {
	conn := testPlan9Version(t)
	err := conn.Flush(0xffff)
	if err != nil {
		t.Fatalf("error: %s\n", err)
	}
}

func TestPlan9Error(t *testing.T) {
	t.Skip("client doesnt send errors to server")
	conn := testPlan9Version(t)
	err := conn.Error("because")
	if err != nil {
		t.Fatalf("error %s", err)
	}
}
