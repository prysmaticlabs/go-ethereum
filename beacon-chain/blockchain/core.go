package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/gogo/protobuf/proto"
	"github.com/prysmaticlabs/prysm/beacon-chain/params"
	"github.com/prysmaticlabs/prysm/beacon-chain/types"
	"github.com/prysmaticlabs/prysm/beacon-chain/utils"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/sirupsen/logrus"
)

var activeStateLookupKey = "beacon-active-state"
var crystallizedStateLookupKey = "beacon-crystallized-state"

var clock utils.Clock = &utils.RealClock{}

// BeaconChain represents the core PoS blockchain object containing
// both a crystallized and active state.
type BeaconChain struct {
	state *beaconState
	lock  sync.Mutex
	db    ethdb.Database
}

type beaconState struct {
	// ActiveState captures the beacon state at block processing level,
	// it focuses on verifying aggregated signatures and pending attestations.
	ActiveState *types.ActiveState
	// CrystallizedState captures the beacon state at epoch transition level,
	// it focuses on changes to the validator set, processing cross links and
	// setting up FFG checkpoints.
	CrystallizedState *types.CrystallizedState
}

type BeaconCommittee struct {
	ShardID   int
	Committee []int
}

// NewBeaconChain initializes a beacon chain using genesis state parameters if
// none provided.
func NewBeaconChain(db ethdb.Database) (*BeaconChain, error) {
	beaconChain := &BeaconChain{
		db:    db,
		state: &beaconState{},
	}
	hasActive, err := db.Has([]byte(activeStateLookupKey))
	if err != nil {
		return nil, err
	}
	hasCrystallized, err := db.Has([]byte(crystallizedStateLookupKey))
	if err != nil {
		return nil, err
	}
	hasGenesis, err := db.Has([]byte("genesis"))
	if err != nil {
		return nil, err
	}
	if !hasGenesis {
		log.Info("No genesis block found on disk, initializing genesis block")
		genesisBlock, err := types.NewGenesisBlock()
		if err != nil {
			return nil, err
		}
		genesisMarshall, err := proto.Marshal(genesisBlock.Proto())
		if err != nil {
			return nil, err
		}
		if err := beaconChain.db.Put([]byte("genesis"), genesisMarshall); err != nil {
			return nil, err
		}
	}
	if !hasActive && !hasCrystallized {
		log.Info("No chainstate found on disk, initializing beacon from genesis")
		active, crystallized := types.NewGenesisStates()
		beaconChain.state.ActiveState = active
		beaconChain.state.CrystallizedState = crystallized

		return beaconChain, nil
	}
	if hasActive {
		enc, err := db.Get([]byte(activeStateLookupKey))
		if err != nil {
			return nil, err
		}
		activeData := &pb.ActiveState{}
		err = proto.Unmarshal(enc, activeData)
		if err != nil {
			return nil, err
		}
		beaconChain.state.ActiveState = types.NewActiveState(activeData)
	}
	if hasCrystallized {
		enc, err := db.Get([]byte(crystallizedStateLookupKey))
		if err != nil {
			return nil, err
		}
		crystallizedData := &pb.CrystallizedState{}
		err = proto.Unmarshal(enc, crystallizedData)
		if err != nil {
			return nil, err
		}
		beaconChain.state.CrystallizedState = types.NewCrystallizedState(crystallizedData)
	}
	return beaconChain, nil
}

// GenesisBlock returns the canonical, genesis block.
func (b *BeaconChain) GenesisBlock() (*types.Block, error) {
	genesisExists, err := b.db.Has([]byte("genesis"))
	if err != nil {
		return nil, err
	}
	if genesisExists {
		bytes, err := b.db.Get([]byte("genesis"))
		if err != nil {
			return nil, err
		}
		block := &pb.BeaconBlock{}
		if err := proto.Unmarshal(bytes, block); err != nil {
			return nil, err
		}
		return types.NewBlock(block), nil
	}
	return types.NewGenesisBlock()
}

// ActiveState contains the current state of attestations and changes every block.
func (b *BeaconChain) ActiveState() *types.ActiveState {
	return b.state.ActiveState
}

// CrystallizedState contains epoch dependent validator information, changes every epoch.
func (b *BeaconChain) CrystallizedState() *types.CrystallizedState {
	return b.state.CrystallizedState
}

// SetActiveState is a convenience method which sets and persists the active state on the beacon chain.
func (b *BeaconChain) SetActiveState(activeState *types.ActiveState) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.state.ActiveState = activeState
	return b.PersistActiveState()
}

// SetCrystallizedState is a convenience method which sets and persists the crystallized state on the beacon chain.
func (b *BeaconChain) SetCrystallizedState(crystallizedState *types.CrystallizedState) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.state.CrystallizedState = crystallizedState
	return b.PersistCrystallizedState()
}

// PersistActiveState stores proto encoding of the current beacon chain active state into the db.
func (b *BeaconChain) PersistActiveState() error {
	encodedState, err := b.ActiveState().Marshal()
	if err != nil {
		return err
	}
	return b.db.Put([]byte(activeStateLookupKey), encodedState)
}

// PersistCrystallizedState stores proto encoding of the current beacon chain crystallized state into the db.
func (b *BeaconChain) PersistCrystallizedState() error {
	encodedState, err := b.CrystallizedState().Marshal()
	if err != nil {
		return err
	}
	return b.db.Put([]byte(crystallizedStateLookupKey), encodedState)
}

// IsEpochTransition checks if it's epoch transition time.
func (b *BeaconChain) IsEpochTransition(slotNumber uint64) bool {
	return slotNumber >= b.CrystallizedState().LastStateRecalc()+params.CycleLength
}

// CanProcessBlock is called to decide if an incoming p2p block can be processed into the chain's block trie,
// it checks time stamp, beacon chain parent block hash. It also checks pow chain reference hash if it's a validator.
func (b *BeaconChain) CanProcessBlock(fetcher types.POWBlockFetcher, block *types.Block, isValidator bool) (bool, error) {
	if isValidator {
		if _, err := fetcher.BlockByHash(context.Background(), block.PowChainRef()); err != nil {
			return false, fmt.Errorf("fetching PoW block corresponding to mainchain reference failed: %v", err)
		}
	}

	canProcess, err := b.verifyBlockParentHash(block)
	if err != nil {
		return false, fmt.Errorf("unable to process block: %v", err)
	}
	if !canProcess {
		return false, fmt.Errorf("parent block verification for beacon block %v failed", block.SlotNumber())
	}

	canProcess, err = b.verifyBlockTimeStamp(block)
	if err != nil {
		return false, fmt.Errorf("unable to process block: %v", err)
	}
	if !canProcess {
		return false, fmt.Errorf("time stamp verification for beacon block %v failed", block.SlotNumber())
	}

	canProcess, err = b.verifyBlockActiveHash(block)
	if err != nil {
		return false, fmt.Errorf("unable to process block: %v", err)
	}
	if !canProcess {
		return false, fmt.Errorf("active state verification for beacon block %v failed", block.SlotNumber())
	}

	canProcess, err = b.verifyBlockCrystallizedHash(block)
	if err != nil {
		return false, fmt.Errorf("unable to process block: %v", err)
	}
	if !canProcess {
		return false, fmt.Errorf("crystallized verification for beacon block %v failed", block.SlotNumber())
	}
	return canProcess, nil
}

// verifyBlockTimeStamp verifies node's local time is greater than or equal to
// min timestamp as computed by GENESIS_TIME + slot_number * SLOT_DURATION.
func (b *BeaconChain) verifyBlockTimeStamp(block *types.Block) (bool, error) {
	slotDuration := time.Duration(block.SlotNumber()*params.SlotDuration) * time.Second
	genesis, err := b.GenesisBlock()
	if err != nil {
		return false, err
	}
	genesisTime, err := genesis.Timestamp()
	if err != nil {
		return false, err
	}
	if clock.Now().Before(genesisTime.Add(slotDuration)) {
		return false, nil
	}
	return true, nil
}

// verifyBlockActiveHash verifies block's active state hash equal to
// node's computed active state hash.
func (b *BeaconChain) verifyBlockActiveHash(block *types.Block) (bool, error) {
	hash, err := b.ActiveState().Hash()
	if err != nil {
		return false, err
	}
	if block.ActiveStateHash() != hash {
		return false, nil
	}
	return true, nil
}

// verifyBlockCrystallizedHash verifies block's crystallized state hash equal to
// node's computed crystallized state hash.
func (b *BeaconChain) verifyBlockCrystallizedHash(block *types.Block) (bool, error) {
	hash, err := b.CrystallizedState().Hash()
	if err != nil {
		return false, err
	}
	if block.CrystallizedStateHash() != hash {
		return false, nil
	}
	return true, nil
}

// verifyBlockParentHash verifies parentHash pointed by the beacon block is in the beaconDB.
func (b *BeaconChain) verifyBlockParentHash(block *types.Block) (bool, error) {
	parentHash := block.ParentHash()
	hasParent, err := b.db.Has(parentHash[:])
	if err != nil {
		return false, err
	}
	if !hasParent && block.SlotNumber() != 1 {
		return false, errors.New("parent hash points to nil in beaconDB")
	}
	return true, nil
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

	return types.NewActiveState(&pb.ActiveState{
		PendingAttestations: []*pb.AttestationRecord{},
		RecentBlockHashes:   [][]byte{},
	}), nil
}

// rotateValidatorSet is called every dynasty transition. The primary functions are:
// 1.) Go through queued validator indices and induct them to be active by setting start
// dynasty to current epoch.
// 2.) Remove bad active validator whose balance is below threshold to the exit set by
// setting end dynasty to current epoch.
func (b *BeaconChain) rotateValidatorSet(dynasty uint64) {

	validators := b.CrystallizedState().Validators()
	upperbound := len(activeValidatorIndices(validators, dynasty))/30 + 1

	// Loop through active validator set, remove validator whose balance is below 50%.
	for _, index := range activeValidatorIndices(validators, dynasty) {
		if validators[index].Balance < params.DefaultBalance/2 {
			validators[index].EndDynasty = b.CrystallizedState().CurrentDynasty()
		}
	}
	// Get the total number of validator we can induct.
	inductNum := upperbound
	if len(queuedValidatorIndices(validators, dynasty)) < inductNum {
		inductNum = len(queuedValidatorIndices(validators, dynasty))
	}

	// Induct queued validator to active validator set until the switch dynasty is greater than current number.
	for _, index := range queuedValidatorIndices(validators, dynasty) {
		validators[index].StartDynasty = b.CrystallizedState().CurrentDynasty()
		inductNum--
		if inductNum == 0 {
			break
		}
	}
}

// getAttestersProposer returns lists of random sampled attesters and proposer indices.
func (b *BeaconChain) getAttestersProposer(seed common.Hash) ([]int, int, error) {
	attesterCount := math.Min(params.MinCommiteeSize, float64(b.CrystallizedState().ValidatorsLength()))

	indices, err := utils.ShuffleIndices(seed, b.activeValidatorIndices())
	if err != nil {
		return nil, -1, err
	}
	return indices[:int(attesterCount)], indices[len(indices)-1], nil
}

// getAttestersTotalDeposit returns the total deposit combined by attesters.
// TODO: Consider slashing condition.
func (b *BeaconChain) getAttestersTotalDeposit() uint64 {
	var numOfBits int
	for _, attestation := range b.ActiveState().PendingAttestations() {
		for _, byte := range attestation.AttesterBitfield {
			numOfBits += int(utils.BitSetCount(byte))
		}
	}
	// Assume there's no slashing condition, the following logic will change later phase.
	return uint64(numOfBits) * params.DefaultBalance
}

// calculateRewardsFFG adjusts validators balances by applying rewards or penalties
// based on FFG incentive structure.
func (b *BeaconChain) calculateRewardsFFG(block *types.Block) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	validators := b.CrystallizedState().Validators()
	activeValidators := activeValidatorIndices(validators, b.CrystallizedState().CurrentDynasty())
	attesterDeposits := b.getAttestersTotalDeposit()
	totalDeposit := b.state.CrystallizedState.TotalDeposits()

	attesterFactor := attesterDeposits * 3
	totalFactor := uint64(totalDeposit * 2)
	if attesterFactor >= totalFactor {
		log.Infof("Setting justified epoch to current slot number: %v", block.SlotNumber())
		b.state.CrystallizedState.UpdateJustifiedSlot(block.SlotNumber())

		log.Info("Applying rewards and penalties for the validators from last epoch")

		for i, attesterIndex := range activeValidators {
			voted, err := utils.CheckBit(b.state.ActiveState.LatestPendingAttestation().AttesterBitfield, attesterIndex)
			if err != nil {
				return fmt.Errorf("exiting calculate rewards FFG due to %v", err)
			}
			if voted {
				validators[i].Balance += params.AttesterReward
			} else {
				validators[i].Balance -= params.AttesterReward
			}
		}

		log.Info("Resetting attester bit field to all zeros")
		b.ActiveState().ClearPendingAttestations()
		b.CrystallizedState().SetValidators(validators)
		if err := b.PersistActiveState(); err != nil {
			return err
		}
		if err := b.PersistCrystallizedState(); err != nil {
			return err
		}
	}
	return nil
}

// activeValidatorIndices filters out active validators based on start and end dynasty
// and returns their indices in a list.
func activeValidatorIndices(validators []*pb.ValidatorRecord, dynasty uint64) []int {
	var indices []int
	for i := 0; i < len(validators); i++ {
		if validators[i].StartDynasty <= dynasty && dynasty < validators[i].EndDynasty {
			indices = append(indices, i)
		}
	}
	return indices
}

// exitedValidatorIndices filters out exited validators based on start and end dynasty
// and returns their indices in a list.
func exitedValidatorIndices(validators []*pb.ValidatorRecord, dynasty uint64) []int {
	var indices []int
	for i := 0; i < len(validators); i++ {
		if validators[i].StartDynasty < dynasty && validators[i].EndDynasty < dynasty {
			indices = append(indices, i)
		}
	}
	return indices
}

// queuedValidatorIndices filters out queued validators based on start and end dynasty
// and returns their indices in a list.
func queuedValidatorIndices(validators []*pb.ValidatorRecord, dynasty uint64) []int {
	var indices []int
	for i := 0; i < len(validators); i++ {
		if validators[i].StartDynasty > dynasty {
			indices = append(indices, i)
		}
	}
	return indices
}

// ValidatorsByHeightShard splits a shuffled validator list by height and by shard,
// it ensures there's enough validators per height and per shard, if not, it'll skip
// some heights and shards.
func ValidatorsByHeightShard(seed common.Hash, validators []*pb.ValidatorRecord, dynasty uint64, crosslinkStartShard int) ([]*BeaconCommittee, error) {
	indices := activeValidatorIndices(validators, dynasty)
	var committeesPerSlot int
	var slotsPerCommittee int
	var committees []*BeaconCommittee

	if len(indices) >= params.CycleLength*params.MinCommiteeSize {
		committeesPerSlot = len(indices)/params.CycleLength/(params.MinCommiteeSize*2) + 1
		slotsPerCommittee = 1
	} else {
		committeesPerSlot = 1
		slotsPerCommittee = 1
		for len(indices)*slotsPerCommittee < params.MinCommiteeSize && slotsPerCommittee < params.CycleLength {
			slotsPerCommittee *= 2
		}
	}

	// split the shuffled list for heights.
	shuffledList, err := utils.ShuffleIndices(seed, indices)
	if err != nil {
		return nil, err
	}

	heightList := utils.SplitIndices(shuffledList, params.CycleLength)

	// split the shuffled height list for shards
	for i, subList := range heightList {
		shardList := utils.SplitIndices(subList, params.MinCommiteeSize)
		for _, shardIndex := range shardList {
			shardID := crosslinkStartShard + i*committeesPerSlot/slotsPerCommittee
			committees = append(committees, &BeaconCommittee{
				ShardID:   shardID,
				Committee: shardIndex,
			})
		}
	}
	return committees, nil
}

// GetIndicesForHeight returns the attester set of a given height.
func (b *BeaconChain) GetIndicesForHeight(height uint64) (*pb.ShardAndCommitteeArray, error) {
	lcs := b.CrystallizedState().LastStateRecalc()
	if !(lcs <= height && height < lcs+params.CycleLength*2) {
		return nil, fmt.Errorf("can not return attester set of given height, input height %v has to be in between %v and %v", height, lcs, lcs+params.CycleLength*2)
	}
	return b.CrystallizedState().IndicesForHeights()[height-lcs], nil
}

// GetBlockHash returns the block hash of a given height.
func (b *BeaconChain) GetBlockHash(slot, height uint64) ([]byte, error) {
	sback := slot - params.CycleLength*2
	if !(sback <= height && height < sback+params.CycleLength*2) {
		return nil, fmt.Errorf("can not return attester set of given height, input height %v has to be in between %v and %v", height, sback, sback+params.CycleLength*2)
	}
	return b.ActiveState().RecentBlockHashes()[height-sback].Bytes(), nil
}

// saveBlock puts the passed block into the beacon chain db.
func (b *BeaconChain) saveBlock(block *types.Block) error {
	encodedState, err := block.Marshal()
	if err != nil {
		return err
	}
	hash, err := block.Hash()
	if err != nil {
		return err
	}

	return b.db.Put(hash[:], encodedState)
}
