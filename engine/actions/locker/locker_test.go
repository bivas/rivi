package locker

import (
	"testing"

	"github.com/bivas/rivi/mocks"
	"github.com/bivas/rivi/util/log"
	"github.com/bivas/rivi/util/state"
	"github.com/stretchr/testify/assert"
)

type mockLockableData struct {
	mocks.MockData
	locked                   bool
	lockCalled, unlockCalled bool
}

func (m *mockLockableData) Lock() {
	m.locked = true
	m.lockCalled = true
}

func (m *mockLockableData) Unlock() {
	m.locked = false
	m.unlockCalled = true
}

func (m *mockLockableData) LockState() bool {
	return m.locked
}

func TestSerialization(t *testing.T) {
	input := map[string]interface{}{
		"state": "lock",
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, "lock", s.rule.State, "state")
}

func TestNotLockable(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("locker.test")}
	meta := &mocks.MockData{Labels: []string{}}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.NotNil(t, action.err, "can't merge")
}

func TestLock(t *testing.T) {
	action := action{rule: &rule{State: "lock"}, logger: log.Get("locker.test")}
	meta := &mockLockableData{MockData: mocks.MockData{}}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Nil(t, action.err, "shouldn't error")
	assert.True(t, meta.locked, "should be locked")
	assert.True(t, meta.lockCalled, "should be locked")
}

func TestLockWhenLocked(t *testing.T) {
	action := action{rule: &rule{State: "lock"}, logger: log.Get("locker.test")}
	meta := &mockLockableData{MockData: mocks.MockData{}, locked: true}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Nil(t, action.err, "shouldn't error")
	assert.True(t, meta.locked, "should be locked")
	assert.False(t, meta.lockCalled, "no need to relock")
}

func TestUnlockWhenLocked(t *testing.T) {
	action := action{rule: &rule{State: "unlock"}, logger: log.Get("locker.test")}
	meta := &mockLockableData{MockData: mocks.MockData{}, locked: true}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Nil(t, action.err, "shouldn't error")
	assert.False(t, meta.locked, "should be unlocked")
	assert.True(t, meta.unlockCalled, "should be unlocked")
}

func TestUnlockWhenUnlocked(t *testing.T) {
	action := action{rule: &rule{State: "unlock"}, logger: log.Get("locker.test")}
	meta := &mockLockableData{MockData: mocks.MockData{}}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Nil(t, action.err, "shouldn't error")
	assert.False(t, meta.locked, "should be unlocked")
	assert.False(t, meta.unlockCalled, "no need to re-unlock")
}

func TestStateChangeFromUnlocked(t *testing.T) {
	action := action{rule: &rule{State: "change"}, logger: log.Get("locker.test")}
	meta := &mockLockableData{MockData: mocks.MockData{}}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Nil(t, action.err, "shouldn't error")
	assert.True(t, meta.locked, "should be locked")
	assert.True(t, meta.lockCalled, "lock")
}

func TestStateChangeFromLocked(t *testing.T) {
	action := action{rule: &rule{State: "change"}, logger: log.Get("locker.test")}
	meta := &mockLockableData{MockData: mocks.MockData{}, locked: true}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Nil(t, action.err, "shouldn't error")
	assert.False(t, meta.locked, "should be unlocked")
	assert.True(t, meta.unlockCalled, "unlock")
}
