package powchain

import (
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	contracts "github.com/prysmaticlabs/prysm/contracts/deposit-contract"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/trieutil"
	"github.com/sirupsen/logrus"
	"math/big"
)

var (
	depositEventSignature = []byte("Deposit(bytes,bytes,bytes,bytes,bytes)")
)

// ETH2GenesisTime retrieves the genesis time of the beacon chain
// from a local var.
func (w *Web3Service) ETH2GenesisTime() uint64 {
	return uint64(w.eth2GenesisTime)
}

// HasChainStarted returns the chainStarted var value.
func (w *Web3Service) HasChainStarted() bool {
	return w.chainStarted
}

// ProcessLog is the main method which handles the processing of all
// logs from the deposit contract on the ETH1.0 chain.
func (w *Web3Service) ProcessLog(depositLog gethTypes.Log) error {
	// Process logs according to their event signature.
	if depositLog.Topics[0] == hashutil.HashKeccak256(depositEventSignature) {
		w.processDepositLog(depositLog)
		if !w.chainStarted {
			triggered, err := w.isGenesisTrigger(w.chainStartDeposits)
			if err != nil {
				return err
			}
			if triggered {
				timeStamp := w.blockTime.Unix()
				w.eth2GenesisTime = timeStamp - timeStamp%params.BeaconConfig().SecondsPerDay + 2*params.BeaconConfig().SecondsPerDay
				w.processChainStart(uint64(w.eth2GenesisTime))
				log.Info("Minimum number of validators reached for beacon-chain to start")
				w.chainStarted = true
			}
		}
		return nil
	}
	log.WithField("signature", fmt.Sprintf("%#x", depositLog.Topics[0])).Debug("Not a valid event signature")
	return nil
}

// processDepositLog processes the log which had been received from
// the ETH1.0 chain by trying to ascertain which participant deposited
// in the contract.
func (w *Web3Service) processDepositLog(depositLog gethTypes.Log) {
	pubkey, withdrawalCredentials, amount, signature, merkleTreeIndex, err := contracts.UnpackDepositLogData(depositLog.Data)
	if err != nil {
		log.Errorf("Could not unpack log %v", err)
		return
	}
	// If we have already seen this Merkle index, skip processing the log.
	// This can happen sometimes when we receive the same log twice from the
	// ETH1.0 network, and prevents us from updating our trie
	// with the same log twice, causing an inconsistent state root.
	index := binary.LittleEndian.Uint64(merkleTreeIndex)
	if int64(index) <= w.lastReceivedMerkleIndex {
		return
	}
	w.lastReceivedMerkleIndex = int64(index)

	// We then decode the deposit input in order to create a deposit object
	// we can store in our persistent DB.
	validData := true
	depositData := &pb.DepositData{
		Amount:                bytesutil.FromBytes8(amount),
		Pubkey:                pubkey,
		Signature:             signature,
		WithdrawalCredentials: withdrawalCredentials,
	}

	depositHash, err := hashutil.DepositHash(depositData)
	if err != nil {
		log.Errorf("Unable to determine hashed value of deposit %v", err)
		return
	}

	if err := w.depositTrie.InsertIntoTrie(depositHash[:], int(index)); err != nil {
		log.Errorf("Unable to insert deposit into trie %v", err)
		return
	}

	proof, err := w.depositTrie.MerkleProof(int(index))
	if err != nil {
		log.Errorf("Unable to generate merkle proof for deposit %v", err)
		return
	}

	deposit := &pb.Deposit{
		Data:  depositData,
		Proof: proof,
	}

	// Make sure duplicates are rejected pre-chainstart.
	if !w.chainStarted && validData {
		var pubkey = fmt.Sprintf("#%x", depositData.Pubkey)
		if w.beaconDB.PubkeyInChainstart(w.ctx, pubkey) {
			log.Warnf("Pubkey %#x has already been submitted for chainstart", pubkey)
		} else {
			w.beaconDB.MarkPubkeyForChainstart(w.ctx, pubkey)
		}
	}

	// We always store all historical deposits in the DB.
	w.beaconDB.InsertDeposit(w.ctx, deposit, big.NewInt(int64(depositLog.BlockNumber)), int(index), w.depositTrie.Root())

	if !w.chainStarted {
		w.chainStartDeposits = append(w.chainStartDeposits, deposit)
	} else {
		w.beaconDB.InsertPendingDeposit(w.ctx, deposit, big.NewInt(int64(depositLog.BlockNumber)), int(index), w.depositTrie.Root())
	}
	if validData {
		log.WithFields(logrus.Fields{
			"publicKey":       fmt.Sprintf("%#x", depositData.Pubkey),
			"merkleTreeIndex": index,
		}).Debug("Deposit registered from deposit contract")
		validDepositsCount.Inc()
	} else {
		log.WithFields(logrus.Fields{
			"merkleTreeIndex": index,
		}).Info("Invalid deposit registered in deposit contract")
	}
}

// isGenesisTrigger gets called whenever there's a deposit event,
// it checks whether there's enough effective balance to trigger genesis.
// Spec pseudocode definition:
//   def is_genesis_trigger(deposits: List[Deposit], timestamp: uint64) -> bool:
//    # Process deposits
//    state = BeaconState()
//    for deposit in deposits:
//        process_deposit(state, deposit)
//
//    # Count active validators at genesis
//    active_validator_count = 0
//    for validator in state.validator_registry:
//        if validator.effective_balance == MAX_EFFECTIVE_BALANCE:
//            active_validator_count += 1
//
//    # Check effective balance to trigger genesis
//    GENESIS_ACTIVE_VALIDATOR_COUNT = 2**16
//    return active_validator_count == GENESIS_ACTIVE_VALIDATOR_COUNT
func (w *Web3Service) isGenesisTrigger(deposits []*pb.Deposit) (bool, error) {
	s := state.BeaconState(nil, 0, nil)
	validatorMap := make(map[[32]byte]int)
	var err error
	for _, d := range deposits {
		s, err = blocks.ProcessDeposit(s, d, validatorMap, false, false)
		if err != nil {
			return false, fmt.Errorf("could not process deposit: %v", err)
		}
	}

	activeValidators := uint64(0)
	for _, v := range s.ValidatorRegistry {
		if v.EffectiveBalance == params.BeaconConfig().MaxEffectiveBalance {
			activeValidators++
		}
	}
	triggered := activeValidators == params.BeaconConfig().GenesisActiveValidatorCount
	return triggered, err
}

// processChainStart processes the last final steps before chain start.
func (w *Web3Service) processChainStart(genesisTime uint64) {
	chainStartCount.Inc()

	depHashes, err := w.ChainStartDepositHashes()
	if err != nil {
		log.Errorf("Generating chainstart deposit hashes failed: %v", err)
		return
	}

	// We then update the in-memory deposit trie from the chain start
	// deposits at this point, as this trie will be later needed for
	// incoming, post-chain start deposits.
	sparseMerkleTrie, err := trieutil.GenerateTrieFromItems(
		depHashes,
		int(params.BeaconConfig().DepositContractTreeDepth),
	)
	if err != nil {
		log.Fatalf("Unable to generate deposit trie from ChainStart deposits: %v", err)
	}

	w.depositTrie = sparseMerkleTrie

	log.WithFields(logrus.Fields{
		"ChainStartTime": genesisTime,
	}).Info("Sending chainstart signal to blockchain service")
	w.chainStartFeed.Send(genesisTime)
}

// processPastLogs processes all the past logs from the deposit contract and
// updates the deposit trie with the data from each individual log.
func (w *Web3Service) processPastLogs() error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			w.depositContractAddress,
		},
	}

	logs, err := w.httpLogger.FilterLogs(w.ctx, query)
	if err != nil {
		return err
	}

	for _, l := range logs {
		if err := w.ProcessLog(l); err != nil {
			log.WithError(err).Error("Could not process log")
		}
	}
	w.lastRequestedBlock.Set(w.blockHeight)

	currentState, err := w.beaconDB.HeadState(w.ctx)
	if err != nil {
		return fmt.Errorf("could not get head state: %v", err)
	}

	if currentState != nil && currentState.Eth1DepositIndex > 0 {
		w.beaconDB.PrunePendingDeposits(w.ctx, int(currentState.Eth1DepositIndex))
	}

	return nil
}

// requestBatchedLogs requests and processes all the logs from the period
// last polled to now.
func (w *Web3Service) requestBatchedLogs() error {
	// We request for the nth block behind the current head, in order to have
	// stabilized logs when we retrieve it from the 1.0 chain.
	requestedBlock := big.NewInt(0).Sub(w.blockHeight, big.NewInt(params.BeaconConfig().LogBlockDelay))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			w.depositContractAddress,
		},
		FromBlock: w.lastRequestedBlock.Add(w.lastRequestedBlock, big.NewInt(1)),
		ToBlock:   requestedBlock,
	}
	logs, err := w.httpLogger.FilterLogs(w.ctx, query)
	if err != nil {
		return err
	}

	// Only process log slices which are larger than zero.
	if len(logs) > 0 {
		log.Debug("Processing Batched Logs")
		for _, l := range logs {
			if err := w.ProcessLog(l); err != nil {
				log.WithError(err).Error("Could not process log")
			}
		}
	}
	w.lastRequestedBlock.Set(requestedBlock)
	return nil
}

// ChainStartDepositHashes returns the hashes of all the chainstart deposits
// stored in memory.
func (w *Web3Service) ChainStartDepositHashes() ([][]byte, error) {
	hashes := make([][]byte, len(w.chainStartDeposits))
	for i, dep := range w.chainStartDeposits {
		hash, err := hashutil.DepositHash(dep.Data)
		if err != nil {
			return nil, err
		}
		hashes[i] = hash[:]
	}
	return hashes, nil
}
