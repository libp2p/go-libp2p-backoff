package policy

import (
	"math/rand"
	"time"

	discovery "github.com/libp2p/go-libp2p-discovery"
)

type ExpBackoffPolicy struct {
	Min       time.Duration
	Max       time.Duration
	Jitter    discovery.Jitter
	TimeUnits time.Duration
	Base      float64
	Offset    time.Duration
	RNG       *rand.Rand
}

func (p *ExpBackoffPolicy) NewBackoffState() BackoffState {
	return &expBackoffState{
		underlying: discovery.NewExponentialBackoff(
			p.Min,
			p.Max,
			p.Jitter,
			p.TimeUnits,
			p.Base,
			p.Offset,
			p.RNG,
		)(),
	}
}

type expBackoffState struct {
	offUntil   time.Time
	underlying discovery.BackoffStrategy
}

func (s *expBackoffState) Clear(now time.Time) {
	s.offUntil = now
	s.underlying.Reset()
}

func (s *expBackoffState) Backoff(now time.Time) {
	// Unlike github.com/libp2p/go-libp2p-discovery, we do not
	// count backoffs that occur on top of prior backoffs.
	if now.After(s.offUntil) {
		s.offUntil = now.Add(s.underlying.Delay())
	}
}

func (s *expBackoffState) TimeToClear(now time.Time) time.Duration {
	return time.Duration(0)
}
