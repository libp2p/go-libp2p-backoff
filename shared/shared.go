package shared

import (
	"sync"

	"github.com/libp2p/go-libp2p-backoff/policy"
	"github.com/libp2p/go-libp2p-backoff/tree"
	"github.com/libp2p/go-libp2p-core/backoff"
	ma "github.com/multiformats/go-multiaddr"
)

// CHECKOUT: https://github.com/libp2p/go-libp2p-discovery

var DefaultIPBackoffPolicy = policy.NoBackoffPolicy{}        // TODO: use a real policy
var DefaultTransportBackoffPolicy = policy.NoBackoffPolicy{} // TODO: use a real policy
var DefaultSwarmBackoffPolicy = policy.NoBackoffPolicy{}     // TODO: use a real policy
var DefaultProtocolBackoffPolicy = policy.NoBackoffPolicy{}  // TODO: use a real policy

func NewSharedBackoffs() backoff.SharedBackoffs {
	b := &sharedBackoffs{
		root: tree.NewBackoffTreeTimer(nil, policy.NoBackoffPolicy{}),
	}
	b.root.StartGC()
	return b
}

type sharedBackoffs struct {
	rlk  sync.Mutex
	root *tree.BackoffTreeTimer
}

func (sh *sharedBackoffs) IP(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, _ := ma.SplitFirst(addr)
	return sh.root.
		Subtimer("dial", policy.NoBackoffPolicy{}). // this creates a namespace for dialing-related timers
		Subtimer(ipComp.String(), DefaultIPBackoffPolicy)
}

func (sh *sharedBackoffs) Transport(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, rest := ma.SplitFirst(addr)
	transportComp, _ := ma.SplitFirst(rest)
	return sh.root.
		Subtimer("dial", policy.NoBackoffPolicy{}). // this creates a namespace for dialing-related timers
		Subtimer(ipComp.String(), DefaultIPBackoffPolicy).
		Subtimer(transportComp.String(), DefaultTransportBackoffPolicy)
}

func (sh *sharedBackoffs) Swarm(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, rest := ma.SplitFirst(addr)
	transportComp, _ := ma.SplitFirst(rest)
	return sh.root.
		Subtimer("dial", policy.NoBackoffPolicy{}). // this creates a namespace for dialing-related timers
		Subtimer(ipComp.String(), DefaultIPBackoffPolicy).
		Subtimer(transportComp.String(), DefaultTransportBackoffPolicy).
		Subtimer("swarm", DefaultSwarmBackoffPolicy)
}

func (sh *sharedBackoffs) Protocol(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, rest := ma.SplitFirst(addr)
	transportComp, rest2 := ma.SplitFirst(rest)
	protocolComp, _ := ma.SplitFirst(rest2)
	return sh.root.
		Subtimer("dial", policy.NoBackoffPolicy{}). // this creates a namespace for dialing-related timers
		Subtimer(ipComp.String(), DefaultIPBackoffPolicy).
		Subtimer(transportComp.String(), DefaultTransportBackoffPolicy).
		Subtimer("swarm", DefaultSwarmBackoffPolicy).
		Subtimer(protocolComp.String(), DefaultProtocolBackoffPolicy)
}
