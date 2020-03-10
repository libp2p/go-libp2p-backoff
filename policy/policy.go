package policy

import "time"

// BackoffPolicy is a factory for backoff strategies of a given type.
// A strategy is the state and arithmetic logic of a timer.
type BackoffPolicy interface {
	NewBackoffState() BackoffState
}

// BackoffState implements the runtime state for a specific backoff policy.
//
// Implementations of BackoffState are purely concerned with the "arithmetic"
// of computing when the respective timer should be cleared.
//
// This interface allows for the implementation of flexible backoff policies.
// For instance, a policy could treat a burst of backoffs as a single one.
//
// BackoffState is an analog of github.com/libp2p/go-libp2p-discovery.BackoffStrategy.
// The latter, however, is not able to describe logic that adapts to bursts of backoffs.
type BackoffState interface {

	// Clear informs the policy of the current time and sets its state to cleared.
	Clear(now time.Time)

	// Backoff informs the policy of the current time and sets its state to backing off.
	Backoff(now time.Time)

	// TimeToClear informs the policy of the current time and returns the duration
	// remaining until the back off state is cleared. Zero or negative durations indicate
	// that the state is already cleared.
	TimeToClear(now time.Time) time.Duration
}

type NoBackoffPolicy struct{}

func (NoBackoffPolicy) NewBackoffState() BackoffState {
	return noBackoffState{}
}

type noBackoffState struct{}

func (noBackoffState) Clear(now time.Time) {
}

func (noBackoffState) Backoff(now time.Time) {
}

func (noBackoffState) TimeToClear(now time.Time) time.Duration {
	return time.Duration(0)
}
