package etcd

import (
    "github.com/coreos/go-etcd/etcd"
    "github.com/furio/icycle/server/locks/peer"
)

type EtcdLock struct {
    client *etcd.Client
    rootDirectory string
}

// Docs: https://godoc.org/github.com/coreos/go-etcd/etcd

func (c *EtcdLock) InitClient(ipAddresses []string, directory string) {
    c.rootDirectory = directory
    c.client = etcd.NewClient(ipAddresses)

    // func (c *Client) CreateDir(key string, ttl uint64) (*Response, error)
}

func (c *EtcdLock) Register(worker uint64, datacenter uint64) bool {

    /*
        if _, err := client.Set("/foo", "bar", 0); err != nil {
            log.Fatal(err)
        }
    */

    // func (c *Client) Set(key string, value string, ttl uint64) (*Response, error)

    return nil
}

func (c *EtcdLock) Peers(datacenter int64) map[int64]peer.Peer {

    // func (c *Client) Get(key string, sort, recursive bool) (*Response, error)

    return nil
}

func (c *EtcdLock) Synced(peerMap map[int64]peer.Peer) bool {

    return nil
}