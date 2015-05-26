package records

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/wire"
)

type PongRecord struct {
	stamp time.Time
	ra    *net.TCPAddr
	la    *net.TCPAddr
	cmd   string
	nonce uint64
}

func NewPongRecord(msg *wire.MsgPong, ra *net.TCPAddr,
	la *net.TCPAddr) *PongRecord {
	record := &PongRecord{
		stamp: time.Now(),
		ra:    ra,
		la:    la,
		cmd:   msg.Command(),
		nonce: msg.Nonce,
	}

	return record
}

func (pr *PongRecord) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(pr.stamp.Format(time.RFC3339Nano))
	buf.WriteString(Delimiter1)
	buf.WriteString(pr.cmd)
	buf.WriteString(Delimiter1)
	buf.WriteString(pr.ra.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(pr.la.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatUint(pr.nonce, 10))

	return buf.String()
}

func (pr *PongRecord) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ParseCommand(pr.cmd)) //  1
	binary.Write(buf, binary.LittleEndian, pr.stamp.UnixNano())  //  8
	binary.Write(buf, binary.LittleEndian, pr.ra.IP.To16())      // 16
	binary.Write(buf, binary.LittleEndian, uint16(pr.ra.Port))   //  2
	binary.Write(buf, binary.LittleEndian, pr.la.IP.To16())      // 16
	binary.Write(buf, binary.LittleEndian, uint16(pr.la.Port))   //  2
	binary.Write(buf, binary.LittleEndian, pr.nonce)             //  8

	// total: 53
	return buf.Bytes()
}