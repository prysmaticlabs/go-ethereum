package altair

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	b "github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	iface "github.com/prysmaticlabs/prysm/beacon-chain/state/interface"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"go.opencensus.io/trace"
)

// ExecuteStateTransitionNoVerifyAnySig defines the procedure for a state transition function.
// This does not validate any BLS signatures of attestations, block proposer signature, randao signature,
// it is used for performing a state transition as quickly as possible. This function also returns a signature
// set of all signatures not verified, so that they can be stored and verified later.
func ExecuteStateTransitionNoVerifyAnySig(
	ctx context.Context,
	state iface.BeaconStateAltair,
	signed *ethpb.SignedBeaconBlockAltair,
) (*bls.SignatureSet, iface.BeaconStateAltair, error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}
	if signed == nil || signed.Block == nil {
		return nil, nil, errors.New("nil block")
	}

	ctx, span := trace.StartSpan(ctx, "altair.ExecuteStateTransitionNoVerifyAttSigs")
	defer span.End()
	var err error

	// Next slot cache is not enabled here.
	state, err = ProcessSlots(ctx, state, signed.Block.Slot)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not process slot")
	}

	// Execute per block transition.
	set, state, err := ProcessBlockNoVerifyAnySig(ctx, state, signed)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not process block")
	}

	// State root validation.
	postStateRoot, err := state.HashTreeRoot(ctx)
	if err != nil {
		return nil, nil, err
	}
	if !bytes.Equal(postStateRoot[:], signed.Block.StateRoot) {
		return nil, nil, fmt.Errorf("could not validate state root, wanted: %#x, received: %#x",
			postStateRoot[:], signed.Block.StateRoot)
	}

	return set, state, nil
}

// ProcessBlockNoVerifyAnySig creates a new, modified beacon state by applying block operation
// transformations as defined in the Ethereum Serenity specification. It does not validate
// any block signature except for deposit and slashing signatures. It also returns the relevant
// signature set from all the respective methods.
//
// Spec pseudocode definition:
//
//  def process_block(state: BeaconState, block: BeaconBlock) -> None:
//    process_block_header(state, block)
//    process_randao(state, block.body)
//    process_eth1_data(state, block.body)
//    process_operations(state, block.body)  # [Modified in Altair]
//    process_sync_committee(state, block.body.sync_aggregate)  # [New in Altair]
func ProcessBlockNoVerifyAnySig(
	ctx context.Context,
	state iface.BeaconState,
	signed *ethpb.SignedBeaconBlockAltair,
) (*bls.SignatureSet, iface.BeaconState, error) {
	ctx, span := trace.StartSpan(ctx, "altair.ProcessBlockNoVerifyAnySig")
	defer span.End()

	// Verify block is not nil.
	if err := VerifyNilBeaconBlock(signed); err != nil {
		return nil, nil, err
	}

	blk := signed.Block
	body := blk.Body
	state, err := ProcessBlockForStateRoot(ctx, state, signed)
	if err != nil {
		return nil, nil, err
	}

	state, err = ProcessSyncCommittee(state, signed.Block.Body.SyncAggregate)
	if err != nil {
		return nil, nil, err
	}

	bSet, err := b.BlockSignatureSet(state, blk.ProposerIndex, signed.Signature, blk.HashTreeRoot)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not retrieve block signature set")
	}
	rSet, err := b.RandaoSignatureSet(state, body.RandaoReveal)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not retrieve randao signature set")
	}
	aSet, err := b.AttestationSignatureSet(ctx, state, signed.Block.Body.Attestations)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not retrieve attestation signature set")
	}

	// Merge beacon block, randao and attestations signatures into a set.
	set := bls.NewSet()
	set.Join(bSet).Join(rSet).Join(aSet)

	return set, state, nil
}

// ProcessOperationsNoVerifyAttsSigs processes the operations in the beacon block and updates beacon state
// with the operations in block. It does not verify attestation signatures.
//
// WARNING: This method does not verify attestation signatures.
// This is used to perform the block operations as fast as possible.
//
// Spec pseudocode definition:
//
//  def process_operations(state: BeaconState, body: BeaconBlockBody) -> None:
//    # Verify that outstanding deposits are processed up to the maximum number of deposits
//    assert len(body.deposits) == min(MAX_DEPOSITS, state.eth1_data.deposit_count - state.eth1_deposit_index)
//
//    def for_ops(operations: Sequence[Any], fn: Callable[[BeaconState, Any], None]) -> None:
//        for operation in operations:
//            fn(state, operation)
//
//    for_ops(body.proposer_slashings, process_proposer_slashing)
//    for_ops(body.attester_slashings, process_attester_slashing)
//    for_ops(body.attestations, process_attestation)
//    for_ops(body.deposits, process_deposit)
//    for_ops(body.voluntary_exits, process_voluntary_exit)
func ProcessOperationsNoVerifyAttsSigs(
	ctx context.Context,
	state iface.BeaconState,
	signed *ethpb.SignedBeaconBlockAltair) (iface.BeaconState, error) {
	ctx, span := trace.StartSpan(ctx, "altair.ProcessOperationsNoVerifyAttsSigs")
	defer span.End()

	if err := VerifyNilBeaconBlock(signed); err != nil {
		return nil, err
	}

	blk := signed.Block
	body := blk.Body

	if _, err := VerifyOperationLengths(state, signed); err != nil {
		return nil, errors.Wrap(err, "could not verify operation lengths")
	}

	// Modified in Altair.
	state, err := b.ProcessProposerSlashings(ctx, state, body.ProposerSlashings, SlashValidator)
	if err != nil {
		return nil, errors.Wrap(err, "could not process block proposer slashings")
	}

	// Modified in Altair.
	state, err = b.ProcessAttesterSlashings(ctx, state, body.AttesterSlashings, SlashValidator)
	if err != nil {
		return nil, errors.Wrap(err, "could not process block attester slashings")
	}

	// Modified in Altair.
	state, err = ProcessAttestationsNoVerifySignature(ctx, state, signed)
	if err != nil {
		return nil, errors.Wrap(err, "could not process block attestations")
	}

	// Modified in Altair.
	state, err = ProcessDeposits(ctx, state, signed.Block.Body.Deposits)
	if err != nil {
		return nil, errors.Wrap(err, "could not process block validator deposits")
	}

	state, err = b.ProcessVoluntaryExits(ctx, state, body.VoluntaryExits)
	if err != nil {
		return nil, errors.Wrap(err, "could not process validator exits")
	}
	return state, nil
}
