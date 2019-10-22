package main

import (
	"bodychains/blockchain"
	"bodychains/message"
	"bodychains/nbody"
	"bodychains/requests"
	"bodychains/stream"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"time"

	"bodychains/connection"

	"github.com/robaho/fixed"
	"golang.org/x/net/context"
)

func fillChunks(chunks blockchain.Chunks, univ nbody.Universe, perChunk, count int64) {
	var buf bytes.Buffer

	for i := int64(0); i < count; i++ {
		begin := i * perChunk
		end := begin + perChunk

		// Prepare new chunk
		group := univ[begin:end]
		hash := nbody.HashUniverse(group, begin)

		var chunkHash blockchain.ChunkHash
		copy(chunkHash[:], hash)

		// Serialize Chunk
		buf.Reset()
		enc := gob.NewEncoder(&buf)

		err := enc.Encode(&group)
		if err != nil {
			panic(err)
		}

		chunks[blockchain.ChunkKey{Index: i, Hash: chunkHash}] = buf.Bytes()
	}
}

func fillUniverse(univ nbody.Universe, chunks blockchain.Chunks, perChunk, count int64) {
	for k, v := range chunks {
		begin := k.Index * perChunk
		end := begin + perChunk

		var group nbody.Universe

		r := bytes.NewReader(v)
		dec := gob.NewDecoder(r)

		err := dec.Decode(&group)
		if err != nil {
			panic(err)
		}

		copy(univ[begin:end], group[:])
	}
}

var listConns []*connection.Connection = []*connection.Connection{}

func pollAll(smgr *stream.StreamsManager) {
	for _, v := range listConns {
		var chainHead requests.GetChainHead
		err := requests.Run(&chainHead, message.GetChainHead, v)
		if err != nil {
			panic(err)
		}
		fmt.Println(chainHead)
	}
}

func startBlockchain(ctx context.Context, smgr *stream.StreamsManager) {
	perChunk := int64(blockchain.BodiesMax / blockchain.ChunksTotal)
	totalChunks := int64(blockchain.ChunksTotal)

	univ := make(nbody.Universe, blockchain.BodiesMax)
	chunks := make(blockchain.Chunks, blockchain.ChunksTotal)

	dbFile := "d0.db"
	if len(os.Args) == 2 {
		dbFile = os.Args[1]
	}

	bc := blockchain.NewBlockchain(dbFile)

	head := bc.Iterator().Next()

	// Create genesis block if its a new chain
	if head == nil {
		fillChunks(chunks, univ,
			perChunk,
			totalChunks,
		)

		bc.AddBlock(chunks)
		head = bc.Iterator().Next()
	}

	univ2 := make(nbody.Universe, blockchain.BodiesMax)

	fillUniverse(univ2, head.Data,
		perChunk,
		totalChunks,
	)

	dt := fixed.NewI(1, 4)
	nbody.StepVelocity(dt, univ, univ)
	nbody.ApplyVelocity(dt, univ)

	fillChunks(chunks, univ,
		perChunk,
		totalChunks,
	)
	bc.AddBlock(chunks)
	head = bc.Iterator().Next()

	go func() {
		for {
			fmt.Println("poll")
			pollAll(smgr)

			select {
			case <-time.After(time.Second * 5):
			case <-ctx.Done():
				return
			}
		}
	}()

	fmt.Println("PRINT!")
}
