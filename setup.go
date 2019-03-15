package network

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/mholt/caddy"
)

var log = clog.NewWithPlugin("network")

func init() {
	rand.Seed(time.Now().UnixNano())
	caddy.RegisterPlugin("network", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	os.Stderr = os.Stdout

	if n, err := networkParse(c); err != nil {
		return plugin.Error("network", err)
	} else {
		dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
			n.Next = next
			return n
		})
	}

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
			return nil, err
		}
	}
	return network, nil
}
func ParseStanza(c *caddy.Controller) (*Network, error) {
	fmt.Printf("network dns is in ParseStanza function !!!!\n")

	network := New([]string{""})
	zones := c.RemainingArgs()
	network.Ctx = context.Background()
	if len(zones) != 0 {
		network.Zones = zones
		for i := 0; i < len(network.Zones); i++ {
			network.Zones[i] = plugin.Host(network.Zones[i]).Normalize()
			fmt.Printf("the len is 0 \n")
			fmt.Printf("%v\n", plugin.Host(network.Zones[i]).Normalize())
		}
	} else {
		network.Zones = make([]string, len(c.ServerBlockKeys))
		for i := 0; i < len(c.ServerBlockKeys); i++ {
			network.Zones[i] = plugin.Host(c.ServerBlockKeys[i]).Normalize()
			fmt.Printf("the len is not 0 \n")
			fmt.Printf("%v\n", plugin.Host(c.ServerBlockKeys[i]).Normalize())
		}
	}
	for c.NextBlock() {
		switch c.Val() {
		case "path":
			if !c.NextArg() {
				return nil, c.ArgErr()
			}
			fmt.Printf("%v\n", c.Val())
			network.PathPrefix = c.Val()
		case "endpoint":
			args := c.RemainingArgs()
			if len(args) == 0 {
				return nil, c.ArgErr()
			}
			network.Endpoints = args
		default:
			return nil, c.Errf("unknown property '%s'", c.Val())
		}
	}
	client, err := newEtcdClient(network.Endpoints, nil, "", "")
	// todo
	// client, err := newEtcdClient(endpoints, tlsConfig, username, password)
	if err != nil {
		return nil, c.Errf("etcd client init fail '%v'", err)
	}
	network.Client = client
	return network, nil
}

const defaultEndpoint = "http://localhost:32379"
