package tree

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/backoff"
)

func NewBackoffTreeTimer(parent backoff.BackoffTimer, policy BackoffPolicy) *BackoffTreeTimer {
	return &BackoffTreeTimer{
		parent:   parent,
		state:    policy.NewBackoffState(),
		children: map[string]*BackoffTreeTimer{},
	}
}

type BackoffTreeTimer struct {
	parent backoff.BackoffTimer
	sync.Mutex
	clk      sync.Mutex // lock for children
	children map[string]*BackoffTreeTimer
	slk      sync.Mutex // lock for state
	state    BackoffState
}

func (t *BackoffTreeTimer) Subtimer(childName string, policy BackoffPolicy) *BackoffTreeTimer {
	t.clk.Lock()
	defer t.clk.Unlock()
	if child := t.children[childName]; child != nil {
		return child
	} else {
		d := NewBackoffTreeTimer(t, policy)
		t.children[childName] = d
		return d
	}
}

var BackoffGCInterval = time.Minute

func (t *BackoffTreeTimer) StartGC() {
	go func() {
		for {
			time.Sleep(BackoffGCInterval)
			t.GC()
		}
	}()
}

// GC runs garbage collection on this timer's descendants.
// GC returns the number of subtrees of this node after garbage collection.
func (t *BackoffTreeTimer) GC() int {
	t.clk.Lock()
	defer t.clk.Unlock()
	for name, child := range t.children {
		if child.GC() == 0 && child.TimeToClear() <= 0 {
			delete(t.children, name)
		}
	}
	return len(t.children)
}

func (t *BackoffTreeTimer) NumChildren() int {
	t.clk.Lock()
	defer t.clk.Unlock()
	return len(t.children)
}

// Wait implements BackoffTimer interface.
func (t *BackoffTreeTimer) Wait() {
	time.Sleep(t.TimeToClear())
}

// TimeToClear implements BackoffTimer interface.
func (t *BackoffTreeTimer) TimeToClear() time.Duration {
	pttc := time.Duration(0)
	if p := t.parent; p != nil {
		pttc = p.TimeToClear()
	}
	t.slk.Lock()
	defer t.slk.Unlock()
	return maxDur(pttc, t.state.TimeToClear(time.Now()))
}

func maxDur(dx, dy time.Duration) time.Duration {
	if dx > dy {
		return dx
	}
	return dy
}

// Clear implements BackoffTimer interface.
func (t *BackoffTreeTimer) Clear() {
	t.slk.Lock()
	defer t.slk.Unlock()
	t.state.Clear(time.Now())
}

// Backoff implements BackoffTimer interface.
func (t *BackoffTreeTimer) Backoff() {
	t.slk.Lock()
	defer t.slk.Unlock()
	t.state.Backoff(time.Now())
}
