package stategen

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gogo/protobuf/proto"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	testDB "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestStateByRoot_ColdState(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)
	service.finalizedInfo.slot = 2
	service.slotsPerArchivedPoint = 1

	b := testutil.NewBeaconBlock()
	b.Block.Slot = 1
	require.NoError(t, beaconDB.SaveBlock(ctx, b))
	bRoot, err := b.Block.HashTreeRoot()
	require.NoError(t, err)
	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(1))
	require.NoError(t, service.beaconDB.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b))
	require.NoError(t, service.beaconDB.SaveGenesisBlockRoot(ctx, bRoot))
	loadedState, err := service.StateByRoot(ctx, bRoot)
	require.NoError(t, err)
	if !proto.Equal(loadedState.InnerStateUnsafe(), beaconState.InnerStateUnsafe()) {
		t.Error("Did not correctly save state")
	}
}

func TestStateByRoot_HotStateUsingEpochBoundaryCacheNoReplay(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(10))
	blk := testutil.NewBeaconBlock()
	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Root: blkRoot[:]}))
	require.NoError(t, service.epochBoundaryStateCache.put(blkRoot, beaconState))
	loadedState, err := service.StateByRoot(ctx, blkRoot)
	require.NoError(t, err)
	assert.Equal(t, types.Slot(10), loadedState.Slot(), "Did not correctly load state")
}

func TestStateByRoot_HotStateUsingEpochBoundaryCacheWithReplay(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	blk := testutil.NewBeaconBlock()
	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.epochBoundaryStateCache.put(blkRoot, beaconState))
	targetSlot := types.Slot(10)
	targetBlock := testutil.NewBeaconBlock()
	targetBlock.Block.Slot = 11
	targetBlock.Block.ParentRoot = blkRoot[:]
	targetBlock.Block.ProposerIndex = 8
	require.NoError(t, service.beaconDB.SaveBlock(ctx, targetBlock))
	targetRoot, err := targetBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Slot: targetSlot, Root: targetRoot[:]}))
	loadedState, err := service.StateByRoot(ctx, targetRoot)
	require.NoError(t, err)
	assert.Equal(t, targetSlot, loadedState.Slot(), "Did not correctly load state")
}

func TestStateByRoot_HotStateCached(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	r := [32]byte{'A'}
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Root: r[:]}))
	service.hotStateCache.put(r, beaconState)

	loadedState, err := service.StateByRoot(ctx, r)
	require.NoError(t, err)
	if !proto.Equal(loadedState.InnerStateUnsafe(), beaconState.InnerStateUnsafe()) {
		t.Error("Did not correctly cache state")
	}
}

func TestStateByRootInitialSync_UseEpochStateCache(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	targetSlot := types.Slot(10)
	require.NoError(t, beaconState.SetSlot(targetSlot))
	blk := testutil.NewBeaconBlock()
	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.epochBoundaryStateCache.put(blkRoot, beaconState))
	loadedState, err := service.StateByRootInitialSync(ctx, blkRoot)
	require.NoError(t, err)
	assert.Equal(t, targetSlot, loadedState.Slot(), "Did not correctly load state")
}

func TestStateByRootInitialSync_UseCache(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	r := [32]byte{'A'}
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Root: r[:]}))
	service.hotStateCache.put(r, beaconState)

	loadedState, err := service.StateByRootInitialSync(ctx, r)
	require.NoError(t, err)
	if !proto.Equal(loadedState.InnerStateUnsafe(), beaconState.InnerStateUnsafe()) {
		t.Error("Did not correctly cache state")
	}
	if service.hotStateCache.has(r) {
		t.Error("Hot state cache was not invalidated")
	}
}

func TestStateByRootInitialSync_CanProcessUpTo(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	blk := testutil.NewBeaconBlock()
	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.epochBoundaryStateCache.put(blkRoot, beaconState))
	targetSlot := types.Slot(10)
	targetBlk := testutil.NewBeaconBlock()
	targetBlk.Block.Slot = 11
	targetBlk.Block.ParentRoot = blkRoot[:]
	targetRoot, err := targetBlk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveBlock(ctx, targetBlk))
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Slot: targetSlot, Root: targetRoot[:]}))

	loadedState, err := service.StateByRootInitialSync(ctx, targetRoot)
	require.NoError(t, err)
	assert.Equal(t, targetSlot, loadedState.Slot(), "Did not correctly load state")
}

func TestStateBySlot_ColdState(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)
	service.slotsPerArchivedPoint = params.BeaconConfig().SlotsPerEpoch * 2
	service.finalizedInfo.slot = service.slotsPerArchivedPoint + 1

	beaconState, pks := testutil.DeterministicGenesisState(t, 32)
	genesisStateRoot, err := beaconState.HashTreeRoot(ctx)
	require.NoError(t, err)
	genesis := blocks.NewGenesisBlock(genesisStateRoot[:])
	assert.NoError(t, beaconDB.SaveBlock(ctx, genesis))
	gRoot, err := genesis.Block.HashTreeRoot()
	require.NoError(t, err)
	assert.NoError(t, beaconDB.SaveState(ctx, beaconState, gRoot))
	assert.NoError(t, beaconDB.SaveGenesisBlockRoot(ctx, gRoot))

	b, err := testutil.GenerateFullBlock(beaconState, pks, testutil.DefaultBlockGenConfig(), 1)
	require.NoError(t, err)
	require.NoError(t, beaconDB.SaveBlock(ctx, b))
	bRoot, err := b.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, beaconDB.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, beaconDB.SaveGenesisBlockRoot(ctx, bRoot))

	r := [32]byte{}
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Slot: service.slotsPerArchivedPoint, Root: r[:]}))

	slot := types.Slot(20)
	loadedState, err := service.StateBySlot(ctx, slot)
	require.NoError(t, err)
	assert.Equal(t, slot, loadedState.Slot(), "Did not correctly save state")
}

func TestStateBySlot_HotStateDB(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	genesisStateRoot, err := beaconState.HashTreeRoot(ctx)
	require.NoError(t, err)
	genesis := blocks.NewGenesisBlock(genesisStateRoot[:])
	assert.NoError(t, beaconDB.SaveBlock(ctx, genesis))
	gRoot, err := genesis.Block.HashTreeRoot()
	require.NoError(t, err)
	assert.NoError(t, beaconDB.SaveState(ctx, beaconState, gRoot))
	assert.NoError(t, beaconDB.SaveGenesisBlockRoot(ctx, gRoot))

	slot := types.Slot(10)
	loadedState, err := service.StateBySlot(ctx, slot)
	require.NoError(t, err)
	assert.Equal(t, slot, loadedState.Slot(), "Did not correctly load state")
}

func TestLoadeStateByRoot_Cached(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	r := [32]byte{'A'}
	service.hotStateCache.put(r, beaconState)

	// This tests where hot state was already cached.
	loadedState, err := service.loadStateByRoot(ctx, r)
	require.NoError(t, err)

	if !proto.Equal(loadedState.InnerStateUnsafe(), beaconState.InnerStateUnsafe()) {
		t.Error("Did not correctly cache state")
	}
}

func TestLoadeStateByRoot_FinalizedState(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	genesisStateRoot, err := beaconState.HashTreeRoot(ctx)
	require.NoError(t, err)
	genesis := blocks.NewGenesisBlock(genesisStateRoot[:])
	assert.NoError(t, beaconDB.SaveBlock(ctx, genesis))
	gRoot, err := genesis.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Slot: 0, Root: gRoot[:]}))

	service.finalizedInfo.state = beaconState
	service.finalizedInfo.slot = beaconState.Slot()
	service.finalizedInfo.root = gRoot

	// This tests where hot state was already cached.
	loadedState, err := service.loadStateByRoot(ctx, gRoot)
	require.NoError(t, err)

	if !proto.Equal(loadedState.InnerStateUnsafe(), beaconState.InnerStateUnsafe()) {
		t.Error("Did not correctly retrieve finalized state")
	}
}

func TestLoadeStateByRoot_EpochBoundaryStateCanProcess(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	gBlk := testutil.NewBeaconBlock()
	gBlkRoot, err := gBlk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.epochBoundaryStateCache.put(gBlkRoot, beaconState))

	blk := testutil.NewBeaconBlock()
	blk.Block.Slot = 11
	blk.Block.ProposerIndex = 8
	blk.Block.ParentRoot = gBlkRoot[:]
	require.NoError(t, service.beaconDB.SaveBlock(ctx, blk))
	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Slot: 10, Root: blkRoot[:]}))

	// This tests where hot state was not cached and needs processing.
	loadedState, err := service.loadStateByRoot(ctx, blkRoot)
	require.NoError(t, err)
	assert.Equal(t, types.Slot(10), loadedState.Slot(), "Did not correctly load state")
}

func TestLoadeStateByRoot_FromDBBoundaryCase(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)

	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	gBlk := testutil.NewBeaconBlock()
	gBlkRoot, err := gBlk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.epochBoundaryStateCache.put(gBlkRoot, beaconState))

	blk := testutil.NewBeaconBlock()
	blk.Block.Slot = 11
	blk.Block.ProposerIndex = 8
	blk.Block.ParentRoot = gBlkRoot[:]
	require.NoError(t, service.beaconDB.SaveBlock(ctx, blk))
	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Slot: 10, Root: blkRoot[:]}))

	// This tests where hot state was not cached and needs processing.
	loadedState, err := service.loadStateByRoot(ctx, blkRoot)
	require.NoError(t, err)
	assert.Equal(t, types.Slot(10), loadedState.Slot(), "Did not correctly load state")
}

func TestLoadeStateBySlot_CanAdvanceSlotUsingDB(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)
	beaconState, _ := testutil.DeterministicGenesisState(t, 32)
	b := testutil.NewBeaconBlock()
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b))
	gRoot, err := b.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveGenesisBlockRoot(ctx, gRoot))
	require.NoError(t, service.beaconDB.SaveState(ctx, beaconState, gRoot))

	slot := types.Slot(10)
	loadedState, err := service.loadStateBySlot(ctx, slot)
	require.NoError(t, err)
	assert.Equal(t, slot, loadedState.Slot(), "Did not correctly load state")
}

func TestLoadeStateBySlot_CanReplayBlock(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)
	genesis, keys := testutil.DeterministicGenesisState(t, 64)
	genesisBlockRoot := bytesutil.ToBytes32(nil)
	require.NoError(t, beaconDB.SaveState(ctx, genesis, genesisBlockRoot))
	stateRoot, err := genesis.HashTreeRoot(ctx)
	require.NoError(t, err)
	genesisBlk := blocks.NewGenesisBlock(stateRoot[:])
	require.NoError(t, beaconDB.SaveBlock(ctx, genesisBlk))
	genesisBlkRoot, err := genesisBlk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, beaconDB.SaveGenesisBlockRoot(ctx, genesisBlkRoot))

	b1, err := testutil.GenerateFullBlock(genesis, keys, testutil.DefaultBlockGenConfig(), 1)
	assert.NoError(t, err)
	require.NoError(t, beaconDB.SaveBlock(ctx, b1))
	r1, err := b1.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveStateSummary(ctx, &pb.StateSummary{Slot: 1, Root: r1[:]}))
	service.hotStateCache.put(bytesutil.ToBytes32(b1.Block.ParentRoot), genesis)

	loadedState, err := service.loadStateBySlot(ctx, 2)
	require.NoError(t, err)
	assert.Equal(t, types.Slot(2), loadedState.Slot(), "Did not correctly load state")
}

func TestLastAncestorState_CanGetUsingDB(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)

	b0 := testutil.NewBeaconBlock()
	b0.Block.ParentRoot = bytesutil.PadTo([]byte{'a'}, 32)
	r0, err := b0.Block.HashTreeRoot()
	require.NoError(t, err)
	b1 := testutil.NewBeaconBlock()
	b1.Block.Slot = 1
	b1.Block.ParentRoot = bytesutil.PadTo(r0[:], 32)
	r1, err := b1.Block.HashTreeRoot()
	require.NoError(t, err)
	b2 := testutil.NewBeaconBlock()
	b2.Block.Slot = 2
	b2.Block.ParentRoot = bytesutil.PadTo(r1[:], 32)
	r2, err := b2.Block.HashTreeRoot()
	require.NoError(t, err)
	b3 := testutil.NewBeaconBlock()
	b3.Block.Slot = 3
	b3.Block.ParentRoot = bytesutil.PadTo(r2[:], 32)
	r3, err := b3.Block.HashTreeRoot()
	require.NoError(t, err)

	b1State, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, b1State.SetSlot(1))

	require.NoError(t, service.beaconDB.SaveBlock(ctx, b0))
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b1))
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b2))
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b3))
	require.NoError(t, service.beaconDB.SaveState(ctx, b1State, r1))

	lastState, err := service.lastAncestorState(ctx, r3)
	require.NoError(t, err)
	assert.Equal(t, b1State.Slot(), lastState.Slot(), "Did not get wanted state")
}

func TestLastAncestorState_CanGetUsingCache(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)

	b0 := testutil.NewBeaconBlock()
	b0.Block.ParentRoot = bytesutil.PadTo([]byte{'a'}, 32)
	r0, err := b0.Block.HashTreeRoot()
	require.NoError(t, err)
	b1 := testutil.NewBeaconBlock()
	b1.Block.Slot = 1
	b1.Block.ParentRoot = bytesutil.PadTo(r0[:], 32)
	r1, err := b1.Block.HashTreeRoot()
	require.NoError(t, err)
	b2 := testutil.NewBeaconBlock()
	b2.Block.Slot = 2
	b2.Block.ParentRoot = bytesutil.PadTo(r1[:], 32)
	r2, err := b2.Block.HashTreeRoot()
	require.NoError(t, err)
	b3 := testutil.NewBeaconBlock()
	b3.Block.Slot = 3
	b3.Block.ParentRoot = bytesutil.PadTo(r2[:], 32)
	r3, err := b3.Block.HashTreeRoot()
	require.NoError(t, err)

	b1State, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, b1State.SetSlot(1))

	require.NoError(t, service.beaconDB.SaveBlock(ctx, b0))
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b1))
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b2))
	require.NoError(t, service.beaconDB.SaveBlock(ctx, b3))
	service.hotStateCache.put(r1, b1State)

	lastState, err := service.lastAncestorState(ctx, r3)
	require.NoError(t, err)
	assert.Equal(t, b1State.Slot(), lastState.Slot(), "Did not get wanted state")
}

func TestState_HasState(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)
	s, err := testutil.NewBeaconState()
	require.NoError(t, err)
	rHit1 := [32]byte{1}
	rHit2 := [32]byte{2}
	rMiss := [32]byte{3}
	service.hotStateCache.put(rHit1, s)
	require.NoError(t, service.epochBoundaryStateCache.put(rHit2, s))

	b := testutil.NewBeaconBlock()
	rHit3, err := b.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, service.beaconDB.SaveState(ctx, s, rHit3))
	tt := []struct {
		root [32]byte
		want bool
	}{
		{rHit1, true},
		{rHit2, true},
		{rMiss, false},
		{rHit3, true},
	}
	for _, tc := range tt {
		got, err := service.HasState(ctx, tc.root)
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
}

func TestState_HasStateInCache(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB)
	s, err := testutil.NewBeaconState()
	require.NoError(t, err)
	rHit1 := [32]byte{1}
	rHit2 := [32]byte{2}
	rMiss := [32]byte{3}
	service.hotStateCache.put(rHit1, s)
	require.NoError(t, service.epochBoundaryStateCache.put(rHit2, s))

	tt := []struct {
		root [32]byte
		want bool
	}{
		{rHit1, true},
		{rHit2, true},
		{rMiss, false},
	}
	for _, tc := range tt {
		got, err := service.HasStateInCache(ctx, tc.root)
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
}

func TestState_StateByStateRoot(t *testing.T) {
	ctx := context.Background()

	// We fill state and block roots with hex representations of natural numbers starting with 1.
	// Example: 16 becomes 0x00...0f
	fillRoots := func(state *pb.BeaconState) {
		rootsLen := params.MainnetConfig().SlotsPerHistoricalRoot
		roots := make([][]byte, rootsLen)
		for i := types.Slot(0); i < rootsLen; i++ {
			roots[i] = make([]byte, 32)
		}
		for j := 0; j < len(roots); j++ {
			// Remove '0x' prefix and left-pad '0' to have 64 chars in total.
			s := fmt.Sprintf("%064s", hexutil.EncodeUint64(uint64(j))[2:])
			h, err := hexutil.Decode("0x" + s)
			require.NoError(t, err, "Failed to decode root "+s)
			roots[j] = h
		}
		state.StateRoots = roots
		state.BlockRoots = roots
	}

	headState, err := testutil.NewBeaconState(fillRoots)
	require.NoError(t, err)

	t.Run("Ok", func(t *testing.T) {
		beaconDB := testDB.SetupDB(t)
		service := New(beaconDB)
		slot := types.Slot(5)

		s := fmt.Sprintf("%064s", hexutil.EncodeUint64(uint64(slot))[2:])
		h, err := hexutil.Decode("0x" + s)
		require.NoError(t, err, "Failed to decode root "+s)
		state, err := testutil.NewBeaconState(func(state *pb.BeaconState) {
			state.Slot = slot
		})
		require.NoError(t, err)
		service.hotStateCache.put(bytesutil.ToBytes32(h), state)

		rootState, err := service.StateByStateRoot(ctx, bytesutil.ToBytes32(h), headState)
		require.NoError(t, err)
		assert.Equal(t, slot, rootState.Slot())
	})

	t.Run("State root not found", func(t *testing.T) {
		beaconDB := testDB.SetupDB(t)
		service := New(beaconDB)

		s := fmt.Sprintf("%064s", hexutil.Encode([]byte("foo"))[2:])
		h, err := hexutil.Decode("0x" + s)
		require.NoError(t, err, "Failed to decode root "+s)

		_, err = service.StateByStateRoot(ctx, bytesutil.ToBytes32(h), headState)
		assert.ErrorContains(t, "could not find state in the last 8192 state roots in head state", err)
	})
}
