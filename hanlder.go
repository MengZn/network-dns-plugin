package network

import (
	"context"
	"math/rand"
	"net"

	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// ServeDNS implements the Handler interface.
func (n *Network) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	path, _ := msg.PathWithWildcard(state.Name(), n.PathPrefix)
	res, err := n.get(path)
	if err != nil {
		return 0, errKeyNotFound
	}
	resp, err := n.parseReslove(res.Kvs)
	if err != nil {
		return 0, errParse
	}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	var rr dns.RR

	rr = new(dns.A)
	rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
	rr.(*dns.A).A = net.ParseIP(resp.Ip[rand.Intn(len(resp.Ip))]).To4()

	a.Answer = []dns.RR{rr}
	w.WriteMsg(a)

	return 0, nil
}

// Name implements the Handler interface.
func (n *Network) Name() string { return "network" }
