package stream

import (
	"encoding/gob"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"golang.org/x/net/context"

	"bodychains/connection"

	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mdnsDiscovery "github.com/libp2p/go-libp2p/p2p/discovery"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}

type StreamWrap struct {
	s network.Stream

	Read  sync.Mutex
	Write sync.Mutex

	Enc *gob.Encoder
	Dec *gob.Decoder
}

type StreamsMap map[peer.ID]*StreamWrap
type ignoreMap map[peer.ID]struct{}

// StreamsManager manages connections streams
type StreamsManager struct {
	sync.Mutex

	ignore     ignoreMap
	list       StreamsMap
	ctx        context.Context
	host       host.Host
	protoID    protocol.ID
	dht        *dht.IpfsDHT
	rendezvous string

	Pipe Pipe
}

// NewStreamsManager creates StreamsManager
func NewStreamsManager(ctx context.Context, host host.Host, proto protocol.ID, rendezvous string) *StreamsManager {
	sm := StreamsManager{
		ignore:     ignoreMap{},
		list:       StreamsMap{},
		ctx:        ctx,
		host:       host,
		protoID:    proto,
		dht:        nil,
		rendezvous: rendezvous,
		Pipe:       newPipe(),
	}

	host.SetStreamHandler(proto, sm.MakeHandleStream())
	// handle streams opened by the remote side.
	// host.Network().SetStreamHandler(sm.MakeHandleStreamRemote())

	// host.Network().Notify()

	mdns, err := mdnsDiscovery.NewMdnsService(ctx, host, time.Second*10, rendezvous)
	if err != nil {
		panic(err)
	}
	mdns.RegisterNotifee(&mdnsNotifee{h: host, ctx: ctx})

	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	sm.dht = kademliaDHT

	return &sm
}

func (sm *StreamsManager) GetStreams() StreamsMap {
	return sm.list
}

// StartDiscover creates streams for a new found peers
func (sm *StreamsManager) StartDiscover() {
	routingDiscovery := discovery.NewRoutingDiscovery(sm.dht)
	discovery.Advertise(sm.ctx, routingDiscovery, sm.rendezvous)

	go func() {
		for {
			peerChan, err := routingDiscovery.FindPeers(sm.ctx, sm.rendezvous)
			if err != nil {
				panic(err)
			}

			for peer := range peerChan {
				if peer.ID == sm.host.ID() || peer.ID == "" {
					continue
				}

				sm.Lock()
				_, found := sm.list[peer.ID]

				if !found {
					sm.makeStream(peer.ID)
				}
				sm.Unlock()
			}

			wait := 10 * time.Second
			select {
			case <-time.After(wait):
			case <-sm.ctx.Done():
				return
			}
		}
	}()
}

func (sm *StreamsManager) MakeHandleStream() network.StreamHandler {
	return func(stream network.Stream) {
		// fmt.Println(stream.Protocol())

		sm.Lock()
		defer sm.Unlock()

		id := stream.Conn().RemotePeer()

		_, found := sm.list[id]

		if found {
			return
		}
		sm.addStream(stream)
	}
}

func (sm *StreamsManager) ignorePeer(id peer.ID) {
	sm.ignore[id] = struct{}{}
}

func (sm *StreamsManager) isIgnored(id peer.ID) bool {
	_, ok := sm.ignore[id]
	return ok
}

func (sm *StreamsManager) closeByPeer(id peer.ID) {
	if stream, ok := sm.list[id]; ok {
		if err := stream.s.Close(); err != nil {
			log.Println(err)
		}

		delete(sm.list, id)

		// peerLong := stream.Conn().RemotePeer()
		sm.Pipe.OnDisconnect <- id
	}
}

func (sm *StreamsManager) makeStream(peerID peer.ID) {
	if _, found := sm.list[peerID]; found {
		return
	}

	stream, err := sm.host.NewStream(sm.ctx, peerID, sm.protoID)
	if err != nil {
		// failed to dial
		// sm.ignorePeer(peerID)
		return
	}
	sm.addStream(stream)
}

func (sm *StreamsManager) addStream(stream network.Stream) {
	id := stream.Conn().RemotePeer()

	if _, found := sm.list[id]; found {
		return
	}

	sm.Pipe.OnConnect <- connection.New(stream)
}
