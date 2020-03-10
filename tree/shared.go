package tree

import (
	"sync"

	"github.com/libp2p/go-libp2p-core/backoff"
	ma "github.com/multiformats/go-multiaddr"
)

// CHECKOUT: https://github.com/libp2p/go-libp2p-discovery

var DefaultIPBackoffPolicy = NoBackoffPolicy{}        // TODO: use a real policy
var DefaultTransportBackoffPolicy = NoBackoffPolicy{} // TODO: use a real policy
var DefaultSwarmBackoffPolicy = NoBackoffPolicy{}     // TODO: use a real policy
var DefaultProtocolBackoffPolicy = NoBackoffPolicy{}  // TODO: use a real policy

func NewSharedBackoffs() backoff.SharedBackoffs {
	b := &sharedBackoffs{
		root: NewBackoffTreeTimer(nil, NoBackoffPolicy{}),
	}
	b.root.StartGC()
	return b
}

type sharedBackoffs struct {
	rlk  sync.Mutex
	root *BackoffTreeTimer
}

func (sh *sharedBackoffs) IP(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, _ := ma.SplitFirst(addr)
	return sh.root.Subtimer(ipComp.String(), DefaultIPBackoffPolicy)
}

func (sh *sharedBackoffs) Transport(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, rest := ma.SplitFirst(addr)
	transportComp, _ := ma.SplitFirst(rest)
	return sh.root.
		Subtimer(ipComp.String(), DefaultIPBackoffPolicy).
		Subtimer(transportComp.String(), DefaultTransportBackoffPolicy)
}

func (sh *sharedBackoffs) Swarm(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, rest := ma.SplitFirst(addr)
	transportComp, _ := ma.SplitFirst(rest)
	return sh.root.
		Subtimer(ipComp.String(), DefaultIPBackoffPolicy).
		Subtimer(transportComp.String(), DefaultTransportBackoffPolicy).
		Subtimer("swarm", DefaultSwarmBackoffPolicy)
}

func (sh *sharedBackoffs) Protocol(addr ma.Multiaddr) backoff.BackoffTimer {
	ipComp, rest := ma.SplitFirst(addr)
	transportComp, rest2 := ma.SplitFirst(rest)
	protocolComp, _ := ma.SplitFirst(rest2)
	return sh.root.
		Subtimer(ipComp.String(), DefaultIPBackoffPolicy).
		Subtimer(transportComp.String(), DefaultTransportBackoffPolicy).
		Subtimer("swarm", DefaultSwarmBackoffPolicy).
		Subtimer(protocolComp.String(), DefaultProtocolBackoffPolicy)
}
