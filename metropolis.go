package main

import (
	"github.com/sean9999/metropolis/event"
	"github.com/sean9999/metropolis/graph"
)

// a Metropolis is a SocialNetwork with an event bus
type Metropolis struct {
	Network *graph.SocialNetwork
	Bus     *event.Bus
}

func NewMetropolis() (Metropolis, error) {
	b := event.NewBus()
	n, err := graph.NewSocialNetwork()
	return Metropolis{n, b}, err
}
