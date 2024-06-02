package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/sean9999/go-oracle"
)

// Hash is the condensed uniqueness of an Actor
type Hash string

// an Actor is a peer in an [SocialNetwork]
type Actor struct {
	oracle.Peer
	Network SocialNetwork
}

func (a Actor) Hash() Hash {
	return Hash(a.Peer.Nickname())
}

func (a Actor) Befriend(b Actor) error {
	e := Edge{
		RelashionshipType: Friend,
		From:              a.Hash(),
		To:                b.Hash(),
	}
	return a.Network.Graph.AddEdge(a.Hash(), b.Hash(), graph.EdgeAttributes(e.Attributes()))
}

func (a Actor) Marry(b Actor) error {
	e := Edge{
		RelashionshipType: Spouse,
		From:              a.Hash(),
		To:                b.Hash(),
	}
	return a.Network.Graph.AddEdge(a.Hash(), b.Hash(), graph.EdgeAttributes(e.Attributes()))
}

func (a Actor) Follow(b Actor) error {
	e := Edge{
		RelashionshipType: Follow,
		From:              a.Hash(),
		To:                b.Hash(),
	}
	return a.Network.Graph.AddEdge(a.Hash(), b.Hash(), graph.EdgeAttributes(e.Attributes()))
}
