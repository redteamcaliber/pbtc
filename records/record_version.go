package records

import (
	"bytes"
	"net"
	"strconv"
	"time"

	"github.com/CIRCL/pbtc/util"

	"github.com/btcsuite/btcd/wire"
)

type VersionRecord struct {
	ra       *net.TCPAddr
	la       *net.TCPAddr
	stamp    time.Time
	cmd      string
	version  int32
	services uint64
	sent     time.Time
	raddr    *net.TCPAddr
	laddr    *net.TCPAddr
	agent    string
	block    int32
	relay    bool
	nonce    uint64
}

func NewVersionRecord(msg *wire.MsgVersion, ra *net.TCPAddr,
	la *net.TCPAddr) *VersionRecord {
	vr := &VersionRecord{
		stamp:    time.Now(),
		ra:       ra,
		la:       la,
		cmd:      msg.Command(),
		version:  msg.ProtocolVersion,
		services: uint64(msg.Services),
		sent:     msg.Timestamp,
		raddr:    util.ParseNetAddress(&msg.AddrYou),
		laddr:    util.ParseNetAddress(&msg.AddrMe),
		agent:    msg.UserAgent,
		block:    msg.LastBlock,
		relay:    !msg.DisableRelayTx,
		nonce:    msg.Nonce,
	}

	return vr
}

func (vr *VersionRecord) Address() *net.TCPAddr {
	return vr.ra
}

func (vr *VersionRecord) Cmd() string {
	return vr.cmd
}

func (vr *VersionRecord) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(vr.stamp.Format(time.RFC3339Nano))
	buf.WriteString(Delimiter1)
	buf.WriteString(vr.cmd)
	buf.WriteString(Delimiter1)
	buf.WriteString(vr.ra.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(vr.la.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(vr.version), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatUint(vr.services, 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(vr.sent.Unix(), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(vr.raddr.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(vr.laddr.String())
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatInt(int64(vr.block), 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatBool(vr.relay))
	buf.WriteString(Delimiter1)
	buf.WriteString(strconv.FormatUint(vr.nonce, 10))
	buf.WriteString(Delimiter1)
	buf.WriteString(vr.agent)

	return buf.String()
}
