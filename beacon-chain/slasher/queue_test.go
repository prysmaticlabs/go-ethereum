package slasher

import (
	"testing"

	types "github.com/prysmaticlabs/eth2-types"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func Test_attestationsQueue(t *testing.T) {
	t.Run("push_and_dequeue", func(tt *testing.T) {
		attQueue := newAttestationsQueue()
		wantedAtts := []*slashertypes.IndexedAttestationWrapper{
			createAttestationWrapper(0, 1, []uint64{1}, make([]byte, 32)),
			createAttestationWrapper(1, 2, []uint64{1}, make([]byte, 32)),
		}
		attQueue.push(wantedAtts[0])
		attQueue.push(wantedAtts[1])
		require.DeepEqual(t, 2, attQueue.size())

		received := attQueue.dequeue()
		require.DeepEqual(t, 0, attQueue.size())
		require.DeepEqual(t, wantedAtts, received)
	})

	t.Run("extend_and_dequeue", func(tt *testing.T) {
		attQueue := newAttestationsQueue()
		wantedAtts := []*slashertypes.IndexedAttestationWrapper{
			createAttestationWrapper(0, 1, []uint64{1}, make([]byte, 32)),
			createAttestationWrapper(1, 2, []uint64{1}, make([]byte, 32)),
		}
		attQueue.extend(wantedAtts)
		require.DeepEqual(t, 2, attQueue.size())

		received := attQueue.dequeue()
		require.DeepEqual(t, 0, attQueue.size())
		require.DeepEqual(t, wantedAtts, received)
	})
}

func Test_blocksQueue(t *testing.T) {
	t.Run("push_and_dequeue", func(tt *testing.T) {
		blkQueue := newBlocksQueue()
		wantedBlks := []*slashertypes.SignedBlockHeaderWrapper{
			createProposalWrapper(0, types.ValidatorIndex(1), make([]byte, 32)),
			createProposalWrapper(1, types.ValidatorIndex(1), make([]byte, 32)),
		}
		blkQueue.push(wantedBlks[0])
		blkQueue.push(wantedBlks[1])
		require.DeepEqual(t, 2, blkQueue.size())

		received := blkQueue.dequeue()
		require.DeepEqual(t, 0, blkQueue.size())
		require.DeepEqual(t, wantedBlks, received)
	})

	t.Run("extend_and_dequeue", func(tt *testing.T) {
		blkQueue := newBlocksQueue()
		wantedBlks := []*slashertypes.SignedBlockHeaderWrapper{
			createProposalWrapper(0, types.ValidatorIndex(1), make([]byte, 32)),
			createProposalWrapper(1, types.ValidatorIndex(1), make([]byte, 32)),
		}
		blkQueue.extend(wantedBlks)
		require.DeepEqual(t, 2, blkQueue.size())

		received := blkQueue.dequeue()
		require.DeepEqual(t, 0, blkQueue.size())
		require.DeepEqual(t, wantedBlks, received)
	})
}
