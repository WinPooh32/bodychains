package main

import (
	"bodychains/stream"
	"context"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	ipfslog "github.com/ipfs/go-log"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var logger = ipfslog.Logger("network")

func startNetwork(ctx context.Context) *stream.StreamsManager {
	var (
		rendezvous     = "funnyuniverse0201"
		bootstrapPeers = dht.DefaultBootstrapPeers
		proto          = protocol.ID("/chains/0.0.1")
	)

	// ipfslog.SetLogLevel("network", "info")
	// ipfslog.SetLogLevel("*", "info")

	defaultIP4ListenAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(ctx,
		libp2p.ListenAddrs(defaultIP4ListenAddr),
		libp2p.NATPortMap(),
	)
	if err != nil {
		panic(err)
	}

	// autonat.NewAutoNAT(ctx, host, nil)

	logger.Info("Host created. We are:", host.ID())
	logger.Info(host.Addrs())

	streamsMgr := stream.NewStreamsManager(ctx, host, proto, rendezvous)

	connectToBootstrapPeers(ctx, bootstrapPeers, host)
	streamsMgr.StartDiscover()

	return streamsMgr
}

func connectToBootstrapPeers(ctx context.Context, addlist []multiaddr.Multiaddr, host host.Host) {
	var wg sync.WaitGroup

	for _, peerAddr := range addlist {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)

		go func() {
			if err := host.Connect(ctx, *peerinfo); err != nil {
				logger.Warning(err)
			} else {
				logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
