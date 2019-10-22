package stream

import (
	"bodychains/connection"

	"github.com/libp2p/go-libp2p-core/peer"
)

// Pipe stores channels for communication with StreamsManager
type Pipe struct {
	OnConnect    chan *connection.Connection
	OnDisconnect chan peer.ID
}

func newPipe() Pipe {
	return Pipe{
		OnConnect:    make(chan *connection.Connection),
		OnDisconnect: make(chan peer.ID),
	}
}
