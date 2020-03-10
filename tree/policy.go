package tree

import "time"

// BackoffPolicy is a factory for backoff policies of a given type.
type BackoffPolicy interface {
	NewBackoffState() BackoffState
}

// BackoffState implements the runtime state of a specific backoff policy.
//
// Implementations of BackoffState are purely concerned with the "arithmetic"
// of computing when the respective timer should be cleared (e.g. for making new connection retries).
//
// This interface allows for the implementation of flexible backoff policies.
// For instance, a policy could treat a burst of backoffs as a single one.
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
