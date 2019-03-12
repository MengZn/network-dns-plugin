package network

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/coredns/coredns/plugin"
)

const (
	priority    = 10  // default priority when nothing is set
	ttl         = 300 // default ttl when nothing is set
	etcdTimeout = 5 * time.Second
)

type Network struct {
	Next      plugin.Handler
	Zones     []string
	Endpoints []string
	Client    *etcdcv3.Client
	Ctx       context.Context
}

func newEtcdClient(endpoints []string, cc *tls.Config, username, password string) (*etcdcv3.Client, error) {
	etcdCfg := etcdcv3.Config{
		Endpoints: endpoints,
		TLS:       cc,
	}
	if username != "" && password != "" {
		etcdCfg.Username = username
		etcdCfg.Password = password
	}
	cli, err := etcdcv3.New(etcdCfg)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
