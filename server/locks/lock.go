package locks

import (
    "github.com/furio/icycle/server/locks/peer"
)

type Lock interface {
    InitClient(ipAddresses []string, directory string)
    Register(worker uint64, datacenter uint64) bool
    Peers(datacenter int64) map[int64]peer.Peer
    Synced(peerMap map[int64]peer.Peer) bool
}

const (
    lockDirectory = "icycle"
)

func NewClient(clientType string, ipAddresses []string, directory string) *Lock {
    return nil
}