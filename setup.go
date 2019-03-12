package network_dns

import (
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("network", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	return nil
}
