// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lasthyphen/DijetsNetworkBinary/chains/atomic"
	"github.com/lasthyphen/DijetsNetworkBinary/database/prefixdb"
	"github.com/lasthyphen/DijetsNetworkBinary/ids"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/crypto"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/logging"
	"github.com/lasthyphen/DijetsNetworkBinary/vms/components/djtx"
	"github.com/lasthyphen/DijetsNetworkBinary/vms/platformvm/status"
	"github.com/lasthyphen/DijetsNetworkBinary/vms/secp256k1fx"
)

func TestAtomicTxImports(t *testing.T) {
	vm, baseDB, _ := defaultVM()
	vm.ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(); err != nil {
			t.Fatal(err)
		}
		vm.ctx.Lock.Unlock()
	}()
	assert := assert.New(t)

	utxoID := djtx.UTXOID{
		TxID:        ids.Empty.Prefix(1),
		OutputIndex: 1,
	}
	amount := uint64(70000)
	recipientKey := keys[1]

	m := &atomic.Memory{}
	err := m.Initialize(logging.NoLog{}, prefixdb.New([]byte{5}, baseDB))
	if err != nil {
		t.Fatal(err)
	}
	vm.ctx.SharedMemory = m.NewSharedMemory(vm.ctx.ChainID)
	vm.AtomicUTXOManager = djtx.NewAtomicUTXOManager(vm.ctx.SharedMemory, Codec)
	peerSharedMemory := m.NewSharedMemory(vm.ctx.XChainID)
	utxo := &djtx.UTXO{
		UTXOID: utxoID,
		Asset:  djtx.Asset{ID: djtxAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: amount,
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{recipientKey.PublicKey().Address()},
			},
		},
	}
	utxoBytes, err := Codec.Marshal(CodecVersion, utxo)
	if err != nil {
		t.Fatal(err)
	}
	inputID := utxo.InputID()
	if err := peerSharedMemory.Apply(map[ids.ID]*atomic.Requests{vm.ctx.ChainID: {PutRequests: []*atomic.Element{{
		Key:   inputID[:],
		Value: utxoBytes,
		Traits: [][]byte{
			recipientKey.PublicKey().Address().Bytes(),
		},
	}}}}); err != nil {
		t.Fatal(err)
	}

	tx, err := vm.newImportTx(
		vm.ctx.XChainID,
		recipientKey.PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{recipientKey},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}
	vm.internalState.SetTimestamp(vm.ApricotPhase5Time.Add(100 * time.Second))

	vm.mempool.AddDecisionTx(tx)
	b, err := vm.BuildBlock()
	assert.NoError(err)
	// Test multiple verify calls work
	err = b.Verify()
	assert.NoError(err)
	err = b.Accept()
	assert.NoError(err)
	_, txStatus, err := vm.internalState.GetTx(tx.ID())
	assert.NoError(err)
	// Ensure transaction is in the committed state
	assert.Equal(txStatus, status.Committed)
	// Ensure standard block contains one atomic transaction
	assert.Equal(b.(*StandardBlock).inputs.Len(), 1)
}