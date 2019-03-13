package network

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"
	etcdcv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	e "google.golang.org/genproto"
)

const (
	priority    = 10  // default priority when nothing is set
	ttl         = 300 // default ttl when nothing is set
	etcdTimeout = 5 * time.Second
)

type Network struct {
	Next       plugin.Handler
	Zones      []string
	Upstream   *upstream.Upstream
	Endpoints  []string
	PathPrefix string
	Client     *etcdcv3.Client
	Ctx        context.Context
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

// Services implements the ServiceBackend interface.
func (n *Network) Services(state request.Request, exact bool, opt plugin.Options) (services []msg.Service, err error) {
	services, err = n.Records(state, exact)
	if err != nil {
		return
	}

	services = msg.Group(services)
	return
}

func (n *Network) Records(state request.Request, exact bool) ([]msg.Service, error) {
	name := state.Name()

	path, star := msg.PathWithWildcard(name, n.PathPrefix)
	r, err := e.get(path)
	if err != nil {
		return nil, err
	}
	segments := strings.Split(msg.Path(name, n.PathPrefix), "/")
	return n.loopNodes(r.Kvs, segments, star, state.QType())
}

func (n *Network) get(path string) (*etcdcv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(n.Ctx, etcdTimeout)
	defer cancel()
	r, err := n.Client.Get(ctx, path, etcdcv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (n *Network) loopNodes(kv []*mvccpb.KeyValue, nameParts []string, star bool, qType uint16) (sx []msg.Service, err error) {
	bx := make(map[msg.Service]struct{})
Nodes:
	for _, n := range kv {
		
		serv := new(msg.Service)
		if err := json.Unmarshal(n.Value, serv); err != nil {
			return nil, fmt.Errorf("%s: %s", n.Key, err.Error())
		}
		serv.Key = string(n.Key)
		if _, ok := bx[*serv]; ok {
			continue
		}
		bx[*serv] = struct{}{}

		serv.TTL = e.TTL(n, serv)
		if serv.Priority == 0 {
			serv.Priority = priority
		}

		if shouldInclude(serv, qType) {
			sx = append(sx, *serv)
		}
	}
	return sx, nil
}
