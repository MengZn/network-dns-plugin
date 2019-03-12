package network

import (
	"github.com/coredns/coredns/plugin"
)

type Network struct {
	Next  plugin.Handler
	Zones []string
}
