// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package testing_test

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/juju/testing"
	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&tcpProxySuite{})

type tcpProxySuite struct{}

func (*tcpProxySuite) TestTCPProxy(c *gc.C) {
	var wg sync.WaitGroup

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, gc.IsNil)
	defer listener.Close()
	wg.Add(1)
	go tcpEcho(&wg, listener)

	p := testing.NewTCPProxy(c, listener.Addr().String())
	c.Assert(p.Addr(), gc.Not(gc.Equals), listener.Addr().String())

	// Dial the proxy and check that we see the text echoed correctly.
	conn, err := net.Dial("tcp", p.Addr())
	c.Assert(err, gc.IsNil)
	defer conn.Close()

	assertEcho(c, conn)

	// Close the connection and check that we see
	// the connection closed for read.
	conn.(*net.TCPConn).CloseWrite()
	assertEOF(c, conn)

	// Make another connection and close the proxy,
	// which should close down the proxy and cause us
	// to get an error.
	conn, err = net.Dial("tcp", p.Addr())
	c.Assert(err, gc.IsNil)
	defer conn.Close()

	p.Close()
	assertEOF(c, conn)

	// Make sure that we cannot dial the proxy address either.
	conn, err = net.Dial("tcp", p.Addr())
	c.Assert(err, gc.ErrorMatches, ".*connection refused")

	listener.Close()
	// Make sure that all our connections have gone away too.
	wg.Wait()
}

func (*tcpProxySuite) TestCloseConns(c *gc.C) {
	var wg sync.WaitGroup

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, gc.IsNil)
	defer listener.Close()
	wg.Add(1)
	go tcpEcho(&wg, listener)

	p := testing.NewTCPProxy(c, listener.Addr().String())
	c.Assert(p.Addr(), gc.Not(gc.Equals), listener.Addr().String())

	// Make a couple of connections through the proxy
	// and test that they work.
	conn1, err := net.Dial("tcp", p.Addr())
	c.Assert(err, gc.IsNil)
	defer conn1.Close()
	assertEcho(c, conn1)

	conn2, err := net.Dial("tcp", p.Addr())
	c.Assert(err, gc.IsNil)
	defer conn1.Close()
	assertEcho(c, conn1)

	p.CloseConns()

	// Assert that both the connections have been broken.
	assertEOF(c, conn1)
	assertEOF(c, conn2)

	// Check that we can still make a connection.
	conn3, err := net.Dial("tcp", p.Addr())
	c.Assert(err, gc.IsNil)
	defer conn3.Close()
	assertEcho(c, conn3)

	// Close the proxy and check that the last connection goes.
	p.Close()
	assertEOF(c, conn3)

	listener.Close()
	// Make sure that all our connections have gone away too.
	wg.Wait()
}

func (*tcpProxySuite) TestPauseConns(c *gc.C) {
	var wg sync.WaitGroup

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, gc.IsNil)
	defer listener.Close()
	wg.Add(1)
	go tcpEcho(&wg, listener)

	p := testing.NewTCPProxy(c, listener.Addr().String())
	c.Assert(p.Addr(), gc.Not(gc.Equals), listener.Addr().String())

	// Make a connection through the proxy
	// and test that it works.
	conn, err := net.Dial("tcp", p.Addr())
	c.Assert(err, gc.IsNil)
	defer conn.Close()
	assertEcho(c, conn)

	p.PauseConns()

	msg := "hello, world\n"
	n, err := fmt.Fprint(conn, msg)
	c.Assert(err, gc.IsNil)
	c.Assert(n, gc.Equals, len(msg))
	assertReadTimeout(c, conn)

	p.ResumeConns()

	buf := make([]byte, n)
	n, err = conn.Read(buf)
	c.Assert(err, gc.IsNil)
	c.Assert(n, gc.Equals, len(msg))
	c.Assert(string(buf), gc.Equals, msg)
}

// tcpEcho listens on the given listener for TCP connections,
// writes all traffic received back to the sender, and calls
// wg.Done when all its goroutines have completed.
func tcpEcho(wg *sync.WaitGroup, listener net.Listener) {
	defer wg.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer conn.Close()
			// Echo anything that was written.
			io.Copy(conn, conn)
		}()
	}
}

func assertEcho(c *gc.C, conn net.Conn) {
	txt := "hello, world\n"
	fmt.Fprint(conn, txt)

	buf := make([]byte, len(txt))
	n, err := io.ReadFull(conn, buf)
	c.Assert(err, gc.IsNil)
	c.Assert(string(buf[0:n]), gc.Equals, txt)
}

func assertEOF(c *gc.C, r io.Reader) {
	n, err := r.Read(make([]byte, 1))
	c.Assert(err, gc.Equals, io.EOF)
	c.Assert(n, gc.Equals, 0)
}

func assertReadTimeout(c *gc.C, conn net.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	c.Assert(err, gc.IsNil)
	defer conn.SetReadDeadline(time.Time{})
	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	c.Assert(n, gc.Equals, 0)
	nerr, ok := err.(net.Error)
	c.Assert(ok, gc.Equals, true)
	c.Assert(nerr.Timeout(), gc.Equals, true)
}
