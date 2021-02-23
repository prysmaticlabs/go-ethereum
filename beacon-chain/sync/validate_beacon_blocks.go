package sync

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/feed"
	blockfeed "github.com/prysmaticlabs/prysm/beacon-chain/core/feed/block"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/timeutils"
	"github.com/prysmaticlabs/prysm/shared/traceutil"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// validateBeaconBlockPubSub checks that the incoming block has a valid BLS signature.
// Blocks that have already been seen are ignored. If the BLS signature is any valid signature,
// this method rebroadcasts the message.
func (s *Service) validateBeaconBlockPubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Record the timestamp to log
	receivedTime := timeutils.Now()
	// Validation runs on publish (not just subscriptions), so we should approve any message from
	// ourselves.
	if pid == s.p2p.PeerID() {
		return pubsub.ValidationAccept
	}

	// We should not attempt to process blocks until fully synced, but propagation is OK.
	if s.initialSync.Syncing() {
		return pubsub.ValidationIgnore
	}

	ctx, span := trace.StartSpan(ctx, "sync.validateBeaconBlockPubSub")
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	s.validateBlockLock.Lock()
	defer s.validateBlockLock.Unlock()

	blk, ok := m.(*ethpb.SignedBeaconBlock)
	if !ok {
		log.WithError(errors.New("msg is not ethpb.SignedBeaconBlock")).Debug("Rejected block")
		return pubsub.ValidationReject
	}

	if blk.Block == nil {
		log.WithError(errors.New("block.Block is nil")).Debug("Rejected block")
		return pubsub.ValidationReject
	}

	// Broadcast the block on a feed to notify other services in the beacon node
	// of a received block (even if it does not process correctly through a state transition).
	s.blockNotifier.BlockFeed().Send(&feed.Event{
		Type: blockfeed.ReceivedBlock,
		Data: &blockfeed.ReceivedBlockData{
			SignedBlock: blk,
		},
	})

	// Verify the block is the first block received for the proposer for the slot.
	if s.hasSeenBlockIndexSlot(blk.Block.Slot, blk.Block.ProposerIndex) {
		return pubsub.ValidationIgnore
	}

	blockRoot, err := blk.Block.HashTreeRoot()
	if err != nil {
		log.WithError(err).WithField("blockSlot", blk.Block.Slot).Debug("Ignored block")
		return pubsub.ValidationIgnore
	}
	if s.db.HasBlock(ctx, blockRoot) {
		return pubsub.ValidationIgnore
	}
	// Check if parent is a bad block and then reject the block.
	if s.hasBadBlock(bytesutil.ToBytes32(blk.Block.ParentRoot)) {
		s.setBadBlock(ctx, blockRoot)
		e := fmt.Errorf("received block with root %#x that has an invalid parent %#x", blockRoot, blk.Block.ParentRoot)
		log.WithError(e).WithField("blockSlot", blk.Block.Slot).Debug("Rejected block")
		return pubsub.ValidationReject
	}

	s.pendingQueueLock.RLock()
	if s.seenPendingBlocks[blockRoot] {
		s.pendingQueueLock.RUnlock()
		return pubsub.ValidationIgnore
	}
	s.pendingQueueLock.RUnlock()

	if err := helpers.VerifySlotTime(uint64(s.chain.GenesisTime().Unix()), blk.Block.Slot, params.BeaconNetworkConfig().MaximumGossipClockDisparity); err != nil {
		log.WithError(err).WithField("blockSlot", blk.Block.Slot).Debug("Ignored block")
		return pubsub.ValidationIgnore
	}

	// Add metrics for block arrival time subtracts slot start time.
	genesisTime := uint64(s.chain.GenesisTime().Unix())
	if err := captureArrivalTimeMetric(genesisTime, blk.Block.Slot); err != nil {
		log.WithError(err).WithField("blockSlot", blk.Block.Slot).Debug("Ignored block")
		return pubsub.ValidationIgnore
	}

	startSlot, err := helpers.StartSlot(s.chain.FinalizedCheckpt().Epoch)
	if err != nil {
		log.WithError(err).WithField("blockSlot", blk.Block.Slot).Debug("Ignored block")
		return pubsub.ValidationIgnore
	}
	if startSlot >= blk.Block.Slot {
		e := fmt.Errorf("finalized slot %d greater or equal to block slot %d", startSlot, blk.Block.Slot)
		log.WithError(e).WithField("blockSlot", blk.Block.Slot).Debug("Ignored block")
		return pubsub.ValidationIgnore
	}

	// Handle block when the parent is unknown.
	if !s.db.HasBlock(ctx, bytesutil.ToBytes32(blk.Block.ParentRoot)) {
		s.pendingQueueLock.Lock()
		if err := s.insertBlockToPendingQueue(blk.Block.Slot, blk, blockRoot); err != nil {
			s.pendingQueueLock.Unlock()
			log.WithError(err).WithField("blockSlot", blk.Block.Slot).Debug("Ignored block")
			return pubsub.ValidationIgnore
		}
		s.pendingQueueLock.Unlock()
		log.WithError(errors.New("unknown parent")).WithField("blockSlot", blk.Block.Slot).Debug("Ignored block")
		return pubsub.ValidationIgnore
	}

	if err := s.validateBeaconBlock(ctx, blk, blockRoot); err != nil {
		log.WithError(err).WithField("blockSlot", blk.Block.Slot).Warn("Rejected block")
		return pubsub.ValidationReject
	}
	// Record attribute of valid block.
	span.AddAttributes(trace.Int64Attribute("slotInEpoch", int64(blk.Block.Slot%params.BeaconConfig().SlotsPerEpoch)))
	msg.ValidatorData = blk // Used in downstream subscriber

	// Log the arrival time of the accepted block
	startTime, err := helpers.SlotToTime(genesisTime, blk.Block.Slot)
	if err != nil {
		log.WithError(err).WithField("blockSlot", blk.Block.Slot).Debug("Couldn't get slot start time")
		return pubsub.ValidationIgnore
	}
	log.WithFields(logrus.Fields{
		"blockSlot":          blk.Block.Slot,
		"sinceSlotStartTime": receivedTime.Sub(startTime),
	}).Debug("Received block")
	return pubsub.ValidationAccept
}

func (s *Service) validateBeaconBlock(ctx context.Context, blk *ethpb.SignedBeaconBlock, blockRoot [32]byte) error {
	ctx, span := trace.StartSpan(ctx, "sync.validateBeaconBlock")
	defer span.End()

	if err := s.chain.VerifyBlkDescendant(ctx, bytesutil.ToBytes32(blk.Block.ParentRoot)); err != nil {
		s.setBadBlock(ctx, blockRoot)
		return err
	}

	hasStateSummaryDB := s.db.HasStateSummary(ctx, bytesutil.ToBytes32(blk.Block.ParentRoot))
	if !hasStateSummaryDB {
		_, err := s.stateGen.RecoverStateSummary(ctx, bytesutil.ToBytes32(blk.Block.ParentRoot))
		if err != nil {
			return err
		}
	}
	parentState, err := s.stateGen.StateByRoot(ctx, bytesutil.ToBytes32(blk.Block.ParentRoot))
	if err != nil {
		return err
	}

	if err := blocks.VerifyBlockSignature(parentState, blk); err != nil {
		s.setBadBlock(ctx, blockRoot)
		return err
	}
	// There is an epoch lookahead for validator proposals
	// for the next epoch from the start of our current epoch. We
	// use the randao mix at the end of the previous epoch as the seed
	// to determine proposals.
	// Seed for Next Epoch => Derived From Randao Mix at the end of the Previous Epoch.
	// Which is why we simply set the slot over here.
	nextEpoch := helpers.NextEpoch(parentState)
	expectedEpoch := helpers.SlotToEpoch(blk.Block.Slot)
	if expectedEpoch <= nextEpoch {
		err = parentState.SetSlot(blk.Block.Slot)
		if err != nil {
			return err
		}
	} else {
		// In the event the block is more than an epoch ahead from its
		// parent state, we have to advance the state forward.
		parentState, err = state.ProcessSlots(ctx, parentState, blk.Block.Slot)
		if err != nil {
			return err
		}
	}
	idx, err := helpers.BeaconProposerIndex(parentState)
	if err != nil {
		return err
	}
	if blk.Block.ProposerIndex != idx {
		s.setBadBlock(ctx, blockRoot)
		return errors.New("incorrect proposer index")
	}

	return nil
}

// Returns true if the block is not the first block proposed for the proposer for the slot.
func (s *Service) hasSeenBlockIndexSlot(slot types.Slot, proposerIdx types.ValidatorIndex) bool {
	s.seenBlockLock.RLock()
	defer s.seenBlockLock.RUnlock()
	b := append(bytesutil.Bytes32(uint64(slot)), bytesutil.Bytes32(uint64(proposerIdx))...)
	_, seen := s.seenBlockCache.Get(string(b))
	return seen
}

// Set block proposer index and slot as seen for incoming blocks.
func (s *Service) setSeenBlockIndexSlot(slot types.Slot, proposerIdx types.ValidatorIndex) {
	s.seenBlockLock.Lock()
	defer s.seenBlockLock.Unlock()
	b := append(bytesutil.Bytes32(uint64(slot)), bytesutil.Bytes32(uint64(proposerIdx))...)
	s.seenBlockCache.Add(string(b), true)
}

// Returns true if the block is marked as a bad block.
func (s *Service) hasBadBlock(root [32]byte) bool {
	s.badBlockLock.RLock()
	defer s.badBlockLock.RUnlock()
	_, seen := s.badBlockCache.Get(string(root[:]))
	return seen
}

// Set bad block in the cache.
func (s *Service) setBadBlock(ctx context.Context, root [32]byte) {
	s.badBlockLock.Lock()
	defer s.badBlockLock.Unlock()
	if ctx.Err() != nil { // Do not mark block as bad if it was due to context error.
		return
	}
	s.badBlockCache.Add(string(root[:]), true)
}

// This captures metrics for block arrival time by subtracts slot start time.
func captureArrivalTimeMetric(genesisTime uint64, currentSlot types.Slot) error {
	startTime, err := helpers.SlotToTime(genesisTime, currentSlot)
	if err != nil {
		return err
	}
	ms := timeutils.Now().Sub(startTime) / time.Millisecond
	arrivalBlockPropagationHistogram.Observe(float64(ms))

	return nil
}
