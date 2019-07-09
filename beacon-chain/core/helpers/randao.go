package helpers

import (
	"encoding/binary"
	"fmt"

	"github.com/prysmaticlabs/prysm/beacon-chain/cache"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

var currentEpochSeed = cache.NewSeedCache()

// Seed returns the randao seed used for shuffling of a given epoch.
//
// Spec pseudocode definition:
//  def get_seed(state: BeaconState, epoch: Epoch) -> Hash:
//    """
//    Return the seed at ``epoch``.
//    """
//    mix = get_randao_mix(state, Epoch(epoch + EPOCHS_PER_HISTORICAL_VECTOR - MIN_SEED_LOOKAHEAD))  # Avoid underflow
//    active_index_root = state.active_index_roots[epoch % EPOCHS_PER_HISTORICAL_VECTOR]
//    return hash(mix + active_index_root + int_to_bytes(epoch, length=32))
func Seed(state *pb.BeaconState, epoch uint64) ([32]byte, error) {
	seed, err := currentEpochSeed.SeedInEpoch(epoch)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not retrieve total balance from cache: %v", err)
	}
	if seed != nil {
		return bytesutil.ToBytes32(seed), nil
	}

	lookAheadEpoch := epoch + params.BeaconConfig().EpochsPerHistoricalVector -
		params.BeaconConfig().MinSeedLookahead

	randaoMix := RandaoMix(state, lookAheadEpoch)

	indexRoot := ActiveIndexRoot(state, epoch)

	th := append(randaoMix, indexRoot...)
	th = append(th, bytesutil.Bytes32(epoch)...)

	seed32 := hashutil.Hash(th)

	if err := currentEpochSeed.AddSeed(&cache.SeedByEpoch{
		Epoch: epoch,
		Seed:  seed32[:],
	}); err != nil {
		return [32]byte{}, fmt.Errorf("could not save active balance for cache: %v", err)
	}

	return seed32, nil
}

// ActiveIndexRoot returns the index root of a given epoch.
//
// Spec pseudocode definition:
//   def get_active_index_root(state: BeaconState,
//                          epoch: Epoch) -> Bytes32:
//    """
//    Return the index root at a recent ``epoch``.
//    ``epoch`` expected to be between
//    (current_epoch - LATEST_ACTIVE_INDEX_ROOTS_LENGTH + ACTIVATION_EXIT_DELAY, current_epoch + ACTIVATION_EXIT_DELAY].
//    """
//    return state.latest_active_index_roots[epoch % LATEST_ACTIVE_INDEX_ROOTS_LENGTH]
func ActiveIndexRoot(state *pb.BeaconState, epoch uint64) []byte {
	newRootLength := len(state.ActiveIndexRoots[epoch%params.BeaconConfig().EpochsPerHistoricalVector])
	newRoot := make([]byte, newRootLength)
	copy(newRoot, state.ActiveIndexRoots[epoch%params.BeaconConfig().EpochsPerHistoricalVector])
	return newRoot
}

// RandaoMix returns the randao mix (xor'ed seed)
// of a given slot. It is used to shuffle validators.
//
// Spec pseudocode definition:
//   def get_randao_mix(state: BeaconState, epoch: Epoch) -> Hash:
//    """
//    Return the randao mix at a recent ``epoch``.
//    """
//    return state.randao_mixes[epoch % EPOCHS_PER_HISTORICAL_VECTOR]
func RandaoMix(state *pb.BeaconState, epoch uint64) []byte {
	newMixLength := len(state.RandaoMixes[epoch%params.BeaconConfig().EpochsPerHistoricalVector])
	newMix := make([]byte, newMixLength)
	copy(newMix, state.RandaoMixes[epoch%params.BeaconConfig().EpochsPerHistoricalVector])
	return newMix
}

// CreateRandaoReveal generates a epoch signature using the beacon proposer priv key.
func CreateRandaoReveal(beaconState *pb.BeaconState, epoch uint64, privKeys []*bls.SecretKey) ([]byte, error) {
	// We fetch the proposer's index as that is whom the RANDAO will be verified against.
	proposerIdx, err := BeaconProposerIndex(beaconState)
	if err != nil {
		return []byte{}, fmt.Errorf("could not get beacon proposer index: %v", err)
	}
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint64(buf, epoch)
	domain := Domain(beaconState, epoch, params.BeaconConfig().DomainRandao)
	// We make the previous validator's index sign the message instead of the proposer.
	epochSignature := privKeys[proposerIdx].Sign(buf, domain)
	return epochSignature.Marshal(), nil
}
