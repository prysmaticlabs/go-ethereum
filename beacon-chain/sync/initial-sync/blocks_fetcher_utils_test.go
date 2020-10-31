package initialsync

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p-core/peer"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	mock "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	dbtest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/flags"
	p2pt "github.com/prysmaticlabs/prysm/beacon-chain/p2p/testing"
	p2ppb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestBlocksFetcher_nonSkippedSlotAfter(t *testing.T) {
	peersGen := func(size int) []*peerData {
		blocks := append(makeSequence(1, 64), makeSequence(500, 640)...)
		blocks = append(blocks, makeSequence(51200, 51264)...)
		blocks = append(blocks, 55000)
		blocks = append(blocks, makeSequence(57000, 57256)...)
		var peersData []*peerData
		for i := 0; i < size; i++ {
			peersData = append(peersData, &peerData{
				blocks:         blocks,
				finalizedEpoch: 1800,
				headSlot:       57000,
			})
		}
		return peersData
	}
	chainConfig := struct {
		peers []*peerData
	}{
		peers: peersGen(5),
	}

	mc, p2p, _ := initializeTestServices(t, []uint64{}, chainConfig.peers)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fetcher := newBlocksFetcher(
		ctx,
		&blocksFetcherConfig{
			chain: mc,
			p2p:   p2p,
		},
	)
	fetcher.rateLimiter = leakybucket.NewCollector(6400, 6400, false)
	seekSlots := map[uint64]uint64{
		0:     1,
		10:    11,
		31:    32,
		32:    33,
		63:    64,
		64:    500,
		160:   500,
		352:   500,
		480:   500,
		512:   513,
		639:   640,
		640:   51200,
		6640:  51200,
		51200: 51201,
	}
	for seekSlot, expectedSlot := range seekSlots {
		t.Run(fmt.Sprintf("range: %d (%d-%d)", expectedSlot-seekSlot, seekSlot, expectedSlot), func(t *testing.T) {
			slot, err := fetcher.nonSkippedSlotAfter(ctx, seekSlot)
			assert.NoError(t, err)
			assert.Equal(t, expectedSlot, slot, "Unexpected slot")
		})
	}

	t.Run("test isolated non-skipped slot", func(t *testing.T) {
		seekSlot := uint64(51264)
		expectedSlot := uint64(55000)
		found := false
		var i int
		for i = 0; i < 100; i++ {
			slot, err := fetcher.nonSkippedSlotAfter(ctx, seekSlot)
			assert.NoError(t, err)
			if slot == expectedSlot {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Isolated non-skipped slot not found in %d iterations: %v", i, expectedSlot)
		} else {
			log.Debugf("Isolated non-skipped slot found in %d iterations", i)
		}
	})

	t.Run("no peers with higher target epoch available", func(t *testing.T) {
		peers := []*peerData{
			{finalizedEpoch: 3, headSlot: 160},
			{finalizedEpoch: 3, headSlot: 160},
			{finalizedEpoch: 3, headSlot: 160},
			{finalizedEpoch: 8, headSlot: 320},
			{finalizedEpoch: 8, headSlot: 320},
			{finalizedEpoch: 10, headSlot: 320},
			{finalizedEpoch: 10, headSlot: 640},
		}
		p2p := p2pt.NewTestP2P(t)
		connectPeers(t, p2p, peers, p2p.Peers())
		fetcher := newBlocksFetcher(
			ctx,
			&blocksFetcherConfig{
				chain: mc,
				p2p:   p2p,
			},
		)
		mc.FinalizedCheckPoint = &eth.Checkpoint{
			Epoch: 10,
		}
		require.NoError(t, mc.State.SetSlot(12*params.BeaconConfig().SlotsPerEpoch))

		fetcher.mode = modeStopOnFinalizedEpoch
		slot, err := fetcher.nonSkippedSlotAfter(ctx, 160)
		assert.ErrorContains(t, errSlotIsTooHigh.Error(), err)
		assert.Equal(t, uint64(0), slot)

		fetcher.mode = modeNonConstrained
		require.NoError(t, mc.State.SetSlot(20*params.BeaconConfig().SlotsPerEpoch))
		slot, err = fetcher.nonSkippedSlotAfter(ctx, 160)
		assert.ErrorContains(t, errSlotIsTooHigh.Error(), err)
		assert.Equal(t, uint64(0), slot)
	})
}

func TestBlocksFetcher_findFork(t *testing.T) {
	// Chain graph:
	// A - B - C - D - E
	//      \
	//       - C'- D'- E'- F'- G'
	// Allow fetcher to proceed till E, then connect peer having alternative, branch.
	// Test that G' slot can be reached i.e. fetcher can track back and explore alternative paths.
	beaconDB, _ := dbtest.SetupDB(t)
	p2p := p2pt.NewTestP2P(t)

	// Chain contains blocks from 8 epochs (from 0 to 7, 256 is the start slot of epoch8).
	chain1 := extendBlockSequence(t, []*eth.SignedBeaconBlock{}, 250)
	finalizedSlot := uint64(63)
	finalizedEpoch := helpers.SlotToEpoch(finalizedSlot)

	genesisBlock := chain1[0]
	require.NoError(t, beaconDB.SaveBlock(context.Background(), genesisBlock))
	genesisRoot, err := genesisBlock.Block.HashTreeRoot()
	require.NoError(t, err)

	st := testutil.NewBeaconState()
	mc := &mock.ChainService{
		State: st,
		Root:  genesisRoot[:],
		DB:    beaconDB,
		FinalizedCheckPoint: &eth.Checkpoint{
			Epoch: finalizedEpoch,
			Root:  []byte(fmt.Sprintf("finalized_root %d", finalizedEpoch)),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fetcher := newBlocksFetcher(
		ctx,
		&blocksFetcherConfig{
			chain: mc,
			p2p:   p2p,
			db:    beaconDB,
		},
	)
	fetcher.rateLimiter = leakybucket.NewCollector(6400, 6400, false)

	// Consume all chain1 blocks from many peers (alternative fork will be featured by a single peer,
	// and should still be enough to explore alternative paths).
	peers := make([]peer.ID, 0)
	for i := 0; i < 5; i++ {
		peers = append(peers, connectPeerHavingBlocks(t, p2p, chain1, finalizedSlot, p2p.Peers()))
	}

	blockBatchLimit := uint64(flags.Get().BlockBatchLimit) * 2
	pidInd := 0
	for i := uint64(1); i < uint64(len(chain1)); i += blockBatchLimit {
		req := &p2ppb.BeaconBlocksByRangeRequest{
			StartSlot: i,
			Step:      1,
			Count:     blockBatchLimit,
		}
		blocks, err := fetcher.requestBlocks(ctx, req, peers[pidInd%len(peers)])
		require.NoError(t, err)
		for _, blk := range blocks {
			require.NoError(t, beaconDB.SaveBlock(ctx, blk))
			require.NoError(t, st.SetSlot(blk.Block.Slot))
		}
		pidInd++
	}

	// Assert that all the blocks from chain1 are known.
	for _, blk := range chain1 {
		blkRoot, err := blk.Block.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, true, beaconDB.HasBlock(ctx, blkRoot) || mc.HasInitSyncBlock(blkRoot))
	}
	assert.Equal(t, uint64(250), mc.HeadSlot())

	// Assert no blocks on further requests, disallowing to progress.
	req := &p2ppb.BeaconBlocksByRangeRequest{
		StartSlot: 251,
		Step:      1,
		Count:     blockBatchLimit,
	}
	blocks, err := fetcher.requestBlocks(ctx, req, peers[pidInd%len(peers)])
	require.NoError(t, err)
	assert.Equal(t, 0, len(blocks))

	// If no peers with unexplored paths exist, error should be returned.
	fork, err := fetcher.findFork(ctx, 251)
	require.ErrorContains(t, errNoPeersAvailable.Error(), err)
	require.Equal(t, (*forkData)(nil), fork)

	// Add peer that has blocks after 250, but those blocks are orphaned i.e. they do not have common
	// ancestor with what we already have. So, no common ancestor exists.
	chain1a := extendBlockSequence(t, []*eth.SignedBeaconBlock{}, 265)
	connectPeerHavingBlocks(t, p2p, chain1a, finalizedSlot, p2p.Peers())
	fork, err = fetcher.findFork(ctx, 251)
	require.ErrorContains(t, errNoPeersWithAltBlocks.Error(), err)
	require.Equal(t, (*forkData)(nil), fork)

	// Add peer which has blocks after 250. It is not on another fork, but algorithm
	// is smart enough to link back to common ancestor, w/o discriminating between forks. This is
	// by design: fork exploration is undertaken when FSMs are stuck, so any progress is good.
	chain1b := extendBlockSequence(t, chain1, 64)
	curForkMoreBlocksPeer := connectPeerHavingBlocks(t, p2p, chain1b, finalizedSlot, p2p.Peers())
	fork, err = fetcher.findFork(ctx, 251)
	require.NoError(t, err)
	require.Equal(t, 2, len(fork.blocks))
	require.Equal(t, curForkMoreBlocksPeer, fork.peer)
	// Save all chain1b blocks (so that they do not interfere with alternative fork)
	for _, blk := range chain1b {
		require.NoError(t, beaconDB.SaveBlock(ctx, blk))
		require.NoError(t, st.SetSlot(blk.Block.Slot))
	}

	forkSlot := 129
	chain2 := extendBlockSequence(t, chain1[:forkSlot], 165)
	// Assert that forked blocks from chain2 are unknown.
	assert.Equal(t, 294, len(chain2))
	for _, blk := range chain2[forkSlot:] {
		blkRoot, err := blk.Block.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, false, beaconDB.HasBlock(ctx, blkRoot) || mc.HasInitSyncBlock(blkRoot))
	}
	// Search for alternative paths (add single peer having alternative path).
	alternativePeer := connectPeerHavingBlocks(t, p2p, chain2, finalizedSlot, p2p.Peers())
	fmt.Printf("altpeer: %v\n", alternativePeer)
	fork, err = fetcher.findFork(ctx, 251)
	require.NoError(t, err)
	assert.Equal(t, alternativePeer, fork.peer)
	assert.Equal(t, 65, len(fork.blocks))
	ind := forkSlot
	for _, blk := range fork.blocks {
		require.Equal(t, blk.Block.Slot, chain2[ind].Block.Slot)
		ind++
	}

	// Process returned blocks and then attempt to extend chain (ensuring that parent block exists).
	for _, blk := range fork.blocks {
		require.NoError(t, beaconDB.SaveBlock(ctx, blk))
		require.NoError(t, st.SetSlot(blk.Block.Slot))
	}
	assert.Equal(t, uint64(forkSlot+len(fork.blocks)-1), mc.HeadSlot())
	for i := forkSlot + len(fork.blocks); i < len(chain2); i++ {
		blk := chain2[i]
		require.Equal(t, blk.Block.Slot, uint64(i), "incorrect block selected for slot %d", i)
		// Only save is parent block exists.
		parentRoot := bytesutil.ToBytes32(blk.Block.ParentRoot)
		if beaconDB.HasBlock(ctx, parentRoot) || mc.HasInitSyncBlock(parentRoot) {
			require.NoError(t, beaconDB.SaveBlock(ctx, blk))
			require.NoError(t, st.SetSlot(blk.Block.Slot))
		}
	}

	// Assert that all the blocks from chain2 are known.
	for _, blk := range chain2 {
		blkRoot, err := blk.Block.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, true, beaconDB.HasBlock(ctx, blkRoot) || mc.HasInitSyncBlock(blkRoot), "slot %d", blk.Block.Slot)
	}
}

func TestBlocksFetcher_currentHeadAndTargetEpochs(t *testing.T) {
	tests := []struct {
		name               string
		syncMode           syncMode
		peers              []*peerData
		ourFinalizedEpoch  uint64
		ourHeadSlot        uint64
		expectedHeadEpoch  uint64
		targetEpoch        uint64
		targetEpochSupport int
	}{
		{
			name:               "ignore lower epoch peers in best finalized",
			syncMode:           modeStopOnFinalizedEpoch,
			ourHeadSlot:        5 * params.BeaconConfig().SlotsPerEpoch,
			expectedHeadEpoch:  4,
			ourFinalizedEpoch:  4,
			targetEpoch:        10,
			targetEpochSupport: 3,
			peers: []*peerData{
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 8, headSlot: 320},
				{finalizedEpoch: 8, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
			},
		},
		{
			name:               "resolve ties in best finalized",
			syncMode:           modeStopOnFinalizedEpoch,
			ourHeadSlot:        5 * params.BeaconConfig().SlotsPerEpoch,
			expectedHeadEpoch:  4,
			ourFinalizedEpoch:  4,
			targetEpoch:        10,
			targetEpochSupport: 3,
			peers: []*peerData{
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 8, headSlot: 320},
				{finalizedEpoch: 8, headSlot: 320},
				{finalizedEpoch: 8, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
			},
		},
		{
			name:               "best non-finalized",
			syncMode:           modeNonConstrained,
			ourHeadSlot:        5 * params.BeaconConfig().SlotsPerEpoch,
			expectedHeadEpoch:  5,
			ourFinalizedEpoch:  4,
			targetEpoch:        20,
			targetEpochSupport: 1,
			peers: []*peerData{
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 3, headSlot: 160},
				{finalizedEpoch: 8, headSlot: 320},
				{finalizedEpoch: 8, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
				{finalizedEpoch: 10, headSlot: 320},
				{finalizedEpoch: 15, headSlot: 640},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, p2p, _ := initializeTestServices(t, []uint64{}, tt.peers)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			fetcher := newBlocksFetcher(
				ctx,
				&blocksFetcherConfig{
					chain: mc,
					p2p:   p2p,
				},
			)
			mc.FinalizedCheckPoint = &eth.Checkpoint{
				Epoch: tt.ourFinalizedEpoch,
			}
			require.NoError(t, mc.State.SetSlot(tt.ourHeadSlot))
			fetcher.mode = tt.syncMode

			// Head and target epochs calculation.
			headEpoch, targetEpoch, peers := fetcher.calculateHeadAndTargetEpochs()
			assert.Equal(t, tt.expectedHeadEpoch, headEpoch, "Unexpected head epoch")
			assert.Equal(t, tt.targetEpoch, targetEpoch, "Unexpected target epoch")
			assert.Equal(t, tt.targetEpochSupport, len(peers), "Unexpected number of peers supporting target epoch")

			// Best finalized and non-finalized slots.
			finalizedSlot := tt.targetEpoch * params.BeaconConfig().SlotsPerEpoch
			if tt.syncMode == modeStopOnFinalizedEpoch {
				assert.Equal(t, finalizedSlot, fetcher.bestFinalizedSlot(), "Unexpected finalized slot")
			} else {
				assert.Equal(t, finalizedSlot, fetcher.bestNonFinalizedSlot(), "Unexpected non-finalized slot")
			}
		})
	}
}
