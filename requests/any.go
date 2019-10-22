package requests

import (
	"bodychains/connection"
	"bodychains/message"
)


//go:generate gotemplate -outfmt "%v" "bodychains/message/template" NotifyBeginChunk()
//go:generate gotemplate -outfmt "%v" "bodychains/message/template" NotifyDoneChunk()
//go:generate gotemplate -outfmt "%v" "bodychains/message/template" NotifyDoneBlock()

//go:generate gotemplate -outfmt "%v" "bodychains/message/template" GetChain()
//go:generate gotemplate -outfmt "%v" "bodychains/message/template" GetChainHead()
//go:generate gotemplate -outfmt "%v" "bodychains/message/template" GetStateByHash()
//go:generate gotemplate -outfmt "%v" "bodychains/message/template" GetChunkByHash()


func Run(req message.Request, messageType message.MessageEnum, conn *connection.Connection) (err error) {
	// pack message
	m, err := message.DoPack(req, messageType)
	if err != nil {
		return err
	}

	// wait for response
	response := <-conn.Request(m)

	// decode response as RespChainHead
	err = message.DoUnpack(req, response)

	return err
}
