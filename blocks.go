package main

import (
	"bodychains/blockchain"
	"bodychains/nbody"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/robaho/fixed"
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
		enc.Encode(&group)

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
		dec.Decode(&group)

		copy(univ[begin:end], group[:])
	}
}

func startBlockchain() {
	perChunk := int64(blockchain.BodiesMax / blockchain.ChunksTotal)
	totalChunks := int64(blockchain.ChunksTotal)

	univ := make(nbody.Universe, blockchain.BodiesMax)
	chunks := make(blockchain.Chunks, blockchain.ChunksTotal)

	dbFile := os.Args[1]
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

	fmt.Println("PRINT!")
}
