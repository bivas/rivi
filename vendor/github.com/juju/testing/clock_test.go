// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing_test

import (
	"sync"
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
)

type clockSuite struct {
	testing.LoggingSuite
}

var _ = gc.Suite(&clockSuite{})

func (*clockSuite) TestNow(c *gc.C) {
	t0 := time.Now()
	cl := testing.NewClock(t0)
	c.Assert(cl.Now(), gc.Equals, t0)
}

var (
	shortWait = 50 * time.Millisecond
	longWait  = time.Second
)

func (*clockSuite) TestAdvanceLogs(c *gc.C) {
	t0 := time.Now()
	cl := testing.NewClock(t0)

	// Shouldn't log anything.
	t := cl.After(time.Second)
	cl.Advance(time.Minute)
	<-t
	c.Check(c.GetTestLog(), jc.DeepEquals, "")

	// Should log since nothing's waiting.
	cl.Advance(time.Hour)
	c.Check(c.GetTestLog(), jc.Contains, "advancing a clock that has nothing waiting: cf. https://github.com/juju/juju/wiki/Intermittent-failures")
}

func (*clockSuite) TestWaitAdvance(c *gc.C) {
	t0 := time.Now()
	cl := testing.NewClock(t0)

	// Test that no timers errors out.
	err := cl.WaitAdvance(time.Millisecond, 10*time.Millisecond, 1)
	c.Check(err, gc.ErrorMatches, "got 0 timers added after waiting 10ms: wanted 1")

	// Test that a timer doesn't error.
	_ = cl.After(time.Nanosecond)
	err = cl.WaitAdvance(time.Millisecond, 10*time.Millisecond, 1)
	c.Check(err, jc.ErrorIsNil)
}

func (*clockSuite) TestAdvanceWithAfter(c *gc.C) {
	t0 := time.Now()
	cl := testing.NewClock(t0)
	ch := cl.After(time.Second)
	select {
	case <-ch:
		c.Fatalf("received unexpected event")
	case <-time.After(shortWait):
	}

	cl.Advance(time.Second - 1)

	select {
	case <-ch:
		c.Fatalf("received unexpected event")
	case <-time.After(shortWait):
	}

	cl.Advance(1)

	select {
	case <-ch:
	case <-time.After(longWait):
		c.Fatalf("expected event to be triggered")
	}

	cl.Advance(time.Second)
	select {
	case <-ch:
		c.Fatalf("received unexpected event")
	case <-time.After(shortWait):
	}

	// Test that we can do it again
	ch = cl.After(time.Second)
	cl.Advance(2 * time.Second)
	select {
	case <-ch:
	case <-time.After(longWait):
		c.Fatalf("expected event to be triggered")
	}
	c.Assert(cl.Now().UTC(), gc.Equals, t0.Add(4*time.Second).UTC())
}

func (*clockSuite) TestAdvanceWithAfterFunc(c *gc.C) {
	// Most of the details have been checked in TestAdvanceWithAfter,
	// so just check that AfterFunc is wired up correctly.
	t0 := time.Now()
	cl := testing.NewClock(t0)
	fired := make(chan struct{})
	cl.AfterFunc(time.Second, func() {
		close(fired)
	})
	cl.Advance(2 * time.Second)
	select {
	case <-fired:
	case <-time.After(longWait):
		c.Fatalf("expected event to be triggered")
	}
}

func (*clockSuite) TestAfterFuncStop(c *gc.C) {
	t0 := time.Now()
	cl := testing.NewClock(t0)
	fired := make(chan struct{})
	timer := cl.AfterFunc(time.Second, func() {
		close(fired)
	})
	cl.Advance(50 * time.Millisecond)
	timer.Stop()
	select {
	case <-fired:
		c.Fatalf("received unexpected event")
	case <-time.After(shortWait):
	}
}

func (*clockSuite) TestNewTimerReset(c *gc.C) {
	t0 := time.Now()
	cl := testing.NewClock(t0)
	timer := cl.NewTimer(time.Second)
	cl.Advance(time.Second)
	select {
	case t := <-timer.Chan():
		c.Assert(t.UTC(), gc.Equals, t0.Add(time.Second).UTC())
	case <-time.After(longWait):
		c.Fatalf("expected event to be triggered")
	}

	timer.Reset(50 * time.Millisecond)
	cl.Advance(100 * time.Millisecond)
	select {
	case t := <-timer.Chan():
		c.Assert(t.UTC(), gc.Equals, t0.Add(time.Second+100*time.Millisecond).UTC())
	case <-time.After(longWait):
		c.Fatalf("expected event to be triggered")
	}
}

func (*clockSuite) TestMultipleWaiters(c *gc.C) {
	var wg sync.WaitGroup
	t0 := time.Date(2000, 01, 01, 01, 0, 0, 0, time.UTC)
	cl := testing.NewClock(t0)

	total := 0
	start := func(f func()) {
		total++
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	start(func() {
		<-cl.After(50 * time.Millisecond)
	})
	start(func() {
		ch := make(chan struct{})
		cl.AfterFunc(100*time.Millisecond, func() {
			close(ch)
		})
		<-ch
	})
	start(func() {
		timer := cl.NewTimer(150 * time.Millisecond)
		<-timer.Chan()
		timer.Reset(50 * time.Millisecond)
		<-timer.Chan()
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all the alarms to be waited on.
	for i := 0; i < total; i++ {
		select {
		case <-cl.Alarms():
		case <-time.After(longWait):
			c.Fatalf("expected a notification on the alarms channel")
		}
	}
	select {
	case <-cl.Alarms():
		c.Fatalf("unexpected extra notification on alarms channel")
	case <-time.After(shortWait):
	}

	cl.Advance(150 * time.Millisecond)

	// Wait for the extra notification after reset.
	select {
	case <-cl.Alarms():
	case <-time.After(longWait):
		c.Fatalf("expected a notification on the alarms channel")
	}

	cl.Advance(50 * time.Millisecond)

	select {
	case <-done:
	case <-time.After(longWait):
		c.Fatalf("expected all waits to complete")
	}

}
