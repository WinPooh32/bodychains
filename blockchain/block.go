package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// BodiesMax ...
// PeersMax ...
// ChunksTotal ...
const (
	BodiesMax   = 2E4
	PeersMax    = 10
	ChunksTotal = BodiesMax / PeersMax
)

type (
	BlockHash [32]byte
	// Chunk ...
	Chunk []byte
	// ChunkHash ...
	ChunkHash [32]byte

	ChunkKey struct {
		Index int64
		Hash  ChunkHash
	}

	// Chunks stores chunks by it's hash
	Chunks map[ChunkKey]Chunk
)

func init() {
	gob.Register(Chunk{})
	gob.Register(ChunkHash{})
	gob.Register(ChunkKey{})
	gob.Register(Chunks{})
	gob.Register(Block{})
}

// Block keeps block headers
type Block struct {
	Timestamp     int64
	Data          Chunks
	PrevBlockHash BlockHash
	Hash          BlockHash
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// NewBlock creates and returns Block
func NewBlock(data Chunks, prevBlockHash []byte) *Block {

	block := &Block{
		Timestamp: time.Now().UnixNano(),
		Data:      data,
	}

	copy(block.PrevBlockHash[:], prevBlockHash[:])

	var buf bytes.Buffer
	for k := range data {
		buf.Write(k.Hash[:])
	}

	block.Hash = sha256.Sum256(buf.Bytes())

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(data Chunks) *Block {
	return NewBlock(data, []byte{})
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
