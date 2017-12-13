// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package retry_test

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/clock"
	gc "gopkg.in/check.v1"

	"github.com/juju/retry"
)

type retrySuite struct {
	testing.LoggingSuite
}

var _ = gc.Suite(&retrySuite{})

type mockClock struct {
	now    time.Time
	delays []time.Duration
}

func (mock *mockClock) Now() time.Time {
	return mock.now
}

func (mock *mockClock) After(wait time.Duration) <-chan time.Time {
	mock.delays = append(mock.delays, wait)
	mock.now = mock.now.Add(wait)
	return time.After(time.Microsecond)
}

func (*retrySuite) TestSuccessHasNoDelay(c *gc.C) {
	clock := &mockClock{}
	err := retry.Call(retry.CallArgs{
		Func:     func() error { return nil },
		Attempts: 5,
		Delay:    time.Minute,
		Clock:    clock,
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(clock.delays, gc.HasLen, 0)
}

func (*retrySuite) TestCalledOnceEvenIfStopped(c *gc.C) {
	stop := make(chan struct{})
	clock := &mockClock{}
	called := false
	close(stop)
	err := retry.Call(retry.CallArgs{
		Func: func() error {
			called = true
			return nil
		},
		Attempts: 5,
		Delay:    time.Minute,
		Clock:    clock,
		Stop:     stop,
	})
	c.Assert(called, jc.IsTrue)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(clock.delays, gc.HasLen, 0)
}

func (*retrySuite) TestAttempts(c *gc.C) {
	clock := &mockClock{}
	funcErr := errors.New("bah")
	err := retry.Call(retry.CallArgs{
		Func:     func() error { return funcErr },
		Attempts: 4,
		Delay:    time.Minute,
		Clock:    clock,
	})
	c.Assert(err, jc.Satisfies, retry.IsAttemptsExceeded)
	// We delay between attempts, and don't delay after the last one.
	c.Assert(clock.delays, jc.DeepEquals, []time.Duration{
		time.Minute,
		time.Minute,
		time.Minute,
	})
}

func (*retrySuite) TestAttemptsExceededError(c *gc.C) {
	clock := &mockClock{}
	funcErr := errors.New("bah")
	err := retry.Call(retry.CallArgs{
		Func:     func() error { return funcErr },
		Attempts: 5,
		Delay:    time.Minute,
		Clock:    clock,
	})
	c.Assert(err, gc.ErrorMatches, `attempt count exceeded: bah`)
	c.Assert(err, jc.Satisfies, retry.IsAttemptsExceeded)
	c.Assert(retry.LastError(err), gc.Equals, funcErr)
}

func (*retrySuite) TestFatalErrorsNotRetried(c *gc.C) {
	clock := &mockClock{}
	funcErr := errors.New("bah")
	err := retry.Call(retry.CallArgs{
		Func:         func() error { return funcErr },
		IsFatalError: func(error) bool { return true },
		Attempts:     5,
		Delay:        time.Minute,
		Clock:        clock,
	})
	c.Assert(errors.Cause(err), gc.Equals, funcErr)
	c.Assert(clock.delays, gc.HasLen, 0)
}

func (*retrySuite) TestBackoffFactor(c *gc.C) {
	clock := &mockClock{}
	err := retry.Call(retry.CallArgs{
		Func:        func() error { return errors.New("bah") },
		Clock:       clock,
		Attempts:    5,
		Delay:       time.Minute,
		BackoffFunc: retry.DoubleDelay,
	})
	c.Assert(err, jc.Satisfies, retry.IsAttemptsExceeded)
	c.Assert(clock.delays, jc.DeepEquals, []time.Duration{
		time.Minute,
		time.Minute * 2,
		time.Minute * 4,
		time.Minute * 8,
	})
}

func (*retrySuite) TestStopChannel(c *gc.C) {
	clock := &mockClock{}
	stop := make(chan struct{})
	count := 0
	err := retry.Call(retry.CallArgs{
		Func: func() error {
			if count == 2 {
				close(stop)
			}
			count++
			return errors.New("bah")
		},
		Attempts: 5,
		Delay:    time.Minute,
		Clock:    clock,
		Stop:     stop,
	})
	c.Assert(err, jc.Satisfies, retry.IsRetryStopped)
	c.Assert(clock.delays, gc.HasLen, 3)
}

func (*retrySuite) TestNotifyFunc(c *gc.C) {
	var (
		clock      = &mockClock{}
		funcErr    = errors.New("bah")
		attempts   []int
		funcErrors []error
	)
	err := retry.Call(retry.CallArgs{
		Func: func() error {
			return funcErr
		},
		NotifyFunc: func(lastError error, attempt int) {
			funcErrors = append(funcErrors, lastError)
			attempts = append(attempts, attempt)
		},
		Attempts: 3,
		Delay:    time.Minute,
		Clock:    clock,
	})
	c.Assert(err, jc.Satisfies, retry.IsAttemptsExceeded)
	c.Assert(clock.delays, gc.HasLen, 2)
	c.Assert(funcErrors, jc.DeepEquals, []error{funcErr, funcErr, funcErr})
	c.Assert(attempts, jc.DeepEquals, []int{1, 2, 3})
}

func (*retrySuite) TestInfiniteRetries(c *gc.C) {
	// OK, we can't test infinite, but we'll go for lots.
	clock := &mockClock{}
	stop := make(chan struct{})
	count := 0
	err := retry.Call(retry.CallArgs{
		Func: func() error {
			if count == 111 {
				close(stop)
			}
			count++
			return errors.New("bah")
		},
		Attempts: retry.UnlimitedAttempts,
		Delay:    time.Minute,
		Clock:    clock,
		Stop:     stop,
	})
	c.Assert(err, jc.Satisfies, retry.IsRetryStopped)
	c.Assert(clock.delays, gc.HasLen, count)
}

func (*retrySuite) TestMaxDuration(c *gc.C) {
	clock := &mockClock{}
	err := retry.Call(retry.CallArgs{
		Func:        func() error { return errors.New("bah") },
		Delay:       time.Minute,
		MaxDuration: 5 * time.Minute,
		Clock:       clock,
	})
	c.Assert(err, jc.Satisfies, retry.IsDurationExceeded)
	c.Assert(clock.delays, jc.DeepEquals, []time.Duration{
		time.Minute,
		time.Minute,
		time.Minute,
		time.Minute,
		time.Minute,
	})
}

func (*retrySuite) TestMaxDurationDoubling(c *gc.C) {
	clock := &mockClock{}
	err := retry.Call(retry.CallArgs{
		Func:        func() error { return errors.New("bah") },
		Delay:       time.Minute,
		MaxDuration: 10 * time.Minute,
		BackoffFunc: retry.DoubleDelay,
		Clock:       clock,
	})
	c.Assert(err, jc.Satisfies, retry.IsDurationExceeded)
	// Stops after seven minutes, because the next wait time
	// would take it to 15 minutes.
	c.Assert(clock.delays, jc.DeepEquals, []time.Duration{
		time.Minute,
		2 * time.Minute,
		4 * time.Minute,
	})
}

func (*retrySuite) TestMaxDelay(c *gc.C) {
	clock := &mockClock{}
	err := retry.Call(retry.CallArgs{
		Func:        func() error { return errors.New("bah") },
		Attempts:    7,
		Delay:       time.Minute,
		MaxDelay:    10 * time.Minute,
		BackoffFunc: retry.DoubleDelay,
		Clock:       clock,
	})
	c.Assert(err, jc.Satisfies, retry.IsAttemptsExceeded)
	c.Assert(clock.delays, jc.DeepEquals, []time.Duration{
		time.Minute,
		2 * time.Minute,
		4 * time.Minute,
		8 * time.Minute,
		10 * time.Minute,
		10 * time.Minute,
	})
}

func (*retrySuite) TestWithWallClock(c *gc.C) {
	var attempts []int
	err := retry.Call(retry.CallArgs{
		Func: func() error { return errors.New("bah") },
		NotifyFunc: func(lastError error, attempt int) {
			attempts = append(attempts, attempt)
		},
		Attempts: 5,
		Delay:    time.Microsecond,
		Clock:    clock.WallClock,
	})
	c.Assert(err, jc.Satisfies, retry.IsAttemptsExceeded)
	c.Assert(attempts, jc.DeepEquals, []int{1, 2, 3, 4, 5})
}

func (*retrySuite) TestMissingFuncNotValid(c *gc.C) {
	err := retry.Call(retry.CallArgs{
		Attempts: 5,
		Delay:    time.Minute,
		Clock:    clock.WallClock,
	})
	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `missing Func not valid`)
}

func (*retrySuite) TestMissingAttemptsNotValid(c *gc.C) {
	err := retry.Call(retry.CallArgs{
		Func:  func() error { return errors.New("bah") },
		Delay: time.Minute,
		Clock: clock.WallClock,
	})
	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `missing Attempts or MaxDuration not valid`)
}

func (*retrySuite) TestMissingDelayNotValid(c *gc.C) {
	err := retry.Call(retry.CallArgs{
		Func:     func() error { return errors.New("bah") },
		Attempts: 5,
		Clock:    clock.WallClock,
	})
	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `missing Delay not valid`)
}

func (*retrySuite) TestMissingClockNotValid(c *gc.C) {
	err := retry.Call(retry.CallArgs{
		Func:     func() error { return errors.New("bah") },
		Attempts: 5,
		Delay:    time.Minute,
	})
	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `missing Clock not valid`)
}
