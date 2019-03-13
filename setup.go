package network

import (
	"fmt"
	"os"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/mholt/caddy"
)

var log = clog.NewWithPlugin("network")

func init() {
	fmt.Printf("network dns setup!!!!")
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
		n.Next = next
		return n
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
			return network, err
		}
	}
	return network, nil
}
func ParseStanza(c *caddy.Controller) (*Network, error) {
	network := New([]string{""})
	return network, nil
}

const defaultEndpoint = "http://localhost:8888"
