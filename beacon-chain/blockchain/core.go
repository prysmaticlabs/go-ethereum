package blockchain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash"
	"math"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/prysmaticlabs/prysm/beacon-chain/params"
	"github.com/prysmaticlabs/prysm/beacon-chain/powchain"
	"github.com/prysmaticlabs/prysm/beacon-chain/types"
	"github.com/prysmaticlabs/prysm/beacon-chain/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/blake2b"
)

var stateLookupKey = "beaconchainstate"

// BeaconChain represents the core PoS blockchain object containing
// both a crystallized and active state.
type BeaconChain struct {
	state *beaconState
	lock  sync.Mutex
	db    ethdb.Database
}

type beaconState struct {
	ActiveState       *types.ActiveState
	CrystallizedState *types.CrystallizedState
}

// NewBeaconChain initializes an instance using genesis state parameters if
// none provided.
func NewBeaconChain(db ethdb.Database) (*BeaconChain, error) {
	beaconChain := &BeaconChain{
		db:    db,
		state: &beaconState{},
	}
	has, err := db.Has([]byte(stateLookupKey))
	if err != nil {
		return nil, err
	}
	if !has {
		log.Info("No chainstate found on disk, initializing beacon from genesis")
		active, crystallized := types.NewGenesisStates()
		beaconChain.state.ActiveState = active
		beaconChain.state.CrystallizedState = crystallized
		return beaconChain, nil
	}
	enc, err := db.Get([]byte(stateLookupKey))
	if err != nil {
		return nil, err
	}
	// Deserializes the encoded object into a beacon chain.
	if err := rlp.DecodeBytes(enc, &beaconChain.state); err != nil {
		return nil, fmt.Errorf("could not deserialize chainstate from disk: %v", err)
	}
	return beaconChain, nil
}

// ActiveState exposes a getter to external services.
func (b *BeaconChain) ActiveState() *types.ActiveState {
	return b.state.ActiveState
}

// CrystallizedState exposes a getter to external services.
func (b *BeaconChain) CrystallizedState() *types.CrystallizedState {
	return b.state.CrystallizedState
}

// ActiveValidatorCount exposes a getter to total number of active validator.
func (b *BeaconChain) ActiveValidatorCount() int {
	return len(b.state.CrystallizedState.ActiveValidators)
}

// QueuedValidatorCount exposes a getter to total number of queued validator.
func (b *BeaconChain) QueuedValidatorCount() int {
	return len(b.state.CrystallizedState.QueuedValidators)
}

// ExitedValidatorCount exposes a getter to total number of exited validator.
func (b *BeaconChain) ExitedValidatorCount() int {
	return len(b.state.CrystallizedState.ExitedValidators)
}

// GenesisBlock returns the canonical, genesis block.
func (b *BeaconChain) GenesisBlock() *types.Block {
	return types.NewGenesisBlock()
}

func (b *BeaconChain) isEpochTransition(slotNumber uint64) bool {
	currentEpoch := b.state.CrystallizedState.CurrentEpoch
	isTransition := (slotNumber / params.SlotLength) > currentEpoch
	return isTransition
}

// MutateActiveState allows external services to modify the active state.
func (b *BeaconChain) MutateActiveState(activeState *types.ActiveState) error {
	defer b.lock.Unlock()
	b.lock.Lock()
	b.state.ActiveState = activeState
	return b.persist()
}

// MutateCrystallizedState allows external services to modify the crystallized state.
func (b *BeaconChain) MutateCrystallizedState(crystallizedState *types.CrystallizedState) error {
	defer b.lock.Unlock()
	b.lock.Lock()
	b.state.CrystallizedState = crystallizedState
	return b.persist()
}

// CanProcessBlock decides if an incoming p2p block can be processed into the chain's block trie.
func (b *BeaconChain) CanProcessBlock(fetcher powchain.POWBlockFetcher, block *types.Block) (bool, error) {
	mainchainBlock, err := fetcher.BlockByHash(context.Background(), block.Data().MainChainRef)
	if err != nil {
		return false, err
	}
	// There needs to be a valid mainchain block for the reference hash in a beacon block.
	if mainchainBlock == nil {
		return false, nil
	}
	// TODO: check if the parentHash pointed by the beacon block is in the beaconDB.

	// Calculate the timestamp validity condition.
	slotDuration := time.Duration(block.Data().SlotNumber*params.SlotLength) * time.Second
	validTime := time.Now().After(b.GenesisBlock().Data().Timestamp.Add(slotDuration))

	// Verify state hashes from the block are correct
	hash, err := hashActiveState(*b.ActiveState())
	if err != nil {
		return false, err
	}

	if !bytes.Equal(block.Data().ActiveStateHash.Sum(nil), hash.Sum(nil)) {
		return false, fmt.Errorf("Active state hash mismatched, wanted: %v, got: %v", hash.Sum(nil), block.Data().ActiveStateHash.Sum(nil))
	}
	hash, err = hashCrystallizedState(*b.CrystallizedState())
	if err != nil {
		return false, err
	}
	if !bytes.Equal(block.Data().CrystallizedStateHash.Sum(nil), hash.Sum(nil)) {
		return false, fmt.Errorf("Crystallized state hash mismatched, wanted: %v, got: %v", hash.Sum(nil), block.Data().CrystallizedStateHash.Sum(nil))
	}

	return validTime, nil
}

// RotateValidatorSet is called  every dynasty transition. It's primary function is
// to go through queued validators and induct them to be active, and remove bad
// active validator whose balance is below threshold to the exit set. It also cross checks
// every validator's switch dynasty before induct or remove.
func (b *BeaconChain) RotateValidatorSet() ([]types.ValidatorRecord, []types.ValidatorRecord, []types.ValidatorRecord) {

	var newExitedValidators = b.CrystallizedState().ExitedValidators
	var newActiveValidators []types.ValidatorRecord
	upperbound := b.ActiveValidatorCount()/30 + 1
	exitCount := 0

	// Loop through active validator set, remove validator whose balance is below 50% and switch dynasty > current dynasty.
	for _, validator := range b.state.CrystallizedState.ActiveValidators {
		if validator.Balance < params.DefaultBalance/2 {
			newExitedValidators = append(newExitedValidators, validator)
		} else if validator.SwitchDynasty == b.CrystallizedState().Dynasty+1 && exitCount < upperbound {
			newExitedValidators = append(newExitedValidators, validator)
			exitCount++
		} else {
			newActiveValidators = append(newActiveValidators, validator)
		}
	}
	// Get the total number of validator we can induct.
	inductNum := upperbound
	if b.QueuedValidatorCount() < inductNum {
		inductNum = b.QueuedValidatorCount()
	}

	// Induct queued validator to active validator set until the switch dynasty is greater than current number.
	for i := 0; i < inductNum; i++ {
		if b.CrystallizedState().QueuedValidators[i].SwitchDynasty > b.CrystallizedState().Dynasty+1 {
			inductNum = i
			break
		}
		newActiveValidators = append(newActiveValidators, b.CrystallizedState().QueuedValidators[i])
	}
	newQueuedValidators := b.CrystallizedState().QueuedValidators[inductNum:]

	return newQueuedValidators, newActiveValidators, newExitedValidators
}

// persist stores the RLP encoding of the latest beacon chain state into the db.
func (b *BeaconChain) persist() error {
	encodedState, err := rlp.EncodeToBytes(b.state)
	if err != nil {
		return err
	}
	return b.db.Put([]byte(stateLookupKey), encodedState)
}

// computeNewActiveState computes a new active state for every beacon block.
func (b *BeaconChain) computeNewActiveState(seed common.Hash) (*types.ActiveState, error) {
	attesters, proposer, err := b.getAttestersProposer(seed)
	if err != nil {
		return nil, err
	}
	// TODO: Verify attestations from attesters.
	log.WithFields(logrus.Fields{"attestersIndices": attesters}).Debug("Attester indices")

	// TODO: Verify main signature from proposer.
	log.WithFields(logrus.Fields{"proposerIndex": proposer}).Debug("Proposer index")

	// TODO: Update crosslink records (post Ruby release).

	// TODO: Track reward for the proposer that just proposed the latest beacon block.

	// TODO: Verify randao reveal from validator's hash pre image.

	return &types.ActiveState{
		TotalAttesterDeposits: 0,
		AttesterBitfields:     []byte{},
	}, nil
}

// hashActiveState serializes the active state object then uses
// blake2b to hash the serialized object.
func hashActiveState(state types.ActiveState) (hash.Hash, error) {
	serializedState, err := rlp.EncodeToBytes(state)
	if err != nil {
		return nil, err
	}
	return blake2b.New256(serializedState)
}

// hashCrystallizedState serializes the crystallized state object
// then uses blake2b to hash the serialized object.
func hashCrystallizedState(state types.CrystallizedState) (hash.Hash, error) {
	serializedState, err := rlp.EncodeToBytes(state)
	if err != nil {
		return nil, err
	}
	return blake2b.New256(serializedState)
}

// getAttestersProposer returns lists of random sampled attesters and proposer indices.
func (b *BeaconChain) getAttestersProposer(seed common.Hash) ([]int, int, error) {
	attesterCount := math.Min(params.AttesterCount, float64(len(b.CrystallizedState().ActiveValidators)))
	indices, err := utils.ShuffleIndices(seed, len(b.CrystallizedState().ActiveValidators))
	if err != nil {
		return nil, -1, err
	}
	return indices[:int(attesterCount)], indices[len(indices)-1], nil
}

// hasVoted checks if the attester has voted by looking at the bitfield.
func hasVoted(bitfields []byte, attesterBlock int, attesterFieldIndex int) bool {
	voted := false

	fields := bitfields[attesterBlock-1]
	attesterField := fields >> (8 - uint(attesterFieldIndex))
	if attesterField%2 != 0 {
		voted = true
	}

	return voted
}

// applyRewardAndPenalty applies the appropriate rewards and penalties according to
// whether the attester has voted or not.
func (b *BeaconChain) applyRewardAndPenalty(index int, voted bool) error {
	defer b.lock.Unlock()
	b.lock.Lock()

	if voted {
		b.state.CrystallizedState.ActiveValidators[index].Balance += params.AttesterReward
	} else {
		// TODO : Change this when penalties are specified for not voting
		b.state.CrystallizedState.ActiveValidators[index].Balance -= params.AttesterReward
	}

	return b.persist()
}

// resetAttesterBitfields resets the attester bitfields in the ActiveState to zero.
func (b *BeaconChain) resetAttesterBitfields() error {

	bitfields := b.state.ActiveState.AttesterBitfields
	length := int(len(bitfields) / 8)
	if len(bitfields)%8 != 0 {
		length += 1
	}

	defer b.lock.Unlock()
	b.lock.Lock()

	newbitfields := make([]byte, length)
	b.state.ActiveState.AttesterBitfields = newbitfields

	return b.persist()
}

// resetTotalDeposit clears and resets the total attester deposit to zero.
func (b *BeaconChain) resetTotalDeposit() error {
	defer b.lock.Unlock()
	b.lock.Lock()
	b.state.ActiveState.TotalAttesterDeposits = 0

	return b.persist()
}

// setJustifiedEpoch sets the justified epoch during an epoch transition.
func (b *BeaconChain) setJustifiedEpoch() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	justifiedEpoch := b.state.CrystallizedState.LastJustifiedEpoch
	b.state.CrystallizedState.LastJustifiedEpoch = b.state.CrystallizedState.CurrentEpoch

	if b.state.CrystallizedState.CurrentEpoch == (justifiedEpoch + 1) {
		b.state.CrystallizedState.LastFinalizedEpoch = justifiedEpoch
	}

	return b.persist()
}

// setRewardsAndPenalties checks if the attester has voted and then applies the
// rewards and penalties for them.
func (b *BeaconChain) setRewardsAndPenalties(index int) error {
	bitfields := b.state.ActiveState.AttesterBitfields
	attesterBlock := (index + 1) / 8
	attesterFieldIndex := (index + 1) % 8
	if attesterFieldIndex == 0 {
		attesterFieldIndex = 8
	}

	if len(bitfields) < attesterBlock {
		return errors.New("attester index does not exist")
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	voted := hasVoted(bitfields, attesterBlock, attesterFieldIndex)
	if err := b.applyRewardAndPenalty(index, voted); err != nil {
		return fmt.Errorf("unable to apply rewards and penalties: %v", err)
	}

	return nil
}

// Slashing Condtions
// TODO: Implement all the conditions when the spec is updated
func (b *BeaconChain) heightEquivocationCondition(validatorIndex int) error {
	testbool := false
	testSlash := uint64(1)
	if !testbool {
		if err := b.slashStake(validatorIndex, testSlash); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconChain) ffgSurroundCondition(validatorIndex int) error {
	testbool := false
	testSlash := uint64(1)
	if !testbool {
		if err := b.slashStake(validatorIndex, testSlash); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconChain) beaconProposalCondition(validatorIndex int) error {
	testbool := false
	testSlash := uint64(1)
	if !testbool {
		if err := b.slashStake(validatorIndex, testSlash); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconChain) pocSecretLeakCondtion(validatorIndex int) error {
	testbool := false
	testSlash := uint64(1)
	if !testbool {
		if err := b.slashStake(validatorIndex, testSlash); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconChain) pocWrongCustodyCondtion(validatorIndex int) error {
	testbool := false
	testSlash := uint64(1)
	if !testbool {
		if err := b.slashStake(validatorIndex, testSlash); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconChain) pocNoRevealCondtion(validatorIndex int) error {
	testbool := false
	testSlash := uint64(1)
	if !testbool {
		if err := b.slashStake(validatorIndex, testSlash); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconChain) randaoLeakCondtion(validatorIndex int) error {
	testbool := false
	testSlash := uint64(1)
	if !testbool {
		if err := b.slashStake(validatorIndex, testSlash); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeaconChain) slashStake(validatorIndex int, stakeToBeSlashed uint64) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.state.CrystallizedState.ActiveValidators[validatorIndex].Balance -= stakeToBeSlashed
	return b.persist()
}

// This will be changed or removed once the spec is more defined, this is just a placeholder
// to apply all the slashing conditions
func (b *BeaconChain) applySlashingConditions(validatorIndex int) error {

	if err := b.heightEquivocationCondition(validatorIndex); err != nil {
		return fmt.Errorf("unable to apply height equivocation condition: %v", err)
	}

	if err := b.ffgSurroundCondition(validatorIndex); err != nil {
		return fmt.Errorf("unable to apply FFG surround condition: %v", err)
	}

	if err := b.beaconProposalCondition(validatorIndex); err != nil {
		return fmt.Errorf("unable to apply Beacon Proposal condition: %v", err)
	}

	if err := b.pocSecretLeakCondtion(validatorIndex); err != nil {
		return fmt.Errorf("unable to apply Secret Leak condition: %v", err)
	}

	if err := b.pocWrongCustodyCondtion(validatorIndex); err != nil {
		return fmt.Errorf("unable to apply POC wrong custody condition: %v", err)
	}

	if err := b.pocNoRevealCondtion(validatorIndex); err != nil {
		return fmt.Errorf("unable to apply POC no reveal condition: %v", err)
	}

	if err := b.randaoLeakCondtion(validatorIndex); err != nil {
		return fmt.Errorf("unable to apply randao leak condition: %v", err)
	}

	return nil
}

// computeValidatorRewardsAndPenalties is run every epoch transition and appropriates the
// rewards and penalties, resets the bitfield and deposits and also applies the slashing conditions.
func (b *BeaconChain) computeValidatorRewardsAndPenalties() error {
	activeValidatorSet := b.state.CrystallizedState.ActiveValidators
	attesterDeposits := b.state.ActiveState.TotalAttesterDeposits
	totalDeposit := b.state.CrystallizedState.TotalDeposits

	attesterFactor := attesterDeposits * 3
	totalFactor := uint64(totalDeposit * 2)

	if attesterFactor >= totalFactor {
		log.Info("Justified epoch in the crystallised state is set to the current epoch")

		if err := b.setJustifiedEpoch(); err != nil {
			return fmt.Errorf("error setting justified epoch: %v", err)
		}

		for i, _ := range activeValidatorSet {
			if err := b.setRewardsAndPenalties(i); err != nil {
				log.Error(err)
			}

			if err := b.applySlashingConditions(i); err != nil {
				log.Error(err)
			}
		}

		if err := b.resetAttesterBitfields(); err != nil {
			return fmt.Errorf("error resetting bitfields: %v", err)
		}
		if err := b.resetTotalDeposit(); err != nil {
			return fmt.Errorf("error resetting total deposits: %v", err)
		}
	}
	return nil
}
