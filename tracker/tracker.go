package tracker

import (
	"github.com/btcsuite/btcd/wire"

	"github.com/CIRCL/pbtc/adaptor"
	"github.com/CIRCL/pbtc/parmap"
)

type Tracker struct {
	blocks *parmap.ParMap
	txs    *parmap.ParMap
	log    adaptor.Log
}

func New(options ...func(*Tracker)) (*Tracker, error) {
	tracker := &Tracker{
		blocks: parmap.New(),
		txs:    parmap.New(),
	}

	return tracker, nil
}

func SetLog(log adaptor.Log) func(*Tracker) {
	return func(tracker *Tracker) {
		tracker.log = log
	}
}

func (tracker *Tracker) AddTx(hash wire.ShaHash) {
	tracker.txs.Insert(hash)
}

func (tracker *Tracker) KnowsTx(hash wire.ShaHash) bool {
	return tracker.txs.Has(hash)
}

func (tracker *Tracker) AddBlock(hash wire.ShaHash) {
	tracker.blocks.Insert(hash)
}

func (tracker *Tracker) KnowsBlock(hash wire.ShaHash) bool {
	return tracker.blocks.Has(hash)
}