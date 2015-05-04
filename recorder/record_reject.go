package recorder

import (
	"net"
	"time"

	"github.com/btcsuite/btcd/wire"
)

type RejectRecord struct {
	stamp    time.Time
	ra       *net.TCPAddr
	la       *net.TCPAddr
	msg_t    MsgType
	reject_t MsgType
	hash     [32]byte
	code     uint8
	reason   string
}

func NewRejectRecord(msg *wire.MsgReject, ra *net.TCPAddr,
	la *net.TCPAddr) *RejectRecord {
	record := &RejectRecord{
		stamp:    time.Now(),
		ra:       ra,
		la:       la,
		msg_t:    MsgReject,
		reject_t: ParseCommand(msg.Command()),
		hash:     [32]byte(msg.Hash),
		code:     uint8(msg.Code),
		reason:   msg.Reason,
	}

	return record
}
