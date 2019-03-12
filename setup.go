package network

import (
	"os"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/mholt/caddy"
)

var log = clog.NewWithPlugin("network")

func init() {
	caddy.RegisterPlugin("network", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	os.Stderr = os.Stdout

	n, err := networkParse(c)
	if err != nil {
		return plugin.Error("network", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		k.Next = next
		return k
	})

	return nil
}

func networkParse(c *caddy.Controller) (*Network, error) {
	var (
		network *Network
		err     error
	)

	i := 0
	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++

		network, err = ParseStanza(c)
		if err != nil {
			return k8s, err
		}
	}
	return network, nil
}
func ParseStanza(c *caddy.Controller) (*Network, error) {
	network := New([]string{""})
	zones := c.RemainingArgs()

	//todo
	if len(zones) != 0 {
		network.Zones = zones
		for i := 0; i < len(k8s.Zones); i++ {
			network.Zones[i] = plugin.Host(k8s.Zones[i]).Normalize()
		}
	} else {
		network.Zones = make([]string, len(c.ServerBlockKeys))
		for i := 0; i < len(c.ServerBlockKeys); i++ {
			network.Zones[i] = plugin.Host(c.ServerBlockKeys[i]).Normalize()
		}
	}
	for c.NextBlock {

	}
	return network, nil
}

const defaultEndpoint = "http://localhost:8888"
