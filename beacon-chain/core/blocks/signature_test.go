package blocks_test

import (
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestVerifyBlockHeaderSignature(t *testing.T) {
	beaconState, err := testutil.NewBeaconState()
	require.NoError(t, err)

	privKey, err := bls.RandKey()
	require.NoError(t, err)
	validators := make([]*ethpb.Validator, 1)
	validators[0] = &ethpb.Validator{
		PublicKey:             privKey.PublicKey().Marshal(),
		WithdrawalCredentials: make([]byte, 32),
		EffectiveBalance:      params.BeaconConfig().MaxEffectiveBalance,
	}
	err = beaconState.SetValidators(validators)
	require.NoError(t, err)

	// Sign the block header.
	blockHeader := testutil.HydrateSignedBeaconHeader(&ethpb.SignedBeaconBlockHeader{
		Header: &ethpb.BeaconBlockHeader{
			Slot:          0,
			ProposerIndex: 0,
		},
	})
	domain, err := helpers.Domain(
		beaconState.Fork(),
		0,
		params.BeaconConfig().DomainBeaconProposer,
		beaconState.GenesisValidatorRoot(),
	)
	require.NoError(t, err)
	htr, err := blockHeader.Header.HashTreeRoot()
	require.NoError(t, err)
	container := &pb.SigningData{
		ObjectRoot: htr[:],
		Domain:     domain,
	}
	require.NoError(t, err)
	signingRoot, err := container.HashTreeRoot()
	require.NoError(t, err)

	// Set the signature in the block header.
	blockHeader.Signature = privKey.Sign(signingRoot[:]).Marshal()

	// Sig should verify.
	err = blocks.VerifyBlockHeaderSignature(beaconState, blockHeader)
	require.NoError(t, err)
}
