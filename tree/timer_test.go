package tree

import (
	"testing"
	"time"
)

type TestBackoffPolicy struct{}

func (TestBackoffPolicy) NewBackoffState() BackoffState {
	return &testBackoffState{}
}

type testBackoffState struct {
	backingOff bool
}

func (s *testBackoffState) Clear(now time.Time) {
	s.backingOff = false
}

func (s *testBackoffState) Backoff(now time.Time) {
	s.backingOff = true
}

func (s *testBackoffState) TimeToClear(now time.Time) time.Duration {
	if s.backingOff {
		return time.Second
	} else {
		return time.Duration(0)
	}
}

func TestTimerGC1(t *testing.T) {
	// create root timer
	root := NewBackoffTreeTimer(nil, TestBackoffPolicy{})
	// create a child timer, which is cleared on init
	root.Subtimer("child", TestBackoffPolicy{})
	if root.NumChildren() != 1 {
		t.Errorf("expecting one child")
	}
	// run GC on the root (should remove the child timer)
	if root.GC() != 0 {
		t.Errorf("GC did not remove child")
	}
}

func TestTimerGC2(t *testing.T) {
	// create root timer
	root := NewBackoffTreeTimer(nil, TestBackoffPolicy{})
	// create a child timer, which is cleared on init
	child := root.Subtimer("child", TestBackoffPolicy{})
	// create child of child, in backoff state
	child2 := child.Subtimer("child2", TestBackoffPolicy{})
	child2.Backoff()
	if root.NumChildren() != 1 || child.NumChildren() != 1 {
		t.Errorf("expecting root has one child with one child")
	}
	// run GC on the root (should remove the child timer)
	if root.GC() != 1 {
		t.Errorf("GC removed a child with a descendant in backoff mode")
	}
}
