// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package getter

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/lasthyphen/DijetsNetworkBinary/ids"
	"github.com/lasthyphen/DijetsNetworkBinary/snow"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/choices"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/consensus/snowman"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/engine/common"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/engine/snowman/block"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/engine/snowman/block/mocks"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/validators"
)

var errUnknownBlock = errors.New("unknown block")

type StateSyncEnabledMock struct {
	*block.TestVM
	*mocks.MockStateSyncableVM
}

func testSetup(
	t *testing.T,
	ctrl *gomock.Controller,
) (StateSyncEnabledMock, *common.SenderTest, common.Config) {
	ctx := snow.DefaultConsensusContextTest()

	peers := validators.NewSet()
	sender := &common.SenderTest{}
	vm := StateSyncEnabledMock{
		TestVM:              &block.TestVM{},
		MockStateSyncableVM: mocks.NewMockStateSyncableVM(ctrl),
	}

	sender.T = t

	sender.Default(true)

	isBootstrapped := false
	subnet := &common.SubnetTest{
		T:               t,
		IsBootstrappedF: func() bool { return isBootstrapped },
		BootstrappedF:   func(ids.ID) { isBootstrapped = true },
	}

	sender.CantSendGetAcceptedFrontier = false

	peer := ids.GenerateTestNodeID()
	if err := peers.AddWeight(peer, 1); err != nil {
		t.Fatal(err)
	}

	commonConfig := common.Config{
		Ctx:                            ctx,
		Validators:                     peers,
		Beacons:                        peers,
		SampleK:                        peers.Len(),
		Alpha:                          peers.Weight()/2 + 1,
		Sender:                         sender,
		Subnet:                         subnet,
		Timer:                          &common.TimerTest{},
		AncestorsMaxContainersSent:     2000,
		AncestorsMaxContainersReceived: 2000,
		SharedCfg:                      &common.SharedConfig{},
	}

	return vm, sender, commonConfig
}

func TestAcceptedFrontier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vm, sender, config := testSetup(t, ctrl)

	blkID := ids.GenerateTestID()

	dummyBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV:     blkID,
			StatusV: choices.Accepted,
		},
		HeightV: 0,
		BytesV:  []byte{1, 2, 3},
	}
	vm.CantLastAccepted = false
	vm.LastAcceptedF = func() (ids.ID, error) { return blkID, nil }
	vm.GetBlockF = func(bID ids.ID) (snowman.Block, error) {
		assert.Equal(t, blkID, bID)
		return dummyBlk, nil
	}

	bsIntf, err := New(vm, config)
	if err != nil {
		t.Fatal(err)
	}
	bs, ok := bsIntf.(*getter)
	if !ok {
		t.Fatal("Unexpected get handler")
	}

	var accepted []ids.ID
	sender.SendAcceptedFrontierF = func(_ ids.NodeID, _ uint32, frontier []ids.ID) {
		accepted = frontier
	}

	if err := bs.GetAcceptedFrontier(ids.EmptyNodeID, 0); err != nil {
		t.Fatal(err)
	}

	if len(accepted) != 1 {
		t.Fatalf("Only one block should be accepted")
	}
	if accepted[0] != blkID {
		t.Fatalf("Blk should be accepted")
	}
}

func TestFilterAccepted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vm, sender, config := testSetup(t, ctrl)

	blkID0 := ids.GenerateTestID()
	blkID1 := ids.GenerateTestID()
	blkID2 := ids.GenerateTestID()

	blk0 := &snowman.TestBlock{TestDecidable: choices.TestDecidable{
		IDV:     blkID0,
		StatusV: choices.Accepted,
	}}
	blk1 := &snowman.TestBlock{TestDecidable: choices.TestDecidable{
		IDV:     blkID1,
		StatusV: choices.Accepted,
	}}

	vm.CantLastAccepted = false
	vm.LastAcceptedF = func() (ids.ID, error) { return blk1.ID(), nil }
	vm.GetBlockF = func(blkID ids.ID) (snowman.Block, error) {
		assert.Equal(t, blk1.ID(), blkID)
		return blk1, nil
	}

	bsIntf, err := New(vm, config)
	if err != nil {
		t.Fatal(err)
	}
	bs, ok := bsIntf.(*getter)
	if !ok {
		t.Fatal("Unexpected get handler")
	}

	blkIDs := []ids.ID{blkID0, blkID1, blkID2}
	vm.GetBlockF = func(blkID ids.ID) (snowman.Block, error) {
		switch blkID {
		case blkID0:
			return blk0, nil
		case blkID1:
			return blk1, nil
		case blkID2:
			return nil, errUnknownBlock
		}
		t.Fatal(errUnknownBlock)
		return nil, errUnknownBlock
	}

	var accepted []ids.ID
	sender.SendAcceptedF = func(_ ids.NodeID, _ uint32, frontier []ids.ID) {
		accepted = frontier
	}

	if err := bs.GetAccepted(ids.EmptyNodeID, 0, blkIDs); err != nil {
		t.Fatal(err)
	}

	acceptedSet := ids.Set{}
	acceptedSet.Add(accepted...)

	if acceptedSet.Len() != 2 {
		t.Fatalf("Two blocks should be accepted")
	}
	if !acceptedSet.Contains(blkID0) {
		t.Fatalf("Blk should be accepted")
	}
	if !acceptedSet.Contains(blkID1) {
		t.Fatalf("Blk should be accepted")
	}
	if acceptedSet.Contains(blkID2) {
		t.Fatalf("Blk shouldn't be accepted")
	}
}