package stream

import (
	"bodychains/message"

	"github.com/libp2p/go-libp2p-core/peer"
)

type Noop struct{}

type PeerStream struct {
	H        *message.Header
	Wrap     *StreamWrap
	ReadDone chan Noop
}

// Pipe stores channels for communication with StreamsManager
type Pipe struct {
	OnConnect    chan *StreamWrap
	OnDisconnect chan peer.ID
	OnMessage    chan *PeerStream
}

func newPipe() Pipe {
	return Pipe{
		OnConnect:    make(chan *StreamWrap),
		OnDisconnect: make(chan peer.ID),
		OnMessage:    make(chan *PeerStream),
	}
}
