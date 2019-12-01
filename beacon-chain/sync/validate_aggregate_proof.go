package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// validateAggregateAndProof validates the aggregated signature and the selection proof is valid before forwarding to the
// network and downstream services.
func (r *RegularSync) validateAggregateAndProof(ctx context.Context, msg proto.Message, p p2p.Broadcaster, fromSelf bool) (bool, error) {
	// To process the following it requires the recent blocks to be present in the database, so we'll skip
	// validating or processing aggregated attestations until fully synced.
	if r.initialSync.Syncing() {
		return false, nil
	}

	m := msg.(*pb.AggregateAndProof)
	attSlot := m.Aggregate.Data.Slot

	// Verify aggregate attestation has not already been seen via aggregate gossip, within a block, or through the creation locally.
	// TODO(3835): Blocked by operation pool redesign

	// Verify the block being voted for passes validation. The block should have passed validation if it's in the DB.
	if !r.db.HasBlock(ctx, bytesutil.ToBytes32(m.Aggregate.Data.BeaconBlockRoot)) {
		return false, errPointsToBlockNotInDatabase
	}

	s, err := r.chain.HeadState(ctx)
	if err != nil {
		return false, err
	}
	s, err = state.ProcessSlots(ctx, s, attSlot)
	if err != nil {
		return false, err
	}

	// Verify attestation slot is within the last ATTESTATION_PROPAGATION_SLOT_RANGE slots.
	currentSlot := (uint64(time.Now().Unix()) - s.GenesisTime) / params.BeaconConfig().SecondsPerSlot
	if attSlot > currentSlot || currentSlot > attSlot+params.BeaconConfig().AttestationPropagationSlotRange {
		return false, fmt.Errorf("attestation slot out of range %d < %d < %d",
			attSlot, currentSlot, attSlot+params.BeaconConfig().AttestationPropagationSlotRange)
	}

	// Verify validator index is within the aggregate's committee.
	if err := verifyIndexInCommittee(ctx, s, m.Aggregate, m.AggregatorIndex); err != nil {
		return false, errors.Wrapf(err, "Could not verify index in committee")
	}

	// Verify selection proof reflects to the right validator and signature is valid.
	if err := verifySelection(s, m.Aggregate.Data, m.AggregatorIndex, m.SelectionProof); err != nil {
		return false, errors.Wrapf(err, "Could not verify selection for validator %d", m.AggregatorIndex)
	}

	// Verify aggregated attestation has a valid signature.
	if err := blocks.VerifyAttestation(ctx, s, m.Aggregate); err != nil {
		return false, err
	}

	return true, nil
}

// This verifies selection proof by validating it's from the correct validator index of the slot and selection
// proof is a valid signature.
func verifySelection(s *pb.BeaconState, data *ethpb.AttestationData, validatorIndex uint64, proof []byte) error {
	slotSig, err := bls.SignatureFromBytes(proof)
	if err != nil {
		return err
	}
	aggregator, err := helpers.IsAggregator(s, data.Slot, data.CommitteeIndex, slotSig)
	if err != nil {
		return err
	}
	if !aggregator {
		return fmt.Errorf("validator is not an aggregator for slot %d", data.Slot)
	}

	domain := helpers.Domain(s.Fork, helpers.SlotToEpoch(data.Slot), params.BeaconConfig().DomainBeaconAttester)
	slotMsg, err := ssz.HashTreeRoot(data.Slot)
	if err != nil {
		return err
	}
	pubKey, err := bls.PublicKeyFromBytes(s.Validators[validatorIndex].PublicKey)
	if err != nil {
		return err
	}
	if !slotSig.Verify(slotMsg[:], pubKey, domain) {
		return err
	}

	return nil
}

// This verifies the aggregator's index in state is within the attesting indices of the attestation.
func verifyIndexInCommittee(ctx context.Context, s *pb.BeaconState, a *ethpb.Attestation, validatorIndex uint64) error {
	currentState, err := state.ProcessSlots(ctx, s, a.Data.Slot)
	if err != nil {
		return err
	}
	attestingIndices, err := helpers.AttestingIndices(currentState, a.Data, a.AggregationBits)
	var withinCommittee bool
	for _, i := range attestingIndices {
		if validatorIndex == i {
			withinCommittee = true
			break
		}
	}
	if !withinCommittee {
		return fmt.Errorf("validator index %d is not within the committee %v",
			validatorIndex, attestingIndices)
	}
	return nil
}
