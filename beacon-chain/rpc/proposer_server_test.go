package rpc

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	b "github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/internal"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"github.com/prysmaticlabs/prysm/shared/featureconfig"
	"github.com/prysmaticlabs/prysm/shared/params"
)

func init() {
	featureconfig.InitFeatureConfig(&featureconfig.FeatureFlagConfig{
		EnableComputeStateRoot: true,
	})
}

func TestProposeBlock_OK(t *testing.T) {
	db := internal.SetupDB(t)
	defer internal.TeardownDB(t, db)
	mockChain := &mockChainService{}
	ctx := context.Background()

	genesis := b.NewGenesisBlock([]byte{})
	if err := db.SaveBlock(genesis); err != nil {
		t.Fatalf("Could not save genesis block: %v", err)
	}

	deposits := make([]*pbp2p.Deposit, params.BeaconConfig().DepositsForChainStart)
	for i := 0; i < len(deposits); i++ {
		depositData, err := helpers.EncodeDepositData(
			&pbp2p.DepositInput{
				Pubkey: []byte(strconv.Itoa(i)),
			},
			params.BeaconConfig().MaxDepositAmount,
			time.Now().Unix(),
		)
		if err != nil {
			t.Fatalf("Could not encode deposit input: %v", err)
		}
		deposits[i] = &pbp2p.Deposit{
			DepositData: depositData,
		}
	}

	beaconState, err := state.GenesisBeaconState(deposits, 0, nil)
	if err != nil {
		t.Fatalf("Could not instantiate genesis state: %v", err)
	}

	if err := db.UpdateChainHead(ctx, genesis, beaconState); err != nil {
		t.Fatalf("Could not save genesis state: %v", err)
	}

	proposerServer := &ProposerServer{
		chainService:    mockChain,
		beaconDB:        db,
		powChainService: &mockPOWChainService{},
	}
	req := &pbp2p.BeaconBlock{
		Slot:             5,
		ParentRootHash32: []byte("parent-hash"),
	}
	if err := proposerServer.beaconDB.SaveBlock(req); err != nil {
		t.Fatal(err)
	}
	if _, err := proposerServer.ProposeBlock(context.Background(), req); err != nil {
		t.Errorf("Could not propose block correctly: %v", err)
	}
}

func TestComputeStateRoot_OK(t *testing.T) {
	db := internal.SetupDB(t)
	defer internal.TeardownDB(t, db)
	ctx := context.Background()

	mockChain := &mockChainService{}

	genesis := b.NewGenesisBlock([]byte{})
	if err := db.SaveBlock(genesis); err != nil {
		t.Fatalf("Could not save genesis block: %v", err)
	}

	deposits := make([]*pbp2p.Deposit, params.BeaconConfig().DepositsForChainStart)
	for i := 0; i < len(deposits); i++ {
		depositData, err := helpers.EncodeDepositData(
			&pbp2p.DepositInput{
				Pubkey: []byte(strconv.Itoa(i)),
			},
			params.BeaconConfig().MaxDepositAmount,
			time.Now().Unix(),
		)
		if err != nil {
			t.Fatalf("Could not encode deposit input: %v", err)
		}
		deposits[i] = &pbp2p.Deposit{
			DepositData: depositData,
		}
	}

	beaconState, err := state.GenesisBeaconState(deposits, 0, nil)
	if err != nil {
		t.Fatalf("Could not instantiate genesis state: %v", err)
	}

	beaconState.Slot = 10

	if err := db.UpdateChainHead(ctx, genesis, beaconState); err != nil {
		t.Fatalf("Could not save genesis state: %v", err)
	}

	proposerServer := &ProposerServer{
		chainService:    mockChain,
		beaconDB:        db,
		powChainService: &mockPOWChainService{},
	}

	req := &pbp2p.BeaconBlock{
		ParentRootHash32: nil,
		Slot:             11,
		RandaoReveal:     nil,
		Body: &pbp2p.BeaconBlockBody{
			ProposerSlashings: nil,
			AttesterSlashings: nil,
		},
	}

	_, _ = proposerServer.ComputeStateRoot(context.Background(), req)
}

func TestPendingAttestations_FiltersWithinInclusionDelay(t *testing.T) {
	db := internal.SetupDB(t)
	defer internal.TeardownDB(t, db)
	ctx := context.Background()

	beaconState := &pbp2p.BeaconState{
		Slot: params.BeaconConfig().GenesisSlot + params.BeaconConfig().MinAttestationInclusionDelay + 100,
	}
	proposerServer := &ProposerServer{
		operationService: &mockOperationService{
			pendingAttestations: []*pbp2p.Attestation{
				{Data: &pbp2p.AttestationData{
					Slot: beaconState.Slot - params.BeaconConfig().MinAttestationInclusionDelay,
				}},
			},
		},
		chainService: &mockChainService{},
		beaconDB:     db,
	}
	if err := db.SaveState(ctx, beaconState); err != nil {
		t.Fatal(err)
	}

	blk := &pbp2p.BeaconBlock{
		Slot: beaconState.Slot,
	}

	if err := db.SaveBlock(blk); err != nil {
		t.Fatalf("failed to save block %v", err)
	}

	if err := db.UpdateChainHead(ctx, blk, beaconState); err != nil {
		t.Fatalf("couldnt update chainhead: %v", err)
	}

	res, err := proposerServer.PendingAttestations(context.Background(), &pb.PendingAttestationsRequest{
		FilterReadyForInclusion: true,
	})
	if err != nil {
		t.Fatalf("Unexpected error fetching pending attestations: %v", err)
	}
	if len(res.PendingAttestations) == 0 {
		t.Error("Expected pending attestations list to be non-empty")
	}
}

func TestPendingAttestations_FiltersExpiredAttestations(t *testing.T) {
	db := internal.SetupDB(t)
	defer internal.TeardownDB(t, db)
	ctx := context.Background()

	// Edge case: current slot is at the end of an epoch. The pending attestation
	// for the next slot should come from currentSlot + 1.
	currentSlot := helpers.StartSlot(
		params.BeaconConfig().GenesisEpoch+10,
	) - 1

	expectedEpoch := uint64(100)

	opService := &mockOperationService{
		pendingAttestations: []*pbp2p.Attestation{
			// Expired attestations
			{Data: &pbp2p.AttestationData{Slot: 0, JustifiedEpoch: expectedEpoch}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 10000, JustifiedEpoch: expectedEpoch}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 5000, JustifiedEpoch: expectedEpoch}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 100, JustifiedEpoch: expectedEpoch}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot - params.BeaconConfig().SlotsPerEpoch, JustifiedEpoch: expectedEpoch}},
			// Non-expired attestation with incorrect justified epoch
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 5, JustifiedEpoch: expectedEpoch - 1}},
			// Non-expired attestations with correct justified epoch
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 5, JustifiedEpoch: expectedEpoch}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 2, JustifiedEpoch: expectedEpoch}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot, JustifiedEpoch: expectedEpoch}},
		},
	}
	expectedNumberOfAttestations := 3
	proposerServer := &ProposerServer{
		operationService: opService,
		chainService:     &mockChainService{},
		beaconDB:         db,
	}
	beaconState := &pbp2p.BeaconState{
		Slot:                   currentSlot,
		JustifiedEpoch:         expectedEpoch,
		PreviousJustifiedEpoch: expectedEpoch,
	}
	if err := db.SaveState(ctx, beaconState); err != nil {
		t.Fatal(err)
	}

	blk := &pbp2p.BeaconBlock{
		Slot: beaconState.Slot,
	}

	if err := db.SaveBlock(blk); err != nil {
		t.Fatalf("failed to save block %v", err)
	}

	if err := db.UpdateChainHead(ctx, blk, beaconState); err != nil {
		t.Fatalf("couldnt update chainhead: %v", err)
	}

	res, err := proposerServer.PendingAttestations(
		context.Background(),
		&pb.PendingAttestationsRequest{
			ProposalBlockSlot: currentSlot,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected error fetching pending attestations: %v", err)
	}
	if len(res.PendingAttestations) != expectedNumberOfAttestations {
		t.Errorf(
			"Expected pending attestations list length %d, but was %d",
			expectedNumberOfAttestations,
			len(res.PendingAttestations),
		)
	}

	expectedAtts := []*pbp2p.Attestation{
		{Data: &pbp2p.AttestationData{Slot: currentSlot - 5, JustifiedEpoch: expectedEpoch}},
		{Data: &pbp2p.AttestationData{Slot: currentSlot - 2, JustifiedEpoch: expectedEpoch}},
		{Data: &pbp2p.AttestationData{Slot: currentSlot, JustifiedEpoch: expectedEpoch}},
	}
	if !reflect.DeepEqual(res.PendingAttestations, expectedAtts) {
		t.Error("Did not receive expected attestations")
	}
}

func TestPendingAttestations_OK(t *testing.T) {
	db := internal.SetupDB(t)
	defer internal.TeardownDB(t, db)
	ctx := context.Background()

	proposerServer := &ProposerServer{
		operationService: &mockOperationService{},
		chainService:     &mockChainService{},
		beaconDB:         db,
	}
	beaconState := &pbp2p.BeaconState{
		Slot: params.BeaconConfig().GenesisSlot + params.BeaconConfig().MinAttestationInclusionDelay,
	}
	if err := db.SaveState(ctx, beaconState); err != nil {
		t.Fatal(err)
	}

	blk := &pbp2p.BeaconBlock{
		Slot: beaconState.Slot,
	}

	if err := db.SaveBlock(blk); err != nil {
		t.Fatalf("failed to save block %v", err)
	}

	if err := db.UpdateChainHead(ctx, blk, beaconState); err != nil {
		t.Fatalf("couldnt update chainhead: %v", err)
	}

	res, err := proposerServer.PendingAttestations(context.Background(), &pb.PendingAttestationsRequest{})
	if err != nil {
		t.Fatalf("Unexpected error fetching pending attestations: %v", err)
	}
	if len(res.PendingAttestations) == 0 {
		t.Error("Expected pending attestations list to be non-empty")
	}
}

func TestPendingAttestations_FiltersCanonicalAttestations(t *testing.T) {
	db := internal.SetupDB(t)
	defer internal.TeardownDB(t, db)
	ctx := context.Background()

	// Edge case: current slot is at the end of an epoch. The pending attestation
	// for the next slot should come from currentSlot + 1.
	currentSlot := helpers.StartSlot(
		params.BeaconConfig().GenesisEpoch+10,
	) - 1

	expectedEpoch := uint64(100)

	opService := &mockOperationService{
		pendingAttestations: []*pbp2p.Attestation{
			// Canonical attestations
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 5, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'A'}}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot - 2, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'B'}}},
			// Non canonical attestations
			{Data: &pbp2p.AttestationData{Slot: currentSlot, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'C'}}},
			{Data: &pbp2p.AttestationData{Slot: currentSlot, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'D'}}},
			// Canonical attestation
			{Data: &pbp2p.AttestationData{Slot: currentSlot, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'E'}}},
		},
	}
	expectedNumberOfAttestations := 3
	proposerServer := &ProposerServer{
		operationService: opService,
		chainService:     &mockChainService{canonicalBlocks: make(map[uint64][]byte)},
		beaconDB:         db,
	}
	beaconState := &pbp2p.BeaconState{
		Slot:                   currentSlot,
		JustifiedEpoch:         expectedEpoch,
		PreviousJustifiedEpoch: expectedEpoch,
	}
	if err := db.SaveState(ctx, beaconState); err != nil {
		t.Fatal(err)
	}

	blk := &pbp2p.BeaconBlock{
		Slot: beaconState.Slot,
	}

	if err := db.SaveBlock(blk); err != nil {
		t.Fatalf("failed to save block %v", err)
	}

	if err := db.UpdateChainHead(ctx, blk, beaconState); err != nil {
		t.Fatalf("couldnt update chainhead: %v", err)
	}
	for _, atts := range opService.pendingAttestations {
		// Insert all of the attestations, since last 3 attestations are same slot number
		// the last 2 get overridden, and only the last 1 is canonical.
		proposerServer.chainService.InsertsCanonical(atts.Data.Slot, atts.Data.BeaconBlockRootHash32)
	}

	res, err := proposerServer.PendingAttestations(
		context.Background(),
		&pb.PendingAttestationsRequest{
			ProposalBlockSlot: currentSlot,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected error fetching pending attestations: %v", err)
	}
	if len(res.PendingAttestations) != expectedNumberOfAttestations {
		t.Errorf(
			"Expected pending attestations list length %d, but was %d",
			expectedNumberOfAttestations,
			len(res.PendingAttestations),
		)
	}

	expectedAtts := []*pbp2p.Attestation{
		{Data: &pbp2p.AttestationData{Slot: currentSlot - 5, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'A'}}},
		{Data: &pbp2p.AttestationData{Slot: currentSlot - 2, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'B'}}},
		{Data: &pbp2p.AttestationData{Slot: currentSlot, JustifiedEpoch: expectedEpoch, BeaconBlockRootHash32: []byte{'E'}}},
	}
	if !reflect.DeepEqual(res.PendingAttestations, expectedAtts) {
		t.Error("Did not receive expected attestations")
	}
}
