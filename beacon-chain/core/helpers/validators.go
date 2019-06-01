package helpers

import (
	"fmt"

	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

var activeIndicesCache = make(map[uint64][]uint64)
var activeCountCache = make(map[uint64]uint64)

// IsActiveValidator returns the boolean value on whether the validator
// is active or not.
//
// Spec pseudocode definition:
//  def is_active_validator(validator: Validator, epoch: Epoch) -> bool:
//    """
//    Check if ``validator`` is active.
//    """
//    return validator.activation_epoch <= epoch < validator.exit_epoch
func IsActiveValidator(validator *pb.Validator, epoch uint64) bool {
	return validator.ActivationEpoch <= epoch &&
		epoch < validator.ExitEpoch
}

// IsSlashableValidator returns the boolean value on whether the validator
// is slashable or not.
//
// Spec pseudocode definition:
//  def is_slashable_validator(validator: Validator, epoch: Epoch) -> bool:
//    """
//    Check if ``validator`` is slashable.
//    """
//    return (
//        validator.activation_epoch <= epoch < validator.withdrawable_epoch and
//        validator.slashed is False
// 		)
func IsSlashableValidator(validator *pb.Validator, epoch uint64) bool {
	active := validator.ActivationEpoch <= epoch
	beforeWithdrawable := epoch < validator.WithdrawableEpoch
	return beforeWithdrawable && active && !validator.Slashed
}

// ActiveValidatorIndices filters out active validators based on validator status
// and returns their indices in a list.
//
// WARNING: This method allocates a new copy of the validator index set and is
// considered to be very memory expensive. Avoid using this unless you really
// need the active validator indices for some specific reason.
//
// Spec pseudocode definition:
//  def get_active_validator_indices(state: BeaconState, epoch: Epoch) -> List[ValidatorIndex]:
//    """
//    Get active validator indices at ``epoch``.
//    """
//    return [i for i, v in enumerate(state.validator_registry) if is_active_validator(v, epoch)]
func ActiveValidatorIndices(state *pb.BeaconState, epoch uint64) []uint64 {
	if _, ok := activeIndicesCache[epoch]; ok {
		return activeIndicesCache[epoch]
	}

	indices := make([]uint64, 0, len(state.ValidatorRegistry))
	for i, v := range state.ValidatorRegistry {
		if IsActiveValidator(v, epoch) {
			indices = append(indices, uint64(i))
		}
	}

	activeIndicesCache[epoch] = indices

	return indices
}

// ActiveValidatorCount returns the number of active validators in the state
// at the given epoch.
func ActiveValidatorCount(state *pb.BeaconState, epoch uint64) uint64 {
	if _, ok := activeCountCache[epoch]; ok {
		return activeCountCache[epoch]
	}

	var count uint64
	for _, v := range state.ValidatorRegistry {
		if IsActiveValidator(v, epoch) {
			count++
		}
	}

	activeCountCache[epoch] = count

	return count
}

// DelayedActivationExitEpoch takes in epoch number and returns when
// the validator is eligible for activation and exit.
//
// Spec pseudocode definition:
//  def get_delayed_activation_exit_epoch(epoch: Epoch) -> Epoch:
//    """
//    Return the epoch at which an activation or exit triggered in ``epoch`` takes effect.
//    """
//    return epoch + 1 + ACTIVATION_EXIT_DELAY
func DelayedActivationExitEpoch(epoch uint64) uint64 {
	return epoch + 1 + params.BeaconConfig().ActivationExitDelay
}

// ChurnLimit returns the number of validators that are allowed to
// enter and exit validator pool for an epoch.
//
// Spec pseudocode definition:
// def get_churn_limit(state: BeaconState) -> int:
//    return max(
//        MIN_PER_EPOCH_CHURN_LIMIT,
//        len(get_active_validator_indices(state, get_current_epoch(state))) // CHURN_LIMIT_QUOTIENT
//    )
func ChurnLimit(state *pb.BeaconState) uint64 {
	validatorCount := uint64(ActiveValidatorCount(state, CurrentEpoch(state)))
	if validatorCount/params.BeaconConfig().ChurnLimitQuotient > params.BeaconConfig().MinPerEpochChurnLimit {
		return validatorCount / params.BeaconConfig().ChurnLimitQuotient
	}
	return params.BeaconConfig().MinPerEpochChurnLimit
}

// BeaconProposerIndex returns proposer index of a current slot.
//
// Spec pseudocode definition:
//  def get_beacon_proposer_index(state: BeaconState) -> ValidatorIndex:
//    """
//    Return the current beacon proposer index.
//    """
//    epoch = get_current_epoch(state)
//    committees_per_slot = get_epoch_committee_count(state, epoch) // SLOTS_PER_EPOCH
//    offset = committees_per_slot * (state.slot % SLOTS_PER_EPOCH)
//    shard = (get_epoch_start_shard(state, epoch) + offset) % SHARD_COUNT
//    first_committee = get_crosslink_committee(state, epoch, shard)
//    MAX_RANDOM_BYTE = 2**8 - 1
//    seed = generate_seed(state, epoch)
//    i = 0
//    while True:
//        candidate_index = first_committee[(epoch + i) % len(first_committee)]
//        random_byte = hash(seed + int_to_bytes(i // 32, length=8))[i % 32]
//        effective_balance = state.validator_registry[candidate_index].effective_balance
//        if effective_balance * MAX_RANDOM_BYTE >= MAX_EFFECTIVE_BALANCE * random_byte:
//            return candidate_index
//        i += 1
func BeaconProposerIndex(state *pb.BeaconState) (uint64, error) {
	// Calculate the offset for slot and shard
	e := CurrentEpoch(state)
	committesPerSlot := EpochCommitteeCount(state, e) / params.BeaconConfig().SlotsPerEpoch
	offSet := committesPerSlot * (state.Slot % params.BeaconConfig().SlotsPerEpoch)

	// Calculate which shards get assigned given the epoch start shard
	// and the offset
	startShard, err := EpochStartShard(state, e)
	if err != nil {
		return 0, fmt.Errorf("could not get start shard: %v", err)
	}
	shard := (startShard + offSet) % params.BeaconConfig().ShardCount

	// Use the first committee of the given slot and shard
	// to select proposer
	firstCommittee, err := CrosslinkCommitteeAtEpoch(state, e, shard)
	if err != nil {
		return 0, fmt.Errorf("could not get first committee: %v", err)
	}
	if len(firstCommittee) == 0 {
		return 0, fmt.Errorf("empty first committee at slot %d",
			state.Slot)
	}

	// Use the generated seed to select proposer from the first committee
	maxRandomByte := uint64(1<<8 - 1)
	seed := GenerateSeed(state, e)

	// Looping through the committee to select proposer that has enough
	// effective balance.
	for i := uint64(0); ; i++ {
		candidateIndex := firstCommittee[(e+i)%uint64(len(firstCommittee))]
		b := append(seed[:], bytesutil.Bytes8(i)...)
		randomByte := hashutil.Hash(b)[i%32]
		effectiveBal := state.ValidatorRegistry[candidateIndex].EffectiveBalance
		if effectiveBal*maxRandomByte >= params.BeaconConfig().MaxEffectiveBalance*uint64(randomByte) {
			return candidateIndex, nil
		}
	}
}

// DomainVersion returns the domain version for BLS private key to sign and verify.
//
//  def get_domain(state: BeaconState,
//               domain_type: int,
//               message_epoch: int=None) -> int:
//    """
//    Return the signature domain (fork version concatenated with domain type) of a message.
//    """
//    epoch = get_current_epoch(state) if message_epoch is None else message_epoch
//    fork_version = state.fork.previous_version if epoch < state.fork.epoch else state.fork.current_version
//    return bytes_to_int(fork_version + int_to_bytes(domain_type, length=4))
func DomainVersion(state *pb.BeaconState, epoch uint64, domainType uint64) uint64 {
	if epoch == 0 {
		epoch = CurrentEpoch(state)
	}
	var forkVersion []byte
	if epoch < state.Fork.Epoch {
		forkVersion = state.Fork.PreviousVersion
	} else {
		forkVersion = state.Fork.CurrentVersion
	}
	by := []byte{}
	by = append(by, forkVersion[:4]...)
	by = append(by, bytesutil.Bytes4(domainType)...)
	return bytesutil.FromBytes8(by)
}

// RestartActiveCountCache restarts the active validator count cache from scratch.
func RestartActiveCountCache() {
	activeCountCache = make(map[uint64]uint64)
}

// RestartActiveIndicesCache restarts the active validator indices cache from scratch.
func RestartActiveIndicesCache() {
	activeIndicesCache = make(map[uint64][]uint64)
}
