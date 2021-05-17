package altair

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	iface "github.com/prysmaticlabs/prysm/beacon-chain/state/interface"
	"github.com/prysmaticlabs/prysm/shared/attestationutil"
	"github.com/prysmaticlabs/prysm/shared/mathutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"go.opencensus.io/trace"
)

type matchingTarget bool
type matchingSource bool
type matchingHead bool

// ProcessAttestations applies processing operations to a block's inner attestation
// records.
func ProcessAttestations(
	ctx context.Context,
	beaconState iface.BeaconState,
	b *ethpb.SignedBeaconBlock,
) (iface.BeaconState, error) {
	if err := helpers.VerifyNilBeaconBlock(b); err != nil {
		return nil, err
	}

	var err error
	for idx, attestation := range b.Block.Body.Attestations {
		beaconState, err = ProcessAttestation(ctx, beaconState, attestation)
		if err != nil {
			return nil, errors.Wrapf(err, "could not verify attestation at index %d in block", idx)
		}
	}
	return beaconState, nil
}

// ProcessAttestation verifies an input attestation can pass through processing using the given beacon state.
//
// Spec code:
//  def process_attestation(state: BeaconState, attestation: Attestation) -> None:
//    data = attestation.data
//    assert data.target.epoch in (get_previous_epoch(state), get_current_epoch(state))
//    assert data.target.epoch == compute_epoch_at_slot(data.slot)
//    assert data.slot + MIN_ATTESTATION_INCLUSION_DELAY <= state.slot <= data.slot + SLOTS_PER_EPOCH
//    assert data.index < get_committee_count_per_slot(state, data.target.epoch)
//
//    committee = get_beacon_committee(state, data.slot, data.index)
//    assert len(attestation.aggregation_bits) == len(committee)
//
//    if data.target.epoch == get_current_epoch(state):
//        epoch_participation = state.current_epoch_participation
//        justified_checkpoint = state.current_justified_checkpoint
//    else:
//        epoch_participation = state.previous_epoch_participation
//        justified_checkpoint = state.previous_justified_checkpoint
//
//    # Matching roots
//    is_matching_head = data.beacon_block_root == get_block_root_at_slot(state, data.slot)
//    is_matching_source = data.source == justified_checkpoint
//    is_matching_target = data.target.root == get_block_root(state, data.target.epoch)
//    assert is_matching_source
//
//    # Verify signature
//    assert is_valid_indexed_attestation(state, get_indexed_attestation(state, attestation))
//
//    # Participation flag indices
//    participation_flag_indices = []
//    if is_matching_head and is_matching_target and state.slot == data.slot + MIN_ATTESTATION_INCLUSION_DELAY:
//        participation_flag_indices.append(TIMELY_HEAD_FLAG_INDEX)
//    if is_matching_source and state.slot <= data.slot + integer_squareroot(SLOTS_PER_EPOCH):
//        participation_flag_indices.append(TIMELY_SOURCE_FLAG_INDEX)
//    if is_matching_target and state.slot <= data.slot + SLOTS_PER_EPOCH:
//        participation_flag_indices.append(TIMELY_TARGET_FLAG_INDEX)
//
//    # Update epoch participation flags
//    proposer_reward_numerator = 0
//    for index in get_attesting_indices(state, data, attestation.aggregation_bits):
//        for flag_index, weight in get_flag_indices_and_weights():
//            if flag_index in participation_flag_indices and not has_flag(epoch_participation[index], flag_index):
//                epoch_participation[index] = add_flag(epoch_participation[index], flag_index)
//                proposer_reward_numerator += get_base_reward(state, index) * weight
//
//    # Reward proposer
//    proposer_reward_denominator = (WEIGHT_DENOMINATOR - PROPOSER_WEIGHT) * WEIGHT_DENOMINATOR // PROPOSER_WEIGHT
//    proposer_reward = Gwei(proposer_reward_numerator // proposer_reward_denominator)
//    increase_balance(state, get_beacon_proposer_index(state), proposer_reward)
func ProcessAttestation(
	ctx context.Context,
	beaconState iface.BeaconStateAltair,
	att *ethpb.Attestation,
) (iface.BeaconStateAltair, error) {
	beaconState, err := ProcessAttestationNoVerifySignature(ctx, beaconState, att)
	if err != nil {
		return nil, err
	}
	return beaconState, blocks.VerifyAttestationSignature(ctx, beaconState, att)
}

// ProcessAttestationsNoVerifySignature applies processing operations to a block's inner attestation
// records. The only difference would be that the attestation signature would not be verified.
func ProcessAttestationsNoVerifySignature(
	ctx context.Context,
	beaconState iface.BeaconState,
	b *ethpb.SignedBeaconBlockAltair,
) (iface.BeaconState, error) {
	if err := VerifyNilBeaconBlock(b); err != nil {
		return nil, err
	}
	body := b.Block.Body
	var err error
	for idx, attestation := range body.Attestations {
		beaconState, err = ProcessAttestationNoVerifySignature(ctx, beaconState, attestation)
		if err != nil {
			return nil, errors.Wrapf(err, "could not verify attestation at index %d in block", idx)
		}
	}
	return beaconState, nil
}

// ProcessAttestationNoVerifySignature processes the attestation without verifying the attestation signature. This
// method is used to validate attestations whose signatures have already been verified or will be verified later.
func ProcessAttestationNoVerifySignature(
	ctx context.Context,
	beaconState iface.BeaconStateAltair,
	att *ethpb.Attestation,
) (iface.BeaconStateAltair, error) {
	ctx, span := trace.StartSpan(ctx, "core.ProcessAttestationNoVerifySignature")
	defer span.End()

	if err := blocks.VerifyAttestationNoVerifySignature(ctx, beaconState, att); err != nil {
		return nil, err
	}

	currEpoch := helpers.CurrentEpoch(beaconState)
	data := att.Data
	var justifiedCheckpt *ethpb.Checkpoint
	var epochParticipation []byte
	var err error
	if data.Target.Epoch == currEpoch {
		justifiedCheckpt = beaconState.CurrentJustifiedCheckpoint()
		epochParticipation, err = beaconState.CurrentEpochParticipation()
		if err != nil {
			return nil, err
		}
	} else {
		justifiedCheckpt = beaconState.PreviousJustifiedCheckpoint()
		epochParticipation, err = beaconState.PreviousEpochParticipation()
		if err != nil {
			return nil, err
		}
	}

	// Get matching participation flags.
	matchingSource, matchingTarget, matchingHead, err := matchingStatus(beaconState, data, justifiedCheckpt)
	if err != nil {
		return nil, err
	}
	if !matchingSource {
		return nil, errors.New("source epoch does not match")
	}

	// Process matched participation flags.
	participatedFlags := make(map[uint8]bool)
	headFlagIndex := params.BeaconConfig().TimelyHeadFlagIndex
	sourceFlagIndex := params.BeaconConfig().TimelySourceFlagIndex
	targetFlagIndex := params.BeaconConfig().TimelyTargetFlagIndex
	if matchingSource && beaconState.Slot() <= data.Slot.Add(mathutil.IntegerSquareRoot(uint64(params.BeaconConfig().SlotsPerEpoch))) {
		participatedFlags[sourceFlagIndex] = true
	}
	matchingHeadTarget := bool(matchingHead) && bool(matchingTarget)
	if matchingHeadTarget && beaconState.Slot() == data.Slot+params.BeaconConfig().MinAttestationInclusionDelay {
		participatedFlags[headFlagIndex] = true
	}
	if matchingTarget && beaconState.Slot() <= data.Slot+params.BeaconConfig().SlotsPerEpoch {
		participatedFlags[targetFlagIndex] = true
	}

	committee, err := helpers.BeaconCommitteeFromState(beaconState, att.Data.Slot, att.Data.CommitteeIndex)
	if err != nil {
		return nil, err
	}
	indices, err := attestationutil.AttestingIndices(att.AggregationBits, committee)
	if err != nil {
		return nil, err
	}
	proposerRewardNumerator := uint64(0)
	for _, index := range indices {
		br, err := BaseReward(beaconState, types.ValidatorIndex(index))
		if err != nil {
			return nil, err
		}
		if participatedFlags[sourceFlagIndex] && !HasValidatorFlag(epochParticipation[index], sourceFlagIndex) {
			epochParticipation[index] = AddValidatorFlag(epochParticipation[index], sourceFlagIndex)
			proposerRewardNumerator += br * params.BeaconConfig().TimelySourceWeight
		}
		if participatedFlags[targetFlagIndex] && !HasValidatorFlag(epochParticipation[index], targetFlagIndex) {
			epochParticipation[index] = AddValidatorFlag(epochParticipation[index], targetFlagIndex)
			proposerRewardNumerator += br * params.BeaconConfig().TimelyTargetWeight
		}
		if participatedFlags[headFlagIndex] && !HasValidatorFlag(epochParticipation[index], headFlagIndex) {
			epochParticipation[index] = AddValidatorFlag(epochParticipation[index], headFlagIndex)
			proposerRewardNumerator += br * params.BeaconConfig().TimelyHeadWeight
		}
	}

	if data.Target.Epoch == currEpoch {
		if err := beaconState.SetCurrentParticipationBits(epochParticipation); err != nil {
			return nil, err
		}
	} else {
		if err := beaconState.SetPreviousParticipationBits(epochParticipation); err != nil {
			return nil, err
		}
	}

	// Reward proposer.
	if err := rewardProposer(beaconState, proposerRewardNumerator); err != nil {
		return nil, err
	}
	return beaconState, nil
}

// This returns the matching statues for attestation data's source target and head.
func matchingStatus(beaconState iface.BeaconState, data *ethpb.AttestationData, cp *ethpb.Checkpoint) (s matchingSource, t matchingTarget, h matchingHead, err error) {
	s = matchingSource(attestationutil.CheckPointIsEqual(data.Source, cp))

	r, err := helpers.BlockRoot(beaconState, data.Target.Epoch)
	if err != nil {
		return false, false, false, err
	}
	t = matchingTarget(bytes.Equal(r, data.Target.Root))

	r, err = helpers.BlockRootAtSlot(beaconState, data.Slot)
	if err != nil {
		return false, false, false, err
	}
	h = matchingHead(bytes.Equal(r, data.BeaconBlockRoot))
	return
}

// This rewards proposer by increasing proposer's balance with input reward numerator and calculated reward denominator.
func rewardProposer(beaconState iface.BeaconState, proposerRewardNumerator uint64) error {
	proposerRewardDenominator := (params.BeaconConfig().WeightDenominator - params.BeaconConfig().ProposerWeight) * params.BeaconConfig().WeightDenominator / params.BeaconConfig().ProposerWeight
	proposerReward := proposerRewardNumerator / proposerRewardDenominator
	i, err := helpers.BeaconProposerIndex(beaconState)
	if err != nil {
		return err
	}

	return helpers.IncreaseBalance(beaconState, i, proposerReward)
}

// HasValidatorFlag returns true if the flag at position has set.
func HasValidatorFlag(flag, flagPosition uint8) bool {
	return ((flag >> flagPosition) & 1) == 1
}

// AddValidatorFlag adds new validator flag to existing one.
func AddValidatorFlag(flag, flagPosition uint8) uint8 {
	return flag | (1 << flagPosition)
}
