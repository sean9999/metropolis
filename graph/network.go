package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/sean9999/go-oracle"
)

// a SocialNetwork is a graph with a central [Actor]
type SocialNetwork struct {
	Me    oracle.Oracle
	Graph graph.Graph[Hash, Actor]
}

func NewSocialNetwork() (*SocialNetwork, error) {
	me, err := oracle.FromFile("mothership.toml")
	if err != nil {
		return nil, err
	}
	oracleHash := func(a Actor) Hash {
		return a.Hash()
	}
	g := graph.New(oracleHash)
	m := SocialNetwork{
		Me:    *me,
		Graph: g,
	}
	return &m, nil
}
