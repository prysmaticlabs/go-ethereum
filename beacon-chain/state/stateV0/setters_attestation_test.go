package stateV0

import (
	"context"
	"testing"

	types "github.com/prysmaticlabs/eth2-types"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	eth "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestBeaconState_RotateAttestations(t *testing.T) {
	st, err := InitializeFromProto(&pb.BeaconState{
		Slot:                      1,
		CurrentEpochAttestations:  []*pb.PendingAttestation{{Data: &eth.AttestationData{Slot: 456}}},
		PreviousEpochAttestations: []*pb.PendingAttestation{{Data: &eth.AttestationData{Slot: 123}}},
	})
	require.NoError(t, err)

	require.NoError(t, st.RotateAttestations())
	require.Equal(t, 0, len(st.currentEpochAttestations()))
	require.Equal(t, types.Slot(456), st.previousEpochAttestations()[0].Data.Slot)
}

func TestAppendBeyondIndicesLimit(t *testing.T) {
	zeroHash := params.BeaconConfig().ZeroHash
	mockblockRoots := make([][]byte, params.BeaconConfig().SlotsPerHistoricalRoot)
	for i := 0; i < len(mockblockRoots); i++ {
		mockblockRoots[i] = zeroHash[:]
	}

	mockstateRoots := make([][]byte, params.BeaconConfig().SlotsPerHistoricalRoot)
	for i := 0; i < len(mockstateRoots); i++ {
		mockstateRoots[i] = zeroHash[:]
	}
	mockrandaoMixes := make([][]byte, params.BeaconConfig().EpochsPerHistoricalVector)
	for i := 0; i < len(mockrandaoMixes); i++ {
		mockrandaoMixes[i] = zeroHash[:]
	}
	st, err := InitializeFromProto(&pb.BeaconState{
		Slot:                      1,
		CurrentEpochAttestations:  []*pb.PendingAttestation{{Data: &eth.AttestationData{Slot: 456}}},
		PreviousEpochAttestations: []*pb.PendingAttestation{{Data: &eth.AttestationData{Slot: 123}}},
		Validators:                []*eth.Validator{},
		Eth1Data:                  &eth.Eth1Data{},
		LatestExecutionPayloadHeader: &pb.ExecutionPayloadHeader{
			BlockHash:        make([]byte, 32),
			ParentHash:       make([]byte, 32),
			Coinbase:         make([]byte, 20),
			StateRoot:        make([]byte, 32),
			ReceiptRoot:      make([]byte, 32),
			LogsBloom:        make([]byte, 256),
			TransactionsRoot: make([]byte, 32),
		},
		BlockRoots:  mockblockRoots,
		StateRoots:  mockstateRoots,
		RandaoMixes: mockrandaoMixes,
	})
	require.NoError(t, err)
	_, err = st.HashTreeRoot(context.Background())
	require.NoError(t, err)
	for i := fieldIndex(0); i < fieldIndex(params.BeaconConfig().BeaconStateFieldCount); i++ {
		st.dirtyFields[i] = true
	}
	_, err = st.HashTreeRoot(context.Background())
	require.NoError(t, err)
	for i := 0; i < 10; i++ {
		assert.NoError(t, st.AppendValidator(&eth.Validator{}))
	}
	assert.Equal(t, false, st.rebuildTrie[validators])
	assert.NotEqual(t, len(st.dirtyIndices[validators]), 0)

	for i := 0; i < indicesLimit; i++ {
		assert.NoError(t, st.AppendValidator(&eth.Validator{}))
	}
	assert.Equal(t, true, st.rebuildTrie[validators])
	assert.Equal(t, len(st.dirtyIndices[validators]), 0)
}
