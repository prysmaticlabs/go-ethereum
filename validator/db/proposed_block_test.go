package db

import (
	"testing"

	"github.com/prysmaticlabs/prysm/bazel-prysm/external/com_github_gogo_protobuf/proto"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
)

func TestCreateBlock(t *testing.T) {
	handmakeBlock := &pbp2p.BeaconBlock{Slot: 42}
	blockEnc, err := proto.Marshal(handmakeBlock)
	createdBlock, err := createBlock(blockEnc)
	if err != nil {
		t.Fatalf("failed to unmarshal encoding: %v", err)
	}
	if createdBlock.Slot != 42 {
		t.Fatal("incorrect block marshal/unmarshal")
	}
}

func TestSaveAndGetProposedBlock(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	fork := &pbp2p.Fork{}
	pubKey := getRandPubKey(t)
	block := &pbp2p.BeaconBlock{Slot: 42}

	err := db.SaveProposedBlock(fork, pubKey, block)
	if err != nil {
		t.Fatalf("can't save attestation: %v", err)
	}
	loadedProposedBlock, err := db.GetProposedBlock(fork, pubKey, block.Slot/params.BeaconConfig().SlotsPerEpoch)
	if err != nil {
		t.Fatalf("can't read attestation: %v", err)
	}

	if loadedProposedBlock.Slot != 42 {
		t.Fatalf("read the wrong attestation")
	}
}

func TestGetMaxProposedBlockEpoch(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	fork := &pbp2p.Fork{}
	pubKey := getRandPubKey(t)
	// if there were no saves, then 0 is returned
	maxProposedBlockEpoch, err := db.getMaxProposedEpoch(pubKey)
	if err != nil {
		t.Fatalf("can't get max proposed block epoch: %v", err)
	}
	if maxProposedBlockEpoch != 0 {
		t.Fatalf("getMaxProposedEpoch for new key return not 0")
	}

	// for multiple saves, the maximum epoch is returned
	block := &pbp2p.BeaconBlock{Slot: 1}
	err = db.SaveProposedBlock(fork, pubKey, block)
	if err != nil {
		t.Fatalf("can't get max proposed block epoch: %v", err)
	}
	block = &pbp2p.BeaconBlock{Slot: 10 * params.BeaconConfig().SlotsPerEpoch}
	err = db.SaveProposedBlock(fork, pubKey, block)
	if err != nil {
		t.Fatalf("can't get max proposed block epoch: %v", err)
	}
	maxProposedBlockEpoch, err = db.getMaxProposedEpoch(pubKey)
	if err != nil {
		t.Fatalf("can't get max proposed block epoch: %v", err)
	}
	if maxProposedBlockEpoch != 10 {
		t.Fatalf("getMaxProposedEpoch return not max epoch")
	}

	// maximum epoch returns to independence from the order of save
	block = &pbp2p.BeaconBlock{Slot: 5 * params.BeaconConfig().SlotsPerEpoch}
	err = db.SaveProposedBlock(fork, pubKey, block)
	if err != nil {
		t.Fatalf("can't get max proposed block epoch: %v", err)
	}
	maxProposedBlockEpoch, err = db.getMaxProposedEpoch(pubKey)
	if err != nil {
		t.Fatalf("can't get max proposed block epoch: %v", err)
	}
	if maxProposedBlockEpoch != 10 {
		t.Fatalf("getMaxProposedEpoch return not max epoch")
	}
}

func TestInitLastProposedBlockEpoch(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	fork := &pbp2p.Fork{}
	pubKey := getRandPubKey(t)

	_, exists := db.lastProposedBlockEpoch[(*pubKey)]
	if exists {
		t.Fatalf("expect that lastProposedBlockEpoch will not be initiated")
	}

	block := &pbp2p.BeaconBlock{Slot: 10 * params.BeaconConfig().SlotsPerEpoch}
	_ = db.SaveProposedBlock(fork, pubKey, block)

	lastProposedBlockEpoch := db.lastProposedBlockEpoch[(*pubKey)]
	if lastProposedBlockEpoch != 10 {
		t.Fatalf("SaveProposedBlock did not init lastProposedBlockEpoch")
	}

	_ = db.Close()
	newDB, _ := NewDB(db.DatabasePath)
	defer teardownDB(t, newDB)

	_, exists = newDB.lastProposedBlockEpoch[(*pubKey)]
	if exists {
		t.Fatalf("expect that lastProposedBlockEpoch will not be initiated in newDB")
	}

	_, _ = newDB.GetProposedBlock(fork, pubKey, 0)
	lastProposedBlockEpoch = newDB.lastProposedBlockEpoch[(*pubKey)]
	if lastProposedBlockEpoch != 10 {
		t.Fatalf("GetProposedBlock did not init lastProposedBlockEpoch")
	}
}

func TestGetProposedBlock_DontReadDiskIfLastProposedBlockEpochLess(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	fork := &pbp2p.Fork{}
	pubKey := getRandPubKey(t)
	block := &pbp2p.BeaconBlock{Slot: 10 * params.BeaconConfig().SlotsPerEpoch}

	err := db.SaveProposedBlock(fork, pubKey, block)
	if err != nil {
		t.Fatalf("can't save attestation: %v", err)
	}

	db.lastProposedBlockEpoch[(*pubKey)] = 5

	loadedProposedBlock, err := db.GetProposedBlock(fork, pubKey, block.Slot/params.BeaconConfig().SlotsPerEpoch)
	if err != nil {
		t.Fatalf("can't read attestation: %v", err)
	}

	if loadedProposedBlock != nil {
		t.Fatalf("read the block, although the lastProposedBlockEpoch said that it is not necessary to do this")
	}
}
