package p2p

import (
	log "github.com/UBChainNet/UBChain/log/log15"
	"github.com/UBChainNet/UBChain/param"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"strings"
)

// Default boot node list
var DefaultBootstrapPeers []multiaddr.Multiaddr
// Custom boot node list
var CustomBootstrapPeers []multiaddr.Multiaddr

func init() {
	for _, s := range param.Boots {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
}

func IsInBootstrapPeers(id peer.ID) bool {
	bootstrap := DefaultBootstrapPeers
	if len(CustomBootstrapPeers) > 0 {
		bootstrap = CustomBootstrapPeers
	}
	for _, bootstrap := range bootstrap {
		if id.String() == strings.Split(bootstrap.String(), "/")[6] {
			return true
		}
	}
	return false
}

// Start boot node
func (p *P2pServer) StartBootStrap() error {
	var err error
	p.dht, err = dht.New(p.ctx, p.host)
	if err != nil {
		return err
	}
	log.Info("Bootstrapping the DHT")
	if err = p.dht.Bootstrap(p.ctx); err != nil {
		return err
	}
	return nil
}
