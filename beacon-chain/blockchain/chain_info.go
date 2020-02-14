package blockchain

import (
	"bytes"
	"context"
	"time"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/epoch/precompute"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// ChainInfoFetcher defines a common interface for methods in blockchain service which
// directly retrieves chain info related data.
type ChainInfoFetcher interface {
	HeadFetcher
	FinalizationFetcher
}

// TimeFetcher retrieves the Eth2 data that's related to time.
type TimeFetcher interface {
	GenesisTime() time.Time
	CurrentSlot() uint64
}

// HeadFetcher defines a common interface for methods in blockchain service which
// directly retrieves head related data.
type HeadFetcher interface {
	HeadSlot() uint64
	HeadRoot(ctx context.Context) ([]byte, error)
	HeadBlock() *ethpb.SignedBeaconBlock
	HeadState(ctx context.Context) (*state.BeaconState, error)
	HeadValidatorsIndices(epoch uint64) ([]uint64, error)
	HeadSeed(epoch uint64) ([32]byte, error)
}

// ForkFetcher retrieves the current fork information of the Ethereum beacon chain.
type ForkFetcher interface {
	CurrentFork() *pb.Fork
}

// FinalizationFetcher defines a common interface for methods in blockchain service which
// directly retrieves finalization and justification related data.
type FinalizationFetcher interface {
	FinalizedCheckpt() *ethpb.Checkpoint
	CurrentJustifiedCheckpt() *ethpb.Checkpoint
	PreviousJustifiedCheckpt() *ethpb.Checkpoint
}

// ParticipationFetcher defines a common interface for methods in blockchain service which
// directly retrieves validator participation related data.
type ParticipationFetcher interface {
	Participation(epoch uint64) *precompute.Balance
}

// FinalizedCheckpt returns the latest finalized checkpoint from head state.
func (s *Service) FinalizedCheckpt() *ethpb.Checkpoint {
	if s.finalizedCheckpt == nil {
		return &ethpb.Checkpoint{Root: params.BeaconConfig().ZeroHash[:]}
	}

	cpt := state.CopyCheckpoint(s.finalizedCheckpt)
	// If head state exists but there hasn't been a finalized check point,
	// the check point's root should refer to genesis block root.
	if bytes.Equal(cpt.Root, params.BeaconConfig().ZeroHash[:]) {
		return &ethpb.Checkpoint{Root: s.genesisRoot[:]}
	}

	return cpt
}

// CurrentJustifiedCheckpt returns the current justified checkpoint from head state.
func (s *Service) CurrentJustifiedCheckpt() *ethpb.Checkpoint {
	if s.justifiedCheckpt == nil {
		return &ethpb.Checkpoint{Root: params.BeaconConfig().ZeroHash[:]}
	}

	cpt := state.CopyCheckpoint(s.justifiedCheckpt)
	// If head state exists but there hasn't been a justified check point,
	// the check point root should refer to genesis block root.
	if bytes.Equal(cpt.Root, params.BeaconConfig().ZeroHash[:]) {
		return &ethpb.Checkpoint{Root: s.genesisRoot[:]}
	}

	return cpt
}

// PreviousJustifiedCheckpt returns the previous justified checkpoint from head state.
func (s *Service) PreviousJustifiedCheckpt() *ethpb.Checkpoint {
	if s.prevFinalizedCheckpt == nil {
		return &ethpb.Checkpoint{Root: params.BeaconConfig().ZeroHash[:]}
	}

	cpt := state.CopyCheckpoint(s.prevFinalizedCheckpt)
	// If head state exists but there hasn't been a justified check point,
	// the check point root should refer to genesis block root.
	if bytes.Equal(cpt.Root, params.BeaconConfig().ZeroHash[:]) {
		return &ethpb.Checkpoint{Root: s.genesisRoot[:]}
	}

	return cpt
}

// HeadSlot returns the slot of the head of the chain.
func (s *Service) HeadSlot() uint64 {
	if !s.hasHeadState() {
		return 0
	}

	return s.headSlot()
}

// HeadRoot returns the root of the head of the chain.
func (s *Service) HeadRoot(ctx context.Context) ([]byte, error) {
	if s.headRoot() != params.BeaconConfig().ZeroHash {
		r := s.headRoot()
		return r[:], nil
	}

	b, err := s.beaconDB.HeadBlock(ctx)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return params.BeaconConfig().ZeroHash[:], nil
	}

	r, err := ssz.HashTreeRoot(b.Block)
	if err != nil {
		return nil, err
	}

	return r[:], nil
}

// HeadBlock returns the head block of the chain.
func (s *Service) HeadBlock() *ethpb.SignedBeaconBlock {
	return s.headBlock()
}

// HeadState returns the head state of the chain.
// If the head state is nil from service struct,
// it will attempt to get from DB and error if nil again.
func (s *Service) HeadState(ctx context.Context) (*state.BeaconState, error) {
	if s.hasHeadState() {
		return s.headState(), nil
	}

	headState, err := s.beaconDB.HeadState(ctx)
	if err != nil {
		return nil, err
	}
	return headState, nil
}

// HeadValidatorsIndices returns a list of active validator indices from the head view of a given epoch.
func (s *Service) HeadValidatorsIndices(epoch uint64) ([]uint64, error) {
	if !s.hasHeadState() {
		return []uint64{}, nil
	}
	return helpers.ActiveValidatorIndices(s.headState(), epoch)
}

// HeadSeed returns the seed from the head view of a given epoch.
func (s *Service) HeadSeed(epoch uint64) ([32]byte, error) {
	if !s.hasHeadState() {
		return [32]byte{}, nil
	}

	return helpers.Seed(s.headState(), epoch, params.BeaconConfig().DomainBeaconAttester)
}

// GenesisTime returns the genesis time of beacon chain.
func (s *Service) GenesisTime() time.Time {
	return s.genesisTime
}

// CurrentFork retrieves the latest fork information of the beacon chain.
func (s *Service) CurrentFork() *pb.Fork {
	if !s.hasHeadState() {
		return &pb.Fork{
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
		}
	}
	return s.headState().Fork()
}

// Participation returns the participation stats of a given epoch.
func (s *Service) Participation(epoch uint64) *precompute.Balance {
	s.epochParticipationLock.RLock()
	defer s.epochParticipationLock.RUnlock()

	return s.epochParticipation[epoch]
}
