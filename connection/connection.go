package connection

import (
	"bodychains/message"
	"context"
	"fmt"

	"encoding/gob"

	"sync"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type ReqID uint32
type Body []byte

// Connection is bi-directional async message channel
type Connection struct {
	s network.Stream

	requests map[ReqID]chan *message.Pack

	lastID ReqID

	chRead          chan *message.Pack
	chWrite         chan *message.Pack
	chRemoteRequest chan *message.Pack

	dec *gob.Decoder
	enc *gob.Encoder

	m sync.Mutex

	cancel context.CancelFunc
}

func New(s network.Stream) *Connection {
	return &Connection{
		s:   s,
		dec: gob.NewDecoder(s),
		enc: gob.NewEncoder(s),
	}
}

func (c *Connection) ID() peer.ID {
	return c.s.Conn().RemotePeer()
}

func (c *Connection) Open(ctx context.Context) chan *message.Pack {
	c.chRead = make(chan *message.Pack, 8)
	c.chWrite = make(chan *message.Pack, 8)
	c.chRemoteRequest = make(chan *message.Pack, 8)

	ctx, c.cancel = context.WithCancel(ctx)

	go c.read(ctx)
	go c.write(ctx)

	return c.chRemoteRequest
}

func (c *Connection) Request(m *message.Pack) chan *message.Pack {
	ch := make(chan *message.Pack)
	c.lastID++

	m.Header.ID = uint32(c.lastID)

	c.m.Lock()
	c.requests[c.lastID] = ch
	c.m.Unlock()

	return ch
}

func (c *Connection) Response(m *message.Pack) {
	c.chWrite <- m
}

func (c *Connection) Close() {
	c.s.Close()
	c.cancel()
}

func (c *Connection) read(ctx context.Context) {
	defer c.cancel()

	for {
		var p message.Pack

		err := c.dec.Decode(&p)
		if err != nil {
			fmt.Println(err)
			return
		}

		if p.Header.Type == message.Return {
			c.routeResponse(p)
		} else {
			c.routeRemoteRequest(p)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (c *Connection) routeResponse(p message.Pack) {
	c.m.Lock()
	v, ok := c.requests[ReqID(p.Header.ID)]
	if ok {
		delete(c.requests, ReqID(p.Header.ID))
	}
	c.m.Unlock()

	if ok {
		v <- &p
	}
}

func (c *Connection) routeRemoteRequest(p message.Pack) {
	c.chRemoteRequest <- &p
}

func (c *Connection) write(ctx context.Context) {
	defer c.cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case pack := <-c.chWrite:
			err := c.enc.Encode(pack)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
