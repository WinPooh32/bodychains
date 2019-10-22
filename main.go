package main

import (
	"bodychains/stream"
	"context"
	"fmt"
	"os"
	"os/signal"
)

// func decodeGetChain(stream *stream.StreamWrap) (message.GetChainArgs, error) {
// 	var getChain message.GetChainArgs

// 	dec := stream.Dec
// 	err := dec.Decode(&getChain)

// 	return getChain, err
// }

// func route(ps *stream.PeerStream) {
// 	switch ps.H.Type {
// 	case message.ReqChain:
// 		getChain, err := decodeGetChain(ps.Wrap)
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		fmt.Println(getChain)

// 	case message.NotifyBeginChunk:
// 	default:
// 		fmt.Println("Unknown message type with id:", ps.H.Type)
// 	}

// 	ps.ReadDone <- stream.Noop{}
// }

// handler of incoming connections
func conworker(ctx context.Context, pipe stream.Pipe) {
	for {
		select {
		case peer := <-pipe.OnDisconnect:
			fmt.Println("disconnected", peer)

		case <-pipe.OnConnect:
			fmt.Println("connected")

		case <-ctx.Done():
			return
		}
	}
}

func chainWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	smgr := startNetwork(ctx)

	go conworker(ctx, smgr.Pipe)

	startBlockchain(ctx, smgr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for range c {
		// handle ^C
		cancel()
		break
	}
}
