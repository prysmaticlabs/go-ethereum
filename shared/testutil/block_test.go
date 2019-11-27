package testutil

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state/stateutils"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

func TestGenerateFullBlock_PassesStateTransition(t *testing.T) {
	deposits, _, privs := SetupInitialDeposits(t, 128)
	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}
	conf := &BlockGenConfig{
		MaxAttestations: 4,
		Signatures:      true,
	}
	block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot)
	beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateFullBlock_ThousandValidators(t *testing.T) {
	helpers.ClearAllCaches()
	params.OverrideBeaconConfig(params.MinimalSpecConfig())
	defer params.OverrideBeaconConfig(params.MainnetConfig())
	deposits, _, privs := SetupInitialDeposits(t, 1024)
	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}
	conf := &BlockGenConfig{
		MaxAttestations: 16,
		Signatures:      true,
	}
	block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot)
	beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateFullBlock_Passes4Epochs(t *testing.T) {
	helpers.ClearAllCaches()
	// Changing to minimal config as this will process 4 epochs of blocks.
	params.OverrideBeaconConfig(params.MinimalSpecConfig())
	defer params.OverrideBeaconConfig(params.MainnetConfig())
	deposits, _, privs := SetupInitialDeposits(t, 64)
	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}

	conf := &BlockGenConfig{
		MaxAttestations: 2,
		Signatures:      true,
	}
	finalSlot := params.BeaconConfig().SlotsPerEpoch*4 + 3
	for i := 0; i < int(finalSlot); i++ {
		block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot)
		beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Blocks are one slot ahead of beacon state.
	if finalSlot != beaconState.Slot {
		t.Fatalf("expected output slot to be %d, received %d", finalSlot, beaconState.Slot)
	}
	if beaconState.CurrentJustifiedCheckpoint.Epoch != 3 {
		t.Fatalf("expected justified epoch to change to 3, received %d", beaconState.CurrentJustifiedCheckpoint.Epoch)
	}
	if beaconState.FinalizedCheckpoint.Epoch != 2 {
		t.Fatalf("expected finalized epoch to change to 2, received %d", beaconState.CurrentJustifiedCheckpoint.Epoch)
	}
}

func TestGenerateFullBlock_ValidProposerSlashings(t *testing.T) {
	params.OverrideBeaconConfig(params.MinimalSpecConfig())
	defer params.OverrideBeaconConfig(params.MainnetConfig())
	deposits, _, privs := SetupInitialDeposits(t, 32)

	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}
	conf := &BlockGenConfig{
		MaxProposerSlashings: 1,
		Signatures:           true,
	}
	block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot+1)
	beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
	if err != nil {
		t.Fatal(err)
	}

	slashableIndice := block.Body.ProposerSlashings[0].ProposerIndex
	if !beaconState.Validators[slashableIndice].Slashed {
		t.Fatal("expected validator to be slashed")
	}
}

func TestGenerateFullBlock_ValidAttesterSlashings(t *testing.T) {
	params.OverrideBeaconConfig(params.MinimalSpecConfig())
	defer params.OverrideBeaconConfig(params.MainnetConfig())
	deposits, _, privs := SetupInitialDeposits(t, 32)
	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}
	conf := &BlockGenConfig{
		MaxAttesterSlashings: 1,
		Signatures:           true,
	}
	block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot)
	beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
	if err != nil {
		t.Fatal(err)
	}

	slashableIndices := block.Body.AttesterSlashings[0].Attestation_1.CustodyBit_0Indices
	if !beaconState.Validators[slashableIndices[0]].Slashed {
		t.Fatal("expected validator to be slashed")
	}
}

func TestGenerateFullBlock_ValidAttestations(t *testing.T) {
	params.OverrideBeaconConfig(params.MinimalSpecConfig())
	defer params.OverrideBeaconConfig(params.MainnetConfig())
	helpers.ClearAllCaches()
	deposits, _, privs := SetupInitialDeposits(t, 256)

	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}
	conf := &BlockGenConfig{
		MaxAttestations: 4,
		Signatures:      true,
	}
	block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot)
	beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
	if err != nil {
		t.Fatal(err)
	}
	if len(beaconState.CurrentEpochAttestations) != 4 {
		t.Fatal("expected 4 attestations to be saved to the beacon state")
	}
}

func TestGenerateFullBlock_ValidDeposits(t *testing.T) {
	deposits, _, privs := SetupInitialDeposits(t, 256)
	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}
	deposits, _, privs = SetupInitialDeposits(t, 257)
	eth1Data = GenerateEth1Data(t, deposits)
	beaconState.Eth1Data = eth1Data
	conf := &BlockGenConfig{
		MaxDeposits: 1,
		Signatures:  true,
	}
	block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot)
	beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
	if err != nil {
		t.Fatal(err)
	}

	depositedPubkey := block.Body.Deposits[0].Data.PublicKey
	valIndexMap := stateutils.ValidatorIndexMap(beaconState)
	index := valIndexMap[bytesutil.ToBytes48(depositedPubkey)]
	if beaconState.Validators[index].EffectiveBalance != params.BeaconConfig().MaxEffectiveBalance {
		t.Fatalf(
			"expected validator balance to be max effective balance, received %d",
			beaconState.Validators[index].EffectiveBalance,
		)
	}
}

func TestGenerateFullBlock_ValidVoluntaryExits(t *testing.T) {
	deposits, _, privs := SetupInitialDeposits(t, 256)
	eth1Data := GenerateEth1Data(t, deposits)
	beaconState, err := state.GenesisBeaconState(deposits, 0, eth1Data)
	if err != nil {
		t.Fatal(err)
	}
	// Moving the state 2048 epochs forward due to PERSISTENT_COMMITTEE_PERIOD.
	beaconState.Slot = 3 + params.BeaconConfig().PersistentCommitteePeriod*params.BeaconConfig().SlotsPerEpoch
	conf := &BlockGenConfig{
		MaxVoluntaryExits: 1,
		Signatures:        true,
	}
	block := GenerateFullBlock(t, beaconState, privs, conf, beaconState.Slot)
	beaconState, err = state.ExecuteStateTransition(context.Background(), beaconState, block)
	if err != nil {
		t.Fatal(err)
	}

	exitedIndex := block.Body.VoluntaryExits[0].ValidatorIndex
	if beaconState.Validators[exitedIndex].ExitEpoch == params.BeaconConfig().FarFutureEpoch {
		t.Fatal("expected exiting validator index to be marked as exiting")
	}
}
