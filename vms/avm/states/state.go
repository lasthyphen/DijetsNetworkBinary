// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.

package states

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lasthyphen/DijetsNetworkBinary/database"
	"github.com/lasthyphen/DijetsNetworkBinary/database/prefixdb"
	"github.com/lasthyphen/DijetsNetworkBinary/vms/avm/txs"
	"github.com/lasthyphen/DijetsNetworkBinary/vms/components/djtx"
)

var (
	utxoPrefix      = []byte("utxo")
	statusPrefix    = []byte("status")
	singletonPrefix = []byte("singleton")
	txPrefix        = []byte("tx")

	_ State = &state{}
)

// State persistently maintains a set of UTXOs, transaction, statuses, and
// singletons.
type State interface {
	djtx.UTXOState
	djtx.StatusState
	djtx.SingletonState
	TxState
}

type state struct {
	djtx.UTXOState
	djtx.StatusState
	djtx.SingletonState
	TxState
}

func New(db database.Database, parser txs.Parser, metrics prometheus.Registerer) (State, error) {
	utxoDB := prefixdb.New(utxoPrefix, db)
	statusDB := prefixdb.New(statusPrefix, db)
	singletonDB := prefixdb.New(singletonPrefix, db)
	txDB := prefixdb.New(txPrefix, db)

	utxoState, err := djtx.NewMeteredUTXOState(utxoDB, parser.Codec(), metrics)
	if err != nil {
		return nil, err
	}

	statusState, err := djtx.NewMeteredStatusState(statusDB, metrics)
	if err != nil {
		return nil, err
	}

	txState, err := NewTxState(txDB, parser, metrics)
	return &state{
		UTXOState:      utxoState,
		StatusState:    statusState,
		SingletonState: djtx.NewSingletonState(singletonDB),
		TxState:        txState,
	}, err
}
