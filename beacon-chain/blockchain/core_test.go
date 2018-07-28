package blockchain

import (
	"bytes"
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/prysmaticlabs/prysm/beacon-chain/database"
	"github.com/prysmaticlabs/prysm/beacon-chain/params"
	"github.com/prysmaticlabs/prysm/beacon-chain/types"
	"github.com/prysmaticlabs/prysm/beacon-chain/utils"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

type faultyFetcher struct{}

func (f *faultyFetcher) BlockByHash(ctx context.Context, hash common.Hash) (*gethTypes.Block, error) {
	return nil, errors.New("cannot fetch block")
}

type mockFetcher struct{}

func (m *mockFetcher) BlockByHash(ctx context.Context, hash common.Hash) (*gethTypes.Block, error) {
	block := gethTypes.NewBlock(&gethTypes.Header{}, nil, nil, nil)
	return block, nil
}

func TestNewBeaconChain(t *testing.T) {
	hook := logTest.NewGlobal()
	config := &database.BeaconDBConfig{DataDir: "", Name: "", InMemory: true}
	db, err := database.NewBeaconDB(config)
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}
	db.Start()
	beaconChain, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("unable to setup beacon chain: %v", err)
	}

	msg := hook.LastEntry().Message
	want := "No chainstate found on disk, initializing beacon from genesis"
	if msg != want {
		t.Errorf("incorrect log, expected %s, got %s", want, msg)
	}

	hook.Reset()
	active, crystallized := types.NewGenesisStates()
	if !reflect.DeepEqual(beaconChain.ActiveState(), active) {
		t.Errorf("active states not equal. received: %v, wanted: %v", beaconChain.ActiveState(), active)
	}
	if !reflect.DeepEqual(beaconChain.CrystallizedState(), crystallized) {
		t.Errorf("crystallized states not equal. received: %v, wanted: %v", beaconChain.CrystallizedState(), crystallized)
	}
}

func TestMutateActiveState(t *testing.T) {
	config := &database.BeaconDBConfig{DataDir: "", Name: "", InMemory: true}
	db, err := database.NewBeaconDB(config)
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}
	db.Start()
	beaconChain, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("unable to setup beacon chain: %v", err)
	}

	active := &types.ActiveState{
		TotalAttesterDeposits: 4096,
		AttesterBitfields:     []byte{'A', 'B', 'C'},
	}
	if err := beaconChain.MutateActiveState(active); err != nil {
		t.Fatalf("unable to mutate active state: %v", err)
	}
	if !reflect.DeepEqual(beaconChain.state.ActiveState, active) {
		t.Errorf("active state was not updated. wanted %v, got %v", active, beaconChain.state.ActiveState)
	}

	// Initializing a new beacon chain should deserialize persisted state from disk.
	newBeaconChain, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("unable to setup second beacon chain: %v", err)
	}
	// The active state should still be the one we mutated and persited earlier.
	if active.TotalAttesterDeposits != newBeaconChain.state.ActiveState.TotalAttesterDeposits {
		t.Errorf("active state height incorrect. wanted %v, got %v", active.TotalAttesterDeposits, newBeaconChain.state.ActiveState.TotalAttesterDeposits)
	}
	if !bytes.Equal(active.AttesterBitfields, newBeaconChain.state.ActiveState.AttesterBitfields) {
		t.Errorf("active state randao incorrect. wanted %v, got %v", active.AttesterBitfields, newBeaconChain.state.ActiveState.AttesterBitfields)
	}
}

func TestMutateCrystallizedState(t *testing.T) {
	config := &database.BeaconDBConfig{DataDir: "", Name: "", InMemory: true}
	db, err := database.NewBeaconDB(config)
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}
	db.Start()
	beaconChain, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("unable to setup beacon chain: %v", err)
	}

	currentCheckpoint := common.BytesToHash([]byte("checkpoint"))
	crystallized := &types.CrystallizedState{
		Dynasty:           3,
		CurrentCheckpoint: currentCheckpoint,
	}
	if err := beaconChain.MutateCrystallizedState(crystallized); err != nil {
		t.Fatalf("unable to mutate crystallized state: %v", err)
	}
	if !reflect.DeepEqual(beaconChain.state.CrystallizedState, crystallized) {
		t.Errorf("crystallized state was not updated. wanted %v, got %v", crystallized, beaconChain.state.CrystallizedState)
	}

	// Initializing a new beacon chain should deserialize persisted state from disk.
	newBeaconChain, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("unable to setup second beacon chain: %v", err)
	}
	// The crystallized state should still be the one we mutated and persited earlier.
	if crystallized.Dynasty != newBeaconChain.state.CrystallizedState.Dynasty {
		t.Errorf("crystallized state dynasty incorrect. wanted %v, got %v", crystallized.Dynasty, newBeaconChain.state.CrystallizedState.Dynasty)
	}
	if crystallized.CurrentCheckpoint.Hex() != newBeaconChain.state.CrystallizedState.CurrentCheckpoint.Hex() {
		t.Errorf("crystallized state current checkpoint incorrect. wanted %v, got %v", crystallized.CurrentCheckpoint.Hex(), newBeaconChain.state.CrystallizedState.CurrentCheckpoint.Hex())
	}
}

func TestGetAttestersProposer(t *testing.T) {
	config := &database.BeaconDBConfig{DataDir: "", Name: "", InMemory: true}
	db, err := database.NewBeaconDB(config)
	if err != nil {
		t.Fatalf("Unable to setup db: %v", err)
	}
	db.Start()
	beaconChain, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("Unable to setup beacon chain: %v", err)
	}

	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}

	var validators []types.ValidatorRecord
	// Create 1000 validators in ActiveValidators.
	for i := 0; i < 1000; i++ {
		validator := types.ValidatorRecord{WithdrawalAddress: common.Address{'A'}, PubKey: enr.Secp256k1(priv.PublicKey)}
		validators = append(validators, validator)
	}

	beaconChain.MutateCrystallizedState(&types.CrystallizedState{ActiveValidators: validators})

	attesters, propser, err := beaconChain.getAttestersProposer(common.Hash{'A'})
	if err != nil {
		t.Errorf("GetAttestersProposer function failed: %v", err)
	}

	validatorList, err := utils.ShuffleIndices(common.Hash{'A'}, len(validators))
	if err != nil {
		t.Errorf("Shuffle function function failed: %v", err)
	}

	if !reflect.DeepEqual(propser, validatorList[len(validatorList)-1]) {
		t.Errorf("Get proposer failed, expected: %v got: %v", validatorList[len(validatorList)-1], propser)
	}
	if !reflect.DeepEqual(attesters, validatorList[:len(attesters)]) {
		t.Errorf("Get attesters failed, expected: %v got: %v", validatorList[:len(attesters)], attesters)
	}
}

func TestCanProcessBlock(t *testing.T) {
	config := &database.BeaconDBConfig{DataDir: "", Name: "", InMemory: true}
	db, err := database.NewBeaconDB(config)
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}
	db.Start()
	beaconChain, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("Unable to setup beacon chain: %v", err)
	}

	block := types.NewBlock(1)
	// Using a faulty fetcher should throw an error.
	if _, err := beaconChain.CanProcessBlock(&faultyFetcher{}, block); err == nil {
		t.Errorf("Using a faulty fetcher should throw an error, received nil")
	}
	activeState := &types.ActiveState{TotalAttesterDeposits: 10000}
	beaconChain.state.ActiveState = activeState

	activeHash, err := hashActiveState(*activeState)
	if err != nil {
		t.Fatalf("Cannot hash active state: %v", err)
	}

	block.InsertActiveHash(activeHash)

	crystallizedHash, err := hashCrystallizedState(types.CrystallizedState{})
	if err != nil {
		t.Fatalf("Compute crystallized state hash failed: %v", err)
	}
	block.InsertCrystallizedHash(crystallizedHash)
	canProcess, err := beaconChain.CanProcessBlock(&mockFetcher{}, block)
	if err != nil {
		t.Fatalf("CanProcessBlocks failed: %v", err)
	}
	if !canProcess {
		t.Errorf("Should be able to process block, could not")
	}

	// Attempting to try a block with that fails the timestamp validity
	// condition.
	block = types.NewBlock(1000000)
	block.InsertActiveHash(activeHash)
	block.InsertCrystallizedHash(crystallizedHash)
	canProcess, err = beaconChain.CanProcessBlock(&mockFetcher{}, block)
	if err != nil {
		t.Fatalf("CanProcessBlocks failed: %v", err)
	}
	if canProcess {
		t.Errorf("Should not be able to process block with invalid timestamp condition")
	}
}

func TestProcessBlockWithBadHashes(t *testing.T) {
	config := &database.BeaconDBConfig{DataDir: "", Name: "", InMemory: true}
	db, err := database.NewBeaconDB(config)
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}
	db.Start()
	b, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("Unable to setup beacon chain: %v", err)
	}

	// Test negative scenario where active state hash is different than node's compute
	block := types.NewBlock(1)
	activeState := &types.ActiveState{TotalAttesterDeposits: 10000}
	stateHash, err := hashActiveState(*activeState)
	if err != nil {
		t.Fatalf("Cannot hash active state: %v", err)
	}
	block.InsertActiveHash(stateHash)

	b.state.ActiveState = &types.ActiveState{TotalAttesterDeposits: 9999}

	canProcess, err := b.CanProcessBlock(&mockFetcher{}, block)
	if err == nil {
		t.Fatalf("CanProcessBlocks should have failed with diff state hashes")
	}
	if canProcess {
		t.Errorf("CanProcessBlocks should have returned false")
	}

	// Test negative scenario where crystallized state hash is different than node's compute
	crystallizedState := &types.CrystallizedState{CurrentEpoch: 10000}
	stateHash, err = hashCrystallizedState(*crystallizedState)
	if err != nil {
		t.Fatalf("Cannot hash crystallized state: %v", err)
	}
	block.InsertCrystallizedHash(stateHash)

	b.state.CrystallizedState = &types.CrystallizedState{CurrentEpoch: 9999}

	canProcess, err = b.CanProcessBlock(&mockFetcher{}, block)
	if err == nil {
		t.Fatalf("CanProcessBlocks should have failed with diff state hashes")
	}
	if canProcess {
		t.Errorf("CanProcessBlocks should have returned false")
	}
}

func TestRotateValidatorSet(t *testing.T) {
	config := &database.BeaconDBConfig{DataDir: "", Name: "", InMemory: true}
	db, err := database.NewBeaconDB(config)
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}
	db.Start()
	b, err := NewBeaconChain(db.DB())
	if err != nil {
		t.Fatalf("Unable to setup beacon chain: %v", err)
	}

	activeValidators := []types.ValidatorRecord{
		{Balance: 10000, WithdrawalAddress: common.Address{'A'}},
		{Balance: 15000, WithdrawalAddress: common.Address{'B'}},
		{Balance: 20000, WithdrawalAddress: common.Address{'C'}},
		{Balance: 25000, WithdrawalAddress: common.Address{'D'}},
		{Balance: 30000, WithdrawalAddress: common.Address{'E'}},
	}

	queuedValidators := []types.ValidatorRecord{
		{Balance: params.DefaultBalance, WithdrawalAddress: common.Address{'F'}},
		{Balance: params.DefaultBalance, WithdrawalAddress: common.Address{'G'}},
		{Balance: params.DefaultBalance, WithdrawalAddress: common.Address{'H'}},
		{Balance: params.DefaultBalance, WithdrawalAddress: common.Address{'I'}},
		{Balance: params.DefaultBalance, WithdrawalAddress: common.Address{'J'}},
	}

	exitedValidators := []types.ValidatorRecord{
		{Balance: 99999, WithdrawalAddress: common.Address{'K'}},
		{Balance: 99999, WithdrawalAddress: common.Address{'L'}},
		{Balance: 99999, WithdrawalAddress: common.Address{'M'}},
		{Balance: 99999, WithdrawalAddress: common.Address{'N'}},
		{Balance: 99999, WithdrawalAddress: common.Address{'O'}},
	}

	b.CrystallizedState().ActiveValidators = activeValidators
	b.CrystallizedState().QueuedValidators = queuedValidators
	b.CrystallizedState().ExitedValidators = exitedValidators

	if b.ActiveValidatorCount() != 5 {
		t.Errorf("Get active validator count failed, wanted 5, got %v", b.ActiveValidatorCount())
	}
	if b.QueuedValidatorCount() != 5 {
		t.Errorf("Get queued validator count failed, wanted 5, got %v", b.QueuedValidatorCount())
	}
	if b.ExitedValidatorCount() != 5 {
		t.Errorf("Get exited validator count failed, wanted 5, got %v", b.ExitedValidatorCount())
	}

	newQueuedValidators, newActiveValidators, newExitedValidators := b.RotateValidatorSet()

	if len(newActiveValidators) != 4 {
		t.Errorf("Get active validator count failed, wanted 5, got %v", len(newActiveValidators))
	}
	if len(newQueuedValidators) != 4 {
		t.Errorf("Get queued validator count failed, wanted 4, got %v", len(newQueuedValidators))
	}
	if len(newExitedValidators) != 7 {
		t.Errorf("Get exited validator count failed, wanted 6, got %v", len(newExitedValidators))
	}
}

func TestCutOffValidatorSet(t *testing.T) {

	// Test scenario #1: Assume there's enough validators to fill in all the heights.
	validatorCount := params.EpochLength * params.MinCommiteeSize
	cutoffsValidators := GetCutoffs(validatorCount)

	// The length of cutoff list should be 65. Since there is 64 heights per epoch,
	// it means during every height, a new set of 128 validators will form a committee.
	expectedCount := int(math.Ceil(float64(validatorCount)/params.MinCommiteeSize)) + 1
	if len(cutoffsValidators) != expectedCount {
		t.Errorf("Incorrect count for cutoffs validator. Wanted: %v, Got: %v", expectedCount, len(cutoffsValidators))
	}

	// Verify each cutoff is an increment of MinCommiteeSize, it means 128 validators forms a
	// a committee and get to attest per height.
	count := 0
	for _, cutoff := range cutoffsValidators {
		if cutoff != count {
			t.Errorf("cutoffsValidators did not get 128 increment. Wanted: count, Got: %v", cutoff)
		}
		count += params.MinCommiteeSize
	}

	// Test scenario #2: Assume there's not enough validators to fill in all the heights.
	validatorCount = 1000
	cutoffsValidators = unique(GetCutoffs(validatorCount))
	// With 1000 validators, we can't attest every height. Given min committee size is 128,
	// we can only attest 7 heights. round down 1000 / 128 equals to 7, means the length is 8.
	expectedCount = int(math.Ceil(float64(validatorCount) / params.MinCommiteeSize))
	if len(unique(cutoffsValidators)) != expectedCount {
		t.Errorf("Incorrect count for cutoffs validator. Wanted: %v, Got: %v", expectedCount, validatorCount/params.MinCommiteeSize)
	}

	// Verify each cutoff is an increment of 142~143 (1000 / 7).
	count = 0
	for _, cutoff := range cutoffsValidators {
		num := count * validatorCount / (len(cutoffsValidators) - 1)
		if cutoff != num {
			t.Errorf("cutoffsValidators did not get correct increment. Wanted: %v, Got: %v", num, cutoff)
		}
		count++
	}
}

// helper function to remove duplicates in a int slice.
func unique(ints []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, int := range ints {
		if _, value := keys[int]; !value {
			keys[int] = true
			list = append(list, int)
		}
	}
	return list
}
