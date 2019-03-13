package network

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/upstream"
)

type Network struct {
	Next       plugin.Handler
	Zones      []string
	Upstream   *upstream.Upstream
	Endpoints  []string
	PathPrefix string
}

func New(zones []string) *Network {
	n := new(Network)
	n.Zones = zones
	return n
}
