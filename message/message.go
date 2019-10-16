package message

import (
	"bodychains/blockchain"
	"time"
)

type MessageEnum uint8

// MessageEnum values
const (
	_ MessageEnum = iota

	NotifyBeginChunk
	NotifyDoneChunk
	NotifyDoneBlock

	ReqChain
	ReqStateByHash
	ReqChunkByHash
)

type (
	Header struct {
		Timestamp int64
		Type      MessageEnum
	}

	ValueChain struct {
		blocks []blockchain.Block
	}

	GetChainArgs struct {
		HeadHash blockchain.BlockHash
	}
)

func init() {
	// gob.Register(Header{})
	// gob.Register(GetChainArgs{})
}

func NewHeader(mtype MessageEnum) *Header {
	return &Header{
		Timestamp: time.Now().UnixNano(),
		Type:      mtype,
	}
}
