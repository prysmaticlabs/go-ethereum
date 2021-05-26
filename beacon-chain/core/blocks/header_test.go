package blocks_test

import (
	"context"
	"io/ioutil"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	p2ptypes "github.com/prysmaticlabs/prysm/beacon-chain/p2p/types"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/interfaces"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func init() {
	logrus.SetOutput(ioutil.Discard) // Ignore "validator activated" logs
}

func TestProcessBlockHeader_ImproperBlockSlot(t *testing.T) {
	validators := make([]*ethpb.Validator, params.BeaconConfig().MinGenesisActiveValidatorCount)
	for i := 0; i < len(validators); i++ {
		validators[i] = &ethpb.Validator{
			PublicKey:             make([]byte, 32),
			WithdrawalCredentials: make([]byte, 32),
			ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
			Slashed:               true,
		}
	}

	state, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, state.SetSlot(10))
	require.NoError(t, state.SetValidators(validators))
	require.NoError(t, state.SetLatestBlockHeader(testutil.HydrateBeaconHeader(&ethpb.BeaconBlockHeader{
		Slot: 10, // Must be less than block.Slot
	})))

	latestBlockSignedRoot, err := state.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)

	currentEpoch := helpers.CurrentEpoch(state)
	priv, err := bls.RandKey()
	require.NoError(t, err)
	pID, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	block := testutil.NewBeaconBlock()
	block.Block.ProposerIndex = pID
	block.Block.Slot = 10
	block.Block.Body.RandaoReveal = bytesutil.PadTo([]byte{'A', 'B', 'C'}, 96)
	block.Block.ParentRoot = latestBlockSignedRoot[:]
	block.Signature, err = helpers.ComputeDomainAndSign(state, currentEpoch, block.Block, params.BeaconConfig().DomainBeaconProposer, priv)
	require.NoError(t, err)

	proposerIdx, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	validators[proposerIdx].Slashed = false
	validators[proposerIdx].PublicKey = priv.PublicKey().Marshal()
	err = state.UpdateValidatorAtIndex(proposerIdx, validators[proposerIdx])
	require.NoError(t, err)

	_, err = blocks.ProcessBlockHeader(context.Background(), state, interfaces.NewWrappedSignedBeaconBlock(block))
	assert.ErrorContains(t, "block.Slot 10 must be greater than state.LatestBlockHeader.Slot 10", err)
}

func TestProcessBlockHeader_WrongProposerSig(t *testing.T) {
	testutil.ResetCache()
	beaconState, privKeys := testutil.DeterministicGenesisState(t, 100)
	require.NoError(t, beaconState.SetLatestBlockHeader(testutil.HydrateBeaconHeader(&ethpb.BeaconBlockHeader{
		Slot: 9,
	})))
	require.NoError(t, beaconState.SetSlot(10))

	lbhdr, err := beaconState.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)

	proposerIdx, err := helpers.BeaconProposerIndex(beaconState)
	require.NoError(t, err)

	block := testutil.NewBeaconBlock()
	block.Block.ProposerIndex = proposerIdx
	block.Block.Slot = 10
	block.Block.Body.RandaoReveal = bytesutil.PadTo([]byte{'A', 'B', 'C'}, 96)
	block.Block.ParentRoot = lbhdr[:]
	block.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, block.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx+1])
	require.NoError(t, err)

	_, err = blocks.ProcessBlockHeader(context.Background(), beaconState, interfaces.NewWrappedSignedBeaconBlock(block))
	want := "signature did not verify"
	assert.ErrorContains(t, want, err)
}

func TestProcessBlockHeader_DifferentSlots(t *testing.T) {
	validators := make([]*ethpb.Validator, params.BeaconConfig().MinGenesisActiveValidatorCount)
	for i := 0; i < len(validators); i++ {
		validators[i] = &ethpb.Validator{
			PublicKey:             make([]byte, 32),
			WithdrawalCredentials: make([]byte, 32),
			ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
			Slashed:               true,
		}
	}

	state, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, state.SetValidators(validators))
	require.NoError(t, state.SetSlot(10))
	require.NoError(t, state.SetLatestBlockHeader(testutil.HydrateBeaconHeader(&ethpb.BeaconBlockHeader{
		Slot: 9,
	})))

	lbhsr, err := state.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)
	currentEpoch := helpers.CurrentEpoch(state)

	priv, err := bls.RandKey()
	require.NoError(t, err)
	sszBytes := p2ptypes.SSZBytes("hello")
	blockSig, err := helpers.ComputeDomainAndSign(state, currentEpoch, &sszBytes, params.BeaconConfig().DomainBeaconProposer, priv)
	require.NoError(t, err)
	validators[5896].PublicKey = priv.PublicKey().Marshal()
	block := testutil.HydrateSignedBeaconBlock(&ethpb.SignedBeaconBlock{
		Block: &ethpb.BeaconBlock{
			Slot:       1,
			ParentRoot: lbhsr[:],
		},
		Signature: blockSig,
	})

	_, err = blocks.ProcessBlockHeader(context.Background(), state, interfaces.NewWrappedSignedBeaconBlock(block))
	want := "is different than block slot"
	assert.ErrorContains(t, want, err)
}

func TestProcessBlockHeader_PreviousBlockRootNotSignedRoot(t *testing.T) {
	validators := make([]*ethpb.Validator, params.BeaconConfig().MinGenesisActiveValidatorCount)
	for i := 0; i < len(validators); i++ {
		validators[i] = &ethpb.Validator{
			PublicKey:             make([]byte, 48),
			WithdrawalCredentials: make([]byte, 32),
			ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
			Slashed:               true,
		}
	}

	state, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, state.SetValidators(validators))
	require.NoError(t, state.SetSlot(10))
	bh := state.LatestBlockHeader()
	bh.Slot = 9
	require.NoError(t, state.SetLatestBlockHeader(bh))
	currentEpoch := helpers.CurrentEpoch(state)
	priv, err := bls.RandKey()
	require.NoError(t, err)
	sszBytes := p2ptypes.SSZBytes("hello")
	blockSig, err := helpers.ComputeDomainAndSign(state, currentEpoch, &sszBytes, params.BeaconConfig().DomainBeaconProposer, priv)
	require.NoError(t, err)
	validators[5896].PublicKey = priv.PublicKey().Marshal()
	pID, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	block := testutil.NewBeaconBlock()
	block.Block.Slot = 10
	block.Block.ProposerIndex = pID
	block.Block.Body.RandaoReveal = bytesutil.PadTo([]byte{'A', 'B', 'C'}, 96)
	block.Block.ParentRoot = bytesutil.PadTo([]byte{'A'}, 32)
	block.Signature = blockSig

	_, err = blocks.ProcessBlockHeader(context.Background(), state, interfaces.NewWrappedSignedBeaconBlock(block))
	want := "does not match"
	assert.ErrorContains(t, want, err)
}

func TestProcessBlockHeader_SlashedProposer(t *testing.T) {
	validators := make([]*ethpb.Validator, params.BeaconConfig().MinGenesisActiveValidatorCount)
	for i := 0; i < len(validators); i++ {
		validators[i] = &ethpb.Validator{
			PublicKey:             make([]byte, 48),
			WithdrawalCredentials: make([]byte, 32),
			ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
			Slashed:               true,
		}
	}

	state, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, state.SetValidators(validators))
	require.NoError(t, state.SetSlot(10))
	bh := state.LatestBlockHeader()
	bh.Slot = 9
	require.NoError(t, state.SetLatestBlockHeader(bh))
	parentRoot, err := state.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)
	currentEpoch := helpers.CurrentEpoch(state)
	priv, err := bls.RandKey()
	require.NoError(t, err)
	sszBytes := p2ptypes.SSZBytes("hello")
	blockSig, err := helpers.ComputeDomainAndSign(state, currentEpoch, &sszBytes, params.BeaconConfig().DomainBeaconProposer, priv)
	require.NoError(t, err)

	validators[12683].PublicKey = priv.PublicKey().Marshal()
	pID, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	block := testutil.NewBeaconBlock()
	block.Block.Slot = 10
	block.Block.ProposerIndex = pID
	block.Block.Body.RandaoReveal = bytesutil.PadTo([]byte{'A', 'B', 'C'}, 96)
	block.Block.ParentRoot = parentRoot[:]
	block.Signature = blockSig

	_, err = blocks.ProcessBlockHeader(context.Background(), state, interfaces.NewWrappedSignedBeaconBlock(block))
	want := "was previously slashed"
	assert.ErrorContains(t, want, err)
}

func TestProcessBlockHeader_OK(t *testing.T) {
	validators := make([]*ethpb.Validator, params.BeaconConfig().MinGenesisActiveValidatorCount)
	for i := 0; i < len(validators); i++ {
		validators[i] = &ethpb.Validator{
			PublicKey:             make([]byte, 32),
			WithdrawalCredentials: make([]byte, 32),
			ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
			Slashed:               true,
		}
	}

	state, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, state.SetValidators(validators))
	require.NoError(t, state.SetSlot(10))
	require.NoError(t, state.SetLatestBlockHeader(testutil.HydrateBeaconHeader(&ethpb.BeaconBlockHeader{
		Slot: 9,
	})))

	latestBlockSignedRoot, err := state.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)

	currentEpoch := helpers.CurrentEpoch(state)
	priv, err := bls.RandKey()
	require.NoError(t, err)
	pID, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	block := testutil.NewBeaconBlock()
	block.Block.ProposerIndex = pID
	block.Block.Slot = 10
	block.Block.Body.RandaoReveal = bytesutil.PadTo([]byte{'A', 'B', 'C'}, 96)
	block.Block.ParentRoot = latestBlockSignedRoot[:]
	block.Signature, err = helpers.ComputeDomainAndSign(state, currentEpoch, block.Block, params.BeaconConfig().DomainBeaconProposer, priv)
	require.NoError(t, err)
	bodyRoot, err := block.Block.Body.HashTreeRoot()
	require.NoError(t, err, "Failed to hash block bytes got")

	proposerIdx, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	validators[proposerIdx].Slashed = false
	validators[proposerIdx].PublicKey = priv.PublicKey().Marshal()
	err = state.UpdateValidatorAtIndex(proposerIdx, validators[proposerIdx])
	require.NoError(t, err)

	newState, err := blocks.ProcessBlockHeader(context.Background(), state, interfaces.NewWrappedSignedBeaconBlock(block))
	require.NoError(t, err, "Failed to process block header got")
	var zeroHash [32]byte
	nsh := newState.LatestBlockHeader()
	expected := &ethpb.BeaconBlockHeader{
		ProposerIndex: pID,
		Slot:          block.Block.Slot,
		ParentRoot:    latestBlockSignedRoot[:],
		BodyRoot:      bodyRoot[:],
		StateRoot:     zeroHash[:],
	}
	assert.Equal(t, true, proto.Equal(nsh, expected), "Expected %v, received %v", expected, nsh)
}

func TestBlockSignatureSet_OK(t *testing.T) {
	validators := make([]*ethpb.Validator, params.BeaconConfig().MinGenesisActiveValidatorCount)
	for i := 0; i < len(validators); i++ {
		validators[i] = &ethpb.Validator{
			PublicKey:             make([]byte, 32),
			WithdrawalCredentials: make([]byte, 32),
			ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
			Slashed:               true,
		}
	}

	state, err := testutil.NewBeaconState()
	require.NoError(t, err)
	require.NoError(t, state.SetValidators(validators))
	require.NoError(t, state.SetSlot(10))
	require.NoError(t, state.SetLatestBlockHeader(testutil.HydrateBeaconHeader(&ethpb.BeaconBlockHeader{
		Slot:          9,
		ProposerIndex: 0,
	})))

	latestBlockSignedRoot, err := state.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)

	currentEpoch := helpers.CurrentEpoch(state)
	priv, err := bls.RandKey()
	require.NoError(t, err)
	pID, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	block := testutil.NewBeaconBlock()
	block.Block.Slot = 10
	block.Block.ProposerIndex = pID
	block.Block.Body.RandaoReveal = bytesutil.PadTo([]byte{'A', 'B', 'C'}, 96)
	block.Block.ParentRoot = latestBlockSignedRoot[:]
	block.Signature, err = helpers.ComputeDomainAndSign(state, currentEpoch, block.Block, params.BeaconConfig().DomainBeaconProposer, priv)
	require.NoError(t, err)
	proposerIdx, err := helpers.BeaconProposerIndex(state)
	require.NoError(t, err)
	validators[proposerIdx].Slashed = false
	validators[proposerIdx].PublicKey = priv.PublicKey().Marshal()
	err = state.UpdateValidatorAtIndex(proposerIdx, validators[proposerIdx])
	require.NoError(t, err)
	set, err := blocks.BlockSignatureSet(state, block.Block.ProposerIndex, block.Signature, block.Block.HashTreeRoot)
	require.NoError(t, err)

	verified, err := set.Verify()
	require.NoError(t, err)
	assert.Equal(t, true, verified, "Block signature set returned a set which was unable to be verified")
}
