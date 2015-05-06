package recorder

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/wire"
)

type RejectRecord struct {
	stamp  time.Time
	ra     *net.TCPAddr
	la     *net.TCPAddr
	cmd    string
	code   uint8
	reject string
	hash   []byte
	reason string
}

func NewRejectRecord(msg *wire.MsgReject, ra *net.TCPAddr,
	la *net.TCPAddr) *RejectRecord {
	record := &RejectRecord{
		stamp:  time.Now(),
		ra:     ra,
		la:     la,
		cmd:    msg.Command(),
		code:   uint8(msg.Code),
		reject: msg.Cmd,
		hash:   msg.Hash.Bytes(),
		reason: msg.Reason,
	}

	return record
}

func (rr *RejectRecord) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(rr.cmd)
	buf.WriteString(" ")
	buf.WriteString(rr.stamp.Format(time.RFC3339Nano))
	buf.WriteString(" ")
	buf.WriteString(rr.ra.String())
	buf.WriteString(" ")
	buf.WriteString(rr.la.String())
	buf.WriteString(" ")
	buf.WriteString(strconv.FormatInt(int64(rr.code), 10))
	buf.WriteString(" ")
	buf.WriteString(rr.reject)
	buf.WriteString(" ")
	buf.WriteString(hex.EncodeToString(rr.hash))
	buf.WriteString(" ")
	buf.WriteString(rr.reason)

	return buf.String()
}

func (rr *RejectRecord) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rr.stamp.UnixNano())
	binary.Write(buf, binary.LittleEndian, rr.ra.IP.To16())
	binary.Write(buf, binary.LittleEndian, uint16(rr.ra.Port))
	binary.Write(buf, binary.LittleEndian, rr.la.IP.To16())
	binary.Write(buf, binary.LittleEndian, uint16(rr.la.Port))
	binary.Write(buf, binary.LittleEndian, ParseCommand(rr.cmd))
	binary.Write(buf, binary.LittleEndian, rr.code)
	binary.Write(buf, binary.LittleEndian, ParseCommand(rr.reject))
	binary.Write(buf, binary.LittleEndian, rr.hash)
	binary.Write(buf, binary.LittleEndian, len(rr.reason))
	binary.Write(buf, binary.LittleEndian, rr.reason)

	return buf.Bytes()
}