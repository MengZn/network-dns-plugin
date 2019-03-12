package network_dns

import (
	"github.com/coredns/coredns/plugin"
)

type Network struct {
	Next plugin.Handler
}
