package locker

import (
	"testing"

	"github.com/bivas/rivi/bot/mock"
	"github.com/bivas/rivi/util/log"
	"github.com/stretchr/testify/assert"
)

type mockLockableEventData struct {
	mock.MockEventData
	locked                   bool
	lockCalled, unlockCalled bool
}

func (m *mockLockableEventData) Lock() {
	m.locked = true
	m.lockCalled = true
}

func (m *mockLockableEventData) Unlock() {
	m.locked = false
	m.unlockCalled = true
}

func (m *mockLockableEventData) LockState() bool {
	return m.locked
}

func TestNotLockable(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("locker.test")}
	meta := &mock.MockEventData{Labels: []string{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.NotNil(t, action.err, "can't merge")
}

func TestLock(t *testing.T) {
	action := action{rule: &rule{State: "lock"}, logger: log.Get("locker.test")}
	meta := &mockLockableEventData{MockEventData: mock.MockEventData{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Nil(t, action.err, "shouldn't error")
	assert.True(t, meta.locked, "should be locked")
	assert.True(t, meta.lockCalled, "should be locked")
}

func TestLockWhenLocked(t *testing.T) {
	action := action{rule: &rule{State: "lock"}, logger: log.Get("locker.test")}
	meta := &mockLockableEventData{MockEventData: mock.MockEventData{}, locked: true}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Nil(t, action.err, "shouldn't error")
	assert.True(t, meta.locked, "should be locked")
	assert.False(t, meta.lockCalled, "no need to relock")
}

func TestUnlockWhenLocked(t *testing.T) {
	action := action{rule: &rule{State: "unlock"}, logger: log.Get("locker.test")}
	meta := &mockLockableEventData{MockEventData: mock.MockEventData{}, locked: true}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Nil(t, action.err, "shouldn't error")
	assert.False(t, meta.locked, "should be unlocked")
	assert.True(t, meta.unlockCalled, "should be unlocked")
}

func TestUnlockWhenUnlocked(t *testing.T) {
	action := action{rule: &rule{State: "unlock"}, logger: log.Get("locker.test")}
	meta := &mockLockableEventData{MockEventData: mock.MockEventData{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Nil(t, action.err, "shouldn't error")
	assert.False(t, meta.locked, "should be unlocked")
	assert.False(t, meta.unlockCalled, "no need to re-unlock")
}

func TestStateChangeFromUnlocked(t *testing.T) {
	action := action{rule: &rule{State: "change"}, logger: log.Get("locker.test")}
	meta := &mockLockableEventData{MockEventData: mock.MockEventData{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Nil(t, action.err, "shouldn't error")
	assert.True(t, meta.locked, "should be locked")
	assert.True(t, meta.lockCalled, "lock")
}

func TestStateChangeFromLocked(t *testing.T) {
	action := action{rule: &rule{State: "change"}, logger: log.Get("locker.test")}
	meta := &mockLockableEventData{MockEventData: mock.MockEventData{}, locked: true}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Nil(t, action.err, "shouldn't error")
	assert.False(t, meta.locked, "should be unlocked")
	assert.True(t, meta.unlockCalled, "unlock")
}
