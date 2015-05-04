package recorder

import (
	"net"
	"time"

	"github.com/btcsuite/btcd/wire"
)

type GetAddrRecord struct {
	stamp time.Time
	ra    *net.TCPAddr
	la    *net.TCPAddr
	msg_t MsgType
}

func NewGetAddrRecord(msg *wire.MsgGetAddr, ra *net.TCPAddr,
	la *net.TCPAddr) *GetAddrRecord {
	record := &GetAddrRecord{
		stamp: time.Now(),
		ra:    ra,
		la:    la,
		msg_t: MsgGetAddr,
	}

	return record
}
