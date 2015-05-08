package recorder

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/btcsuite/btcd/wire"
)

type TransactionRecord struct {
	stamp   time.Time
	ra      *net.TCPAddr
	la      *net.TCPAddr
	cmd     string
	details *DetailsRecord
}

func NewTransactionRecord(msg *wire.MsgTx, ra *net.TCPAddr,
	la *net.TCPAddr) *TransactionRecord {
	record := &TransactionRecord{
		stamp:   time.Now(),
		ra:      ra,
		la:      la,
		cmd:     msg.Command(),
		details: NewDetailsRecord(msg),
	}

	return record
}

func (tr *TransactionRecord) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(tr.cmd)
	buf.WriteString(" ")
	buf.WriteString(tr.stamp.Format(time.RFC3339Nano))
	buf.WriteString(" ")
	buf.WriteString(tr.ra.String())
	buf.WriteString(" ")
	buf.WriteString(tr.la.String())
	buf.WriteString(" ")
	buf.WriteString(tr.details.String())

	return buf.String()
}

func (tr *TransactionRecord) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, tr.stamp.UnixNano())
	binary.Write(buf, binary.LittleEndian, tr.ra.IP.To16())
	binary.Write(buf, binary.LittleEndian, uint16(tr.ra.Port))
	binary.Write(buf, binary.LittleEndian, tr.la.IP.To16())
	binary.Write(buf, binary.LittleEndian, uint16(tr.la.Port))
	binary.Write(buf, binary.LittleEndian, ParseCommand(tr.cmd))

	return buf.Bytes()
}
