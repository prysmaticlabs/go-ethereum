// Package blockchain defines the life-cycle and status of the beacon chain.
package blockchain

import (
	"context"
	"fmt"

	"github.com/prysmaticlabs/prysm/beacon-chain/powchain"
	"github.com/prysmaticlabs/prysm/beacon-chain/types"
	"github.com/prysmaticlabs/prysm/shared/database"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "blockchain")

// ChainService represents a service that handles the internal
// logic of managing the full PoS beacon chain.
type ChainService struct {
	ctx                                    context.Context
	cancel                                 context.CancelFunc
	beaconDB                               *database.DB
	chain                                  *BeaconChain
	web3Service                            *powchain.Web3Service
	canonicalBlockAnnouncement             chan *types.Block
	canonicalCrystallizedStateAnnouncement chan *types.CrystallizedState
	latestBeaconBlock                      chan *types.Block
	processedBlockHashes                   [][32]byte
	processedActiveStateHashes             [][32]byte
	processedCrystallizedStateHashes       [][32]byte
}

// Config options for the service.
type Config struct {
	BeaconBlockBuf int
}

// DefaultConfig options.
func DefaultConfig() *Config {
	return &Config{
		BeaconBlockBuf: 10,
	}
}

// NewChainService instantiates a new service instance that will
// be registered into a running beacon node.
func NewChainService(ctx context.Context, cfg *Config, beaconDB *database.DB, web3Service *powchain.Web3Service) (*ChainService, error) {
	ctx, cancel := context.WithCancel(ctx)
	return &ChainService{
		ctx:                                    ctx,
		cancel:                                 cancel,
		beaconDB:                               beaconDB,
		web3Service:                            web3Service,
		latestBeaconBlock:                      make(chan *types.Block, cfg.BeaconBlockBuf),
		canonicalBlockAnnouncement:             make(chan *types.Block, 10),
		canonicalCrystallizedStateAnnouncement: make(chan *types.CrystallizedState, 10),
		processedBlockHashes:                   [][32]byte{},
		processedActiveStateHashes:             [][32]byte{},
		processedCrystallizedStateHashes:       [][32]byte{},
	}, nil
}

// Start a blockchain service's main event loop.
func (c *ChainService) Start() {
	log.Infof("Starting service")

	beaconChain, err := NewBeaconChain(c.beaconDB.DB())
	if err != nil {
		log.Errorf("Unable to setup blockchain: %v", err)
	}
	c.chain = beaconChain
	go c.run(c.ctx.Done())
}

// Stop the blockchain service's main event loop and associated goroutines.
func (c *ChainService) Stop() error {
	defer c.cancel()
	log.Info("Stopping service")
	log.Infof("Persisting current active and crystallized states before closing")
	if err := c.chain.PersistActiveState(); err != nil {
		return fmt.Errorf("Error persisting active state: %v", err)
	}
	if err := c.chain.PersistCrystallizedState(); err != nil {
		return fmt.Errorf("Error persisting crystallized state: %v", err)
	}
	return nil
}

// HasStoredState checks if there is any Crystallized/Active State or blocks(not implemented) are
// persisted to the db.
func (c *ChainService) HasStoredState() (bool, error) {

	hasActive, err := c.beaconDB.DB().Has([]byte(activeStateLookupKey))
	if err != nil {
		return false, err
	}
	hasCrystallized, err := c.beaconDB.DB().Has([]byte(crystallizedStateLookupKey))
	if err != nil {
		return false, err
	}
	if !hasActive || !hasCrystallized {
		return false, nil
	}

	return true, nil
}

// ProcessedBlockHashes exposes a getter for the processed block hashes of the chain.
func (c *ChainService) ProcessedBlockHashes() [][32]byte {
	return c.processedBlockHashes
}

// ProcessedCrystallizedStateHashes exposes a getter for the processed crystallized state hashes of the chain.
func (c *ChainService) ProcessedCrystallizedStateHashes() [][32]byte {
	return c.processedCrystallizedStateHashes
}

// ProcessedActiveStateHashes exposes a getter for the processed active state hashes of the chain.
func (c *ChainService) ProcessedActiveStateHashes() [][32]byte {
	return c.processedActiveStateHashes
}

// ProcessBlock accepts a new block for inclusion in the chain.
func (c *ChainService) ProcessBlock(block *types.Block) error {
	h, err := block.Hash()
	if err != nil {
		return fmt.Errorf("could not hash incoming block: %v", err)
	}
	log.WithField("blockHash", fmt.Sprintf("0x%x", h)).Info("Received full block, processing validity conditions")
	canProcess, err := c.chain.CanProcessBlock(c.web3Service.Client(), block)
	if err != nil {
		return err
	}
	if canProcess {
		c.latestBeaconBlock <- block
	}
	return nil
}

// SaveBlock is a mock which saves a block to the local db using the
// blockhash as the key.
func (c *ChainService) SaveBlock(block *types.Block) error {
	return c.chain.saveBlock(block)
}

// ProcessCrystallizedState accepts a new crystallized state object for inclusion in the chain.
// TODO: Implement crystallized state verifier function and apply fork choice rules
func (c *ChainService) ProcessCrystallizedState(state *types.CrystallizedState) error {
	h, err := state.Hash()
	if err != nil {
		return fmt.Errorf("could not hash incoming block: %v", err)
	}
	log.WithField("stateHash", fmt.Sprintf("0x%x", h)).Info("Received crystallized state, processing validity conditions")
	// For now, broadcast all incoming crystallized states to the following channel
	// for gRPC clients to receive then. TODO: change this to actually send
	// canonical crystallized states over the channel.
	c.canonicalCrystallizedStateAnnouncement <- state
	return nil
}

// ProcessActiveState accepts a new active state object for inclusion in the chain.
// TODO: Implement active state verifier function and apply fork choice rules
func (c *ChainService) ProcessActiveState(state *types.ActiveState) error {
	h, err := state.Hash()
	if err != nil {
		return fmt.Errorf("could not hash incoming block: %v", err)
	}
	log.WithField("stateHash", fmt.Sprintf("0x%x", h)).Info("Received active state, processing validity conditions")

	return nil
}

// ContainsBlock checks if a block for the hash exists in the chain.
// This method must be safe to call from a goroutine.
//
// TODO: implement function.
func (c *ChainService) ContainsBlock(h [32]byte) bool {
	return false
}

// ContainsCrystallizedState checks if a crystallized state for the hash exists in the chain.
//
// TODO: implement function.
func (c *ChainService) ContainsCrystallizedState(h [32]byte) bool {
	return false
}

// ContainsActiveState checks if a active state for the hash exists in the chain.
//
// TODO: implement function.
func (c *ChainService) ContainsActiveState(h [32]byte) bool {
	return false
}

// CurrentCrystallizedState of the canonical chain.
func (c *ChainService) CurrentCrystallizedState() *types.CrystallizedState {
	return c.chain.CrystallizedState()
}

// CurrentActiveState of the canonical chain.
func (c *ChainService) CurrentActiveState() *types.ActiveState {
	return c.chain.ActiveState()
}

// CanonicalBlockAnnouncement returns a channel that is written to
// whenever a new block is determined to be canonical in the chain.
func (c *ChainService) CanonicalBlockAnnouncement() <-chan *types.Block {
	return c.canonicalBlockAnnouncement
}

// CanonicalCrystallizedStateAnnouncement returns a channel that is written to
// whenever a new crystallized state is determined to be canonical in the chain.
func (c *ChainService) CanonicalCrystallizedStateAnnouncement() <-chan *types.CrystallizedState {
	return c.canonicalCrystallizedStateAnnouncement
}

// run processes the changes needed every beacon chain block,
// including epoch transition if needed.
func (c *ChainService) run(done <-chan struct{}) {
	for {
		select {
		case block := <-c.latestBeaconBlock:
			// TODO: Using latest block hash for seed, this will eventually be replaced by randao.
			// TODO: Uncomment after there is a reasonable way to bootstrap validators into the
			// protocol. For the first few blocks after genesis, the current approach below
			// will panic as there are no registered validators.
			activeState, err := c.chain.computeNewActiveState(c.web3Service.LatestBlockHash())
			if err != nil {
				log.Errorf("Compute active state failed: %v", err)
			}

			err = c.chain.SetActiveState(activeState)
			if err != nil {
				log.Errorf("Write active state to disk failed: %v", err)
			}

			// TODO: Apply 2.1 fork choice logic using the following.
			validatorsByHeight, err := c.chain.validatorsByHeightShard()
			if err != nil {
				log.Errorf("Unable to get validators by height and by shard: %v", err)
				continue
			}
			log.Debugf("Received the following number of validators by height: %v", len(validatorsByHeight))

			// Entering epoch transitions.
			transition := c.chain.IsEpochTransition(block.SlotNumber())
			if transition {
				c.chain.CrystallizedState().SetStateRecalc(block.SlotNumber())
				if err := c.chain.calculateRewardsFFG(block); err != nil {
					log.Errorf("Error computing validator rewards and penalties %v", err)
					continue
				}
			}
			// Announce the block as "canonical" (TODO: this assumes a fork choice rule
			// occurred successfully).
			c.canonicalBlockAnnouncement <- block

			// SaveBlock to the DB (TODO: this should be done after the fork choice rule and
			// save the fork choice rule).
			if err := c.SaveBlock(block); err != nil {
				log.Errorf("Unable to save block to db: %v", err)
			}

		case <-done:
			log.Debug("Chain service context closed, exiting goroutine")
			return
		}
	}
}
