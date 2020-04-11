package kv

import (
	"context"
	"reflect"
	"testing"

	"github.com/gogo/protobuf/proto"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-ssz"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"gopkg.in/d4l3k/messagediff.v1"
)

func TestState_CanSaveRetrieve(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	r := [32]byte{'A'}

	if db.HasState(context.Background(), r) {
		t.Fatal("Wanted false")
	}

	st := testutil.NewBeaconState()
	st.SetSlot(100)

	if err := db.SaveState(context.Background(), st, r); err != nil {
		t.Fatal(err)
	}

	if !db.HasState(context.Background(), r) {
		t.Fatal("Wanted true")
	}

	savedS, err := db.State(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(st.InnerStateUnsafe(), savedS.InnerStateUnsafe()) {
		diff, _ := messagediff.PrettyDiff(st.InnerStateUnsafe(), savedS.InnerStateUnsafe())
		t.Errorf("Did not retrieve saved state: %v", diff)
	}

	savedS, err = db.State(context.Background(), [32]byte{'B'})
	if err != nil {
		t.Fatal(err)
	}

	if savedS != nil {
		t.Error("Unsaved state should've been nil")
	}
}

func TestHeadState_CanSaveRetrieve(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	headRoot := [32]byte{'A'}

	st := testutil.NewBeaconState()
	st.SetSlot(100)

	if err := db.SaveState(context.Background(), st, headRoot); err != nil {
		t.Fatal(err)
	}

	if err := db.SaveHeadBlockRoot(context.Background(), headRoot); err != nil {
		t.Fatal(err)
	}

	savedHeadS, err := db.HeadState(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(st.InnerStateUnsafe(), savedHeadS.InnerStateUnsafe()) {
		t.Error("did not retrieve saved state")
	}
}

func TestGenesisState_CanSaveRetrieve(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	headRoot := [32]byte{'B'}

	st := testutil.NewBeaconState()
	st.SetSlot(1)

	if err := db.SaveGenesisBlockRoot(context.Background(), headRoot); err != nil {
		t.Fatal(err)
	}

	if err := db.SaveState(context.Background(), st, headRoot); err != nil {
		t.Fatal(err)
	}

	savedGenesisS, err := db.GenesisState(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(st.InnerStateUnsafe(), savedGenesisS.InnerStateUnsafe()) {
		t.Error("did not retrieve saved state")
	}

	if err := db.SaveGenesisBlockRoot(context.Background(), [32]byte{'C'}); err != nil {
		t.Fatal(err)
	}

	savedGenesisS, err = db.HeadState(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if savedGenesisS != nil {
		t.Error("unsaved genesis state should've been nil")
	}
}

func TestStore_StatesBatchDelete(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)
	ctx := context.Background()
	numBlocks := 100
	totalBlocks := make([]*ethpb.SignedBeaconBlock, numBlocks)
	blockRoots := make([][32]byte, 0)
	evenBlockRoots := make([][32]byte, 0)
	for i := 0; i < len(totalBlocks); i++ {
		totalBlocks[i] = &ethpb.SignedBeaconBlock{
			Block: &ethpb.BeaconBlock{
				Slot:       uint64(i),
				ParentRoot: []byte("parent"),
			},
		}
		r, err := ssz.HashTreeRoot(totalBlocks[i].Block)
		if err != nil {
			t.Fatal(err)
		}
		st := testutil.NewBeaconState()
		st.SetSlot(uint64(i))
		if err := db.SaveState(context.Background(), st, r); err != nil {
			t.Fatal(err)
		}
		blockRoots = append(blockRoots, r)
		if i%2 == 0 {
			evenBlockRoots = append(evenBlockRoots, r)
		}
	}
	if err := db.SaveBlocks(ctx, totalBlocks); err != nil {
		t.Fatal(err)
	}
	// We delete all even indexed states.
	if err := db.DeleteStates(ctx, evenBlockRoots); err != nil {
		t.Fatal(err)
	}
	// When we retrieve the data, only the odd indexed state should remain.
	for _, r := range blockRoots {
		s, err := db.State(context.Background(), r)
		if err != nil {
			t.Fatal(err)
		}
		if s == nil {
			continue
		}
		if s.Slot()%2 == 0 {
			t.Errorf("State with slot %d should have been deleted", s.Slot())
		}
	}
}

func TestStore_DeleteGenesisState(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)
	ctx := context.Background()

	genesisBlockRoot := [32]byte{'A'}
	if err := db.SaveGenesisBlockRoot(ctx, genesisBlockRoot); err != nil {
		t.Fatal(err)
	}
	st := testutil.NewBeaconState()
	st.SetSlot(100)
	if err := db.SaveState(ctx, st, genesisBlockRoot); err != nil {
		t.Fatal(err)
	}
	wantedErr := "cannot delete genesis, finalized, or head state"
	if err := db.DeleteState(ctx, genesisBlockRoot); err.Error() != wantedErr {
		t.Error("Did not receive wanted error")
	}
}

func TestStore_DeleteFinalizedState(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)
	ctx := context.Background()

	genesis := bytesutil.ToBytes32([]byte{'G', 'E', 'N', 'E', 'S', 'I', 'S'})
	if err := db.SaveGenesisBlockRoot(ctx, genesis); err != nil {
		t.Fatal(err)
	}

	blk := &ethpb.SignedBeaconBlock{
		Block: &ethpb.BeaconBlock{
			ParentRoot: genesis[:],
			Slot:       100,
		},
	}
	if err := db.SaveBlock(ctx, blk); err != nil {
		t.Fatal(err)
	}

	finalizedBlockRoot, err := ssz.HashTreeRoot(blk.Block)
	if err != nil {
		t.Fatal(err)
	}

	finalizedState := testutil.NewBeaconState()
	finalizedState.SetSlot(100)
	if err := db.SaveState(ctx, finalizedState, finalizedBlockRoot); err != nil {
		t.Fatal(err)
	}
	finalizedCheckpoint := &ethpb.Checkpoint{Root: finalizedBlockRoot[:]}
	if err := db.SaveFinalizedCheckpoint(ctx, finalizedCheckpoint); err != nil {
		t.Fatal(err)
	}
	wantedErr := "cannot delete genesis, finalized, or head state"
	if err := db.DeleteState(ctx, finalizedBlockRoot); err.Error() != wantedErr {
		t.Log(err.Error())
		t.Error("Did not receive wanted error")
	}
}

func TestStore_DeleteHeadState(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)
	ctx := context.Background()

	genesis := bytesutil.ToBytes32([]byte{'G', 'E', 'N', 'E', 'S', 'I', 'S'})
	if err := db.SaveGenesisBlockRoot(ctx, genesis); err != nil {
		t.Fatal(err)
	}

	blk := &ethpb.SignedBeaconBlock{
		Block: &ethpb.BeaconBlock{
			ParentRoot: genesis[:],
			Slot:       100,
		},
	}
	if err := db.SaveBlock(ctx, blk); err != nil {
		t.Fatal(err)
	}

	headBlockRoot, err := ssz.HashTreeRoot(blk.Block)
	if err != nil {
		t.Fatal(err)
	}
	st := testutil.NewBeaconState()
	st.SetSlot(100)
	if err := db.SaveState(ctx, st, headBlockRoot); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveHeadBlockRoot(ctx, headBlockRoot); err != nil {
		t.Fatal(err)
	}
	wantedErr := "cannot delete genesis, finalized, or head state"
	if err := db.DeleteState(ctx, headBlockRoot); err.Error() != wantedErr {
		t.Error("Did not receive wanted error")
	}
}

func TestStore_SaveDeleteState_CanGetHighest(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	b := &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 1}}
	r, _ := ssz.HashTreeRoot(b.Block)
	if err := db.SaveBlock(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	st := testutil.NewBeaconState()
	if err := db.SaveState(context.Background(), st, r); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveGenesisBlockRoot(context.Background(), r); err != nil {
		t.Error(err)
	}
	s0 := st.InnerStateUnsafe()

	b = &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 999}}
	r1, _ := ssz.HashTreeRoot(b.Block)
	if err := db.SaveBlock(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	st = testutil.NewBeaconState()
	st.SetSlot(999)
	s1 := st.InnerStateUnsafe()
	if err := db.SaveState(context.Background(), st, r1); err != nil {
		t.Fatal(err)
	}

	highest, err := db.HighestSlotStates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s1) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, s1)
	}

	b = &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 1000}}
	r2, _ := ssz.HashTreeRoot(b.Block)
	if err := db.SaveBlock(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	st = testutil.NewBeaconState()
	st.SetSlot(1000)
	s2 := st.InnerStateUnsafe()
	if err := db.SaveState(context.Background(), st, r2); err != nil {
		t.Fatal(err)
	}

	highest, err = db.HighestSlotStates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s2) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, s2)
	}

	db.DeleteState(context.Background(), r2)
	highest, err = db.HighestSlotStates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s1) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, s1)
	}

	db.DeleteState(context.Background(), r1)
	highest, err = db.HighestSlotStates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if highest[0] == nil {
		t.Fatal("returned nil state ")
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s0) {
		diff, _ := messagediff.PrettyDiff(highest[0].InnerStateUnsafe(), s0)
		t.Errorf("Did not retrieve saved state: %v", diff)
	}
}

func TestStore_SaveDeleteState_CanGetHighestBelow(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	b := &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 1}}
	r, _ := ssz.HashTreeRoot(b.Block)
	if err := db.SaveBlock(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	st := testutil.NewBeaconState()
	st.SetSlot(1)
	s0 := st.InnerStateUnsafe()
	if err := db.SaveState(context.Background(), st, r); err != nil {
		t.Fatal(err)
	}

	b = &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 100}}
	r1, _ := ssz.HashTreeRoot(b.Block)
	if err := db.SaveBlock(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	st = testutil.NewBeaconState()
	st.SetSlot(100)
	s1 := st.InnerStateUnsafe()
	if err := db.SaveState(context.Background(), st, r1); err != nil {
		t.Fatal(err)
	}

	highest, err := db.HighestSlotStates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s1) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, s1)
	}

	b = &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 1000}}
	r2, _ := ssz.HashTreeRoot(b.Block)
	if err := db.SaveBlock(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	st = testutil.NewBeaconState()
	st.SetSlot(1000)
	s2 := st.InnerStateUnsafe()

	if err := db.SaveState(context.Background(), st, r2); err != nil {
		t.Fatal(err)
	}

	highest, err = db.HighestSlotStatesBelow(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s0) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, s0)
	}

	highest, err = db.HighestSlotStatesBelow(context.Background(), 101)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s1) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, s1)
	}

	highest, err = db.HighestSlotStatesBelow(context.Background(), 1001)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), s2) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, s2)
	}
}

func TestStore_GenesisState_CanGetHighestBelow(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	genesisState := testutil.NewBeaconState()
	genesisRoot := [32]byte{'a'}
	db.SaveGenesisBlockRoot(context.Background(), genesisRoot)
	db.SaveState(context.Background(), genesisState, genesisRoot)

	b := &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 1}}
	r, _ := ssz.HashTreeRoot(b.Block)
	if err := db.SaveBlock(context.Background(), b); err != nil {
		t.Fatal(err)
	}

	st := testutil.NewBeaconState()
	st.SetSlot(1)
	if err := db.SaveState(context.Background(), st, r); err != nil {
		t.Fatal(err)
	}

	highest, err := db.HighestSlotStatesBelow(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), st.InnerStateUnsafe()) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, st.InnerStateUnsafe())
	}

	highest, err = db.HighestSlotStatesBelow(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), genesisState.InnerStateUnsafe()) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, genesisState.InnerStateUnsafe())
	}
	highest, err = db.HighestSlotStatesBelow(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(highest[0].InnerStateUnsafe(), genesisState.InnerStateUnsafe()) {
		t.Errorf("Did not retrieve saved state: %v != %v", highest, genesisState.InnerStateUnsafe())
	}
}

func initializeBeaconState(slot uint64) *pb.BeaconState {
	return &pb.BeaconState{
		Slot:                      slot,
		GenesisValidatorsRoot:     make([]byte, 32),
		BlockRoots:                make([][]byte, params.BeaconConfig().SlotsPerHistoricalRoot),
		StateRoots:                make([][]byte, params.BeaconConfig().SlotsPerHistoricalRoot),
		RandaoMixes:               make([][]byte, params.BeaconConfig().EpochsPerHistoricalVector),
		Slashings:                 make([]uint64, params.BeaconConfig().EpochsPerSlashingsVector),
		PreviousEpochAttestations: make([]*pb.PendingAttestation, 0, params.BeaconConfig().MaxAttestations*params.BeaconConfig().SlotsPerEpoch),
		CurrentEpochAttestations:  make([]*pb.PendingAttestation, 0, params.BeaconConfig().MaxAttestations*params.BeaconConfig().SlotsPerEpoch),
	}
}
