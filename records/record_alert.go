package records

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"net"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/wire"
)

type AlertRecord struct {
	stamp time.Time
	ra    *net.TCPAddr
	la    *net.TCPAddr
	cmd   string

	version    int32
	relayUntil int64
	expiration int64
	id         int32
	cancel     int32
	minVer     int32
	maxVer     int32
	priority   int32
	setCancel  []int32
	setSubVer  []string
	comment    string
	statusBar  string
	reserved   string
}

func NewAlertRecord(msg *wire.MsgAlert, ra *net.TCPAddr,
	la *net.TCPAddr) *AlertRecord {
	record := &AlertRecord{
		stamp: time.Now(),
		ra:    ra,
		la:    la,
		cmd:   msg.Command(),

		version:    msg.Payload.Version,
		relayUntil: msg.Payload.RelayUntil,
		expiration: msg.Payload.Expiration,
		id:         msg.Payload.ID,
		cancel:     msg.Payload.Cancel,
		minVer:     msg.Payload.MinVer,
		maxVer:     msg.Payload.MaxVer,
		priority:   msg.Payload.Priority,
		setCancel:  msg.Payload.SetCancel,
		setSubVer:  msg.Payload.SetSubVer,
		comment:    msg.Payload.Comment,
		statusBar:  msg.Payload.StatusBar,
		reserved:   msg.Payload.Reserved,
	}

	return record
}

func (ar *AlertRecord) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(ar.stamp.Format(time.RFC3339Nano))
	buf.WriteString(Delimiter1)
	buf.WriteString(ar.cmd)
	buf.WriteString(Delimiter1)
	buf.WriteString(ar.ra.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(ar.la.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.version), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.relayUntil), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.expiration), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.id), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.cancel), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.minVer), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.maxVer), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(ar.priority), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(len(ar.setCancel)), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(len(ar.setSubVer)), 10))

	buf.WriteString(Delimiter2)
	for i, cancel := range ar.setCancel {
		if i != 0 {
			buf.WriteString(Delimiter3)
		}

		buf.WriteString(strconv.FormatInt(int64(cancel), 10))
	}

	buf.WriteString(Delimiter2)
	for i, subver := range ar.setSubVer {
		if i != 0 {
			buf.WriteString(Delimiter3)
		}

		buf.WriteString(base64.StdEncoding.EncodeToString([]byte(subver)))
	}

	buf.WriteString(Delimiter2)
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(ar.comment)))

	buf.WriteString(Delimiter2)
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(ar.statusBar)))

	buf.WriteString(Delimiter2)
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(ar.reserved)))

	return buf.String()
}

func (ar *AlertRecord) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ParseCommand(ar.cmd))      //  1
	binary.Write(buf, binary.LittleEndian, ar.stamp.UnixNano())       //  8
	binary.Write(buf, binary.LittleEndian, ar.ra.IP.To16())           // 16
	binary.Write(buf, binary.LittleEndian, uint16(ar.ra.Port))        //  2
	binary.Write(buf, binary.LittleEndian, ar.la.IP.To16())           // 16
	binary.Write(buf, binary.LittleEndian, uint16(ar.la.Port))        //  2
	binary.Write(buf, binary.LittleEndian, ar.version)                //  4
	binary.Write(buf, binary.LittleEndian, ar.relayUntil)             //  8
	binary.Write(buf, binary.LittleEndian, ar.expiration)             //  8
	binary.Write(buf, binary.LittleEndian, ar.id)                     //  4
	binary.Write(buf, binary.LittleEndian, ar.cancel)                 //  4
	binary.Write(buf, binary.LittleEndian, ar.minVer)                 //  4
	binary.Write(buf, binary.LittleEndian, ar.maxVer)                 //  4
	binary.Write(buf, binary.LittleEndian, ar.priority)               //  4
	binary.Write(buf, binary.LittleEndian, uint32(len(ar.comment)))   //  4
	binary.Write(buf, binary.LittleEndian, uint32(len(ar.statusBar))) //  4
	binary.Write(buf, binary.LittleEndian, uint32(len(ar.reserved)))  //  4
	binary.Write(buf, binary.LittleEndian, uint32(len(ar.setCancel))) //  4
	binary.Write(buf, binary.LittleEndian, uint32(len(ar.setSubVer))) //  4

	binary.Write(buf, binary.LittleEndian, ar.comment)   //  X
	binary.Write(buf, binary.LittleEndian, ar.statusBar) //  Y
	binary.Write(buf, binary.LittleEndian, ar.reserved)  //  Z

	for _, cancel := range ar.setCancel { // N
		binary.Write(buf, binary.LittleEndian, cancel) // 4
	}

	for _, subver := range ar.setSubVer { // M
		binary.Write(buf, binary.LittleEndian, uint32(len(subver))) // 4
		binary.Write(buf, binary.LittleEndian, subver)              // V
	}

	// total: 105 + X + Y + Z + N*4 + M*(4+V)
	return buf.Bytes()
}