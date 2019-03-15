package network

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/coredns/coredns/plugin"
	etcdcv3 "github.com/coreos/etcd/clientv3"
)

const (
	priority    = 10  // default priority when nothing is set
	ttl         = 300 // default ttl when nothing is set
	etcdTimeout = 5 * time.Second
)

var errKeyNotFound = errors.New("key not found")
var errParse = errors.New("parse etcd fail")

type Network struct {
	Next       plugin.Handler
	Zones      []string
	Endpoints  []string
	PathPrefix string
	Ctx        context.Context
	Client     *etcdcv3.Client
}

func New(zones []string) *Network {
	n := new(Network)
	n.Zones = zones
	return n
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
