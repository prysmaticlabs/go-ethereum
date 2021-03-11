package slasher

import (
	"context"
	"testing"

	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	dbtest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/event"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func Test_processQueuedBlocks_DetectsDoubleProposals(t *testing.T) {
	hook := logTest.NewGlobal()
	beaconDB := dbtest.SetupDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		serviceCfg: &ServiceConfig{
			Database: beaconDB,
		},
		params:                DefaultParams(),
		blksQueue:             newBlocksQueue(),
		proposerSlashingsFeed: new(event.Feed),
	}
	currentEpochChan := make(chan types.Epoch)
	exitChan := make(chan struct{})
	go func() {
		s.processQueuedBlocks(ctx, currentEpochChan)
		exitChan <- struct{}{}
	}()
	s.blksQueue.extend([]*slashertypes.SignedBlockHeaderWrapper{
		createProposalWrapper(4, 1, []byte{1}),
		createProposalWrapper(4, 1, []byte{1}),
		createProposalWrapper(4, 1, []byte{1}),
		createProposalWrapper(4, 1, []byte{2}),
	})
	currentEpoch := types.Epoch(0)
	currentEpochChan <- currentEpoch
	cancel()
	<-exitChan
	require.LogsContain(t, hook, "Proposer double proposal slashing")
}

func Test_processQueuedBlocks_NotSlashable(t *testing.T) {
	hook := logTest.NewGlobal()
	beaconDB := dbtest.SetupDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		serviceCfg: &ServiceConfig{
			Database: beaconDB,
		},
		params:    DefaultParams(),
		blksQueue: newBlocksQueue(),
	}
	currentEpochChan := make(chan types.Epoch)
	exitChan := make(chan struct{})
	go func() {
		s.processQueuedBlocks(ctx, currentEpochChan)
		exitChan <- struct{}{}
	}()
	s.blksQueue.extend([]*slashertypes.SignedBlockHeaderWrapper{
		createProposalWrapper(4, 1, []byte{1}),
		createProposalWrapper(4, 1, []byte{1}),
	})
	currentEpoch := types.Epoch(4)
	currentEpochChan <- currentEpoch
	cancel()
	<-exitChan
	require.LogsDoNotContain(t, hook, "Proposer double proposal slashing")
}

func createProposalWrapper(slot types.Slot, proposerIndex types.ValidatorIndex, signingRoot []byte) *slashertypes.SignedBlockHeaderWrapper {
	signRoot := bytesutil.ToBytes32(signingRoot)
	if signingRoot == nil {
		signRoot = params.BeaconConfig().ZeroHash
	}
	return &slashertypes.SignedBlockHeaderWrapper{
		SignedBeaconBlockHeader: &ethpb.SignedBeaconBlockHeader{
			Header: &ethpb.BeaconBlockHeader{
				Slot:          slot,
				ProposerIndex: proposerIndex,
				ParentRoot:    params.BeaconConfig().ZeroHash[:],
				StateRoot:     params.BeaconConfig().ZeroHash[:],
				BodyRoot:      params.BeaconConfig().ZeroHash[:],
			},
			Signature: params.BeaconConfig().EmptySignature[:],
		},
		SigningRoot: signRoot,
	}
}
