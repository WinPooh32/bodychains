package message

import (
	"bufio"
	"bytes"
	"encoding/gob"
)

type MessageEnum uint8

// MessageEnum values
//go:generate stringer -type=MessageEnum
const (
	_ MessageEnum = iota

	NotifyBeginChunk
	NotifyDoneChunk
	NotifyDoneBlock

	GetChain
	GetChainHead
	GetStateByHash
	GetChunkByHash

	Return
)

type Pack struct {
	Header Header
	Body   []byte
}

type Header struct {
	Type MessageEnum
	ID   uint32
}

var _MessageStrings map[string]MessageEnum

func MessageEnumFromString(raw string) (out MessageEnum, ok bool) {
	out, ok = _MessageStrings[raw]
	return out, ok
}

func init() {
	count := len(_MessageEnum_index) - 1
	_MessageStrings = make(map[string]MessageEnum, count)

	for i := 0; i < count; i++ {
		enum := MessageEnum(i)
		_MessageStrings[enum.String()] = enum
	}
}

type Request interface {
	Args() interface{}
	Resp() interface{}
}

// encode message
func DoPack(req Request, messageType MessageEnum) (m *Pack, err error) {
	buf := bytes.NewBuffer(make([]byte, 64))

	enc := gob.NewEncoder(buf)
	err = enc.Encode(req.Args())
	if err != nil {
		return m, err
	}

	// pack message to one packet
	m = &Pack{
		Header: Header{
			Type: messageType,
		},
		Body: buf.Bytes(),
	}

	return m, nil
}

// decode response as RespChainHead
func DoUnpack(req Request, m *Pack) (err error) {
	reader := bufio.NewReader(bytes.NewBuffer(m.Body))
	dec := gob.NewDecoder(reader)
	err = dec.Decode(req.Resp())
	return err
}
