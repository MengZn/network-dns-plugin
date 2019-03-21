package network

import (
	"context"
	"encoding/json"
	"fmt"

	etcdcv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type reslove struct {
	Ip []string `json:"ip"`
}

func (n *Network) get(path string) (*etcdcv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(n.Ctx, etcdTimeout)
	defer cancel()
	fmt.Printf("%s\n", path)
	fmt.Printf("====================\n")
	r, err := n.Client.Get(ctx, path)
	if err != nil {
		fmt.Printf("Client.Get Error \n")
		fmt.Errorf("%s", err)
		fmt.Printf("====================\n")
		return nil, err
	}
	if r.Count == 0 {
		fmt.Printf("count = 0")
		fmt.Printf("====================\n")
		return nil, errKeyNotFound
	}
	return r, nil
}

func (n *Network) parseReslove(kv []*mvccpb.KeyValue) (*reslove, error) {
	fmt.Printf("this value is %s\n", kv[0].Value)
	fmt.Printf("====================\n")
	reslove := new(reslove)
	if err := json.Unmarshal(kv[0].Value, reslove); err != nil {
		return nil, fmt.Errorf("%s: %s", kv[0].Key, err.Error())
	}
	fmt.Printf("this struct is %s\n", reslove)
	fmt.Printf("====================\n")
	return reslove, nil
}
