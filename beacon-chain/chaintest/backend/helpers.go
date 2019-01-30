package backend

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	b "github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/trie"
)

// Generates a simulated beacon block to use
// in the next state transition given the current state,
// the previous beacon block, and previous beacon block root.
func generateSimulatedBlock(
	beaconState *pb.BeaconState,
	prevBlockRoot [32]byte,
	randaoReveal [32]byte,
	depositRandaoCommit [32]byte,
	simulatedDeposit *StateTestDeposit,
	depositsTrie *trie.DepositTrie,
	simulatedProposerSlashing *StateTestProposerSlashing,
	simulatedAttesterSlashing *StateTestAttesterSlashing,
	simulatedExit *StateTestValidatorExit,
) (*pb.BeaconBlock, [32]byte, error) {
	encodedState, err := proto.Marshal(beaconState)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("could not marshal beacon state: %v", err)
	}
	stateRoot := hashutil.Hash(encodedState)
	block := &pb.BeaconBlock{
		Slot:               beaconState.Slot + 1,
		RandaoRevealHash32: randaoReveal[:],
		ParentRootHash32:   prevBlockRoot[:],
		StateRootHash32:    stateRoot[:],
		Eth1Data: &pb.Eth1Data{
			DepositRootHash32: []byte{1},
			BlockHash32:       []byte{2},
		},
		Body: &pb.BeaconBlockBody{
			ProposerSlashings: []*pb.ProposerSlashing{},
			AttesterSlashings: []*pb.AttesterSlashing{},
			Attestations:      []*pb.Attestation{},
			Deposits:          []*pb.Deposit{},
			Exits:             []*pb.Exit{},
		},
	}
	if simulatedDeposit != nil {
		depositInput := &pb.DepositInput{
			Pubkey:                      []byte(simulatedDeposit.Pubkey),
			WithdrawalCredentialsHash32: []byte{},
			ProofOfPossession:           []byte{},
			RandaoCommitmentHash32:      depositRandaoCommit[:],
			CustodyCommitmentHash32:     []byte{},
		}

		data, err := b.EncodeDepositData(depositInput, simulatedDeposit.Amount, time.Now().Unix())
		if err != nil {
			return nil, [32]byte{}, fmt.Errorf("could not encode deposit data: %v", err)
		}

		// We then update the deposits Merkle trie with the deposit data and return
		// its Merkle branch leading up to the root of the trie.
		depositsTrie.UpdateDepositTrie(data)
		merkleBranch := depositsTrie.GenerateMerkleBranch(simulatedDeposit.MerkleIndex)

		block.Body.Deposits = append(block.Body.Deposits, &pb.Deposit{
			DepositData:         data,
			MerkleBranchHash32S: merkleBranch,
			MerkleTreeIndex:     simulatedDeposit.MerkleIndex,
		})
	}
	if simulatedProposerSlashing != nil {
		block.Body.ProposerSlashings = append(block.Body.ProposerSlashings, &pb.ProposerSlashing{
			ProposerIndex: simulatedProposerSlashing.ProposerIndex,
			ProposalData_1: &pb.ProposalSignedData{
				Slot:            simulatedProposerSlashing.Proposal1Slot,
				Shard:           simulatedProposerSlashing.Proposal1Shard,
				BlockRootHash32: []byte(simulatedProposerSlashing.Proposal1Root),
			},
			ProposalData_2: &pb.ProposalSignedData{
				Slot:            simulatedProposerSlashing.Proposal2Slot,
				Shard:           simulatedProposerSlashing.Proposal2Shard,
				BlockRootHash32: []byte(simulatedProposerSlashing.Proposal2Root),
			},
		})
	}
	if simulatedAttesterSlashing != nil {
		block.Body.AttesterSlashings = append(block.Body.AttesterSlashings, &pb.AttesterSlashing{
			SlashableVote_1: &pb.SlashableVote{
				Data: &pb.AttestationData{
					Slot:          simulatedAttesterSlashing.SlashableVote1Slot,
					JustifiedSlot: simulatedAttesterSlashing.SlashableVote1JustifiedSlot,
				},
				CustodyBitfield:  []byte(simulatedAttesterSlashing.SlashableVote1CustodyBitField),
				ValidatorIndices: simulatedAttesterSlashing.SlashableVote1ValidatorIndices,
			},
			SlashableVote_2: &pb.SlashableVote{
				Data: &pb.AttestationData{
					Slot:          simulatedAttesterSlashing.SlashableVote2Slot,
					JustifiedSlot: simulatedAttesterSlashing.SlashableVote2JustifiedSlot,
				},
				CustodyBitfield:  []byte(simulatedAttesterSlashing.SlashableVote2CustodyBitField),
				ValidatorIndices: simulatedAttesterSlashing.SlashableVote2ValidatorIndices,
			},
		})
	}
	if simulatedExit != nil {
		block.Body.Exits = append(block.Body.Exits, &pb.Exit{
			Slot:           simulatedExit.Slot,
			ValidatorIndex: simulatedExit.ValidatorIndex,
		})
	}
	encodedBlock, err := proto.Marshal(block)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("could not marshal new block: %v", err)
	}
	return block, hashutil.Hash(encodedBlock), nil
}

// Given a number of slots, we create a list of hash onions from an underlying randao reveal. For example,
// if we have N slots, we create a list of [secret, hash(secret), hash(hash(secret)), hash(...(prev N-1 hashes))].
func generateSimulatedRandaoHashOnions(numSlots uint64) [][32]byte {
	// We create a list of randao hash onions for the given number of epochs
	// we run the state transition.
	numEpochs := numSlots % params.BeaconConfig().EpochLength
	hashOnions := [][32]byte{params.BeaconConfig().SimulatedBlockRandao}

	// We make the length of the hash onions list equal to the number of epochs + 10 to be safe.
	for i := uint64(0); i < numEpochs+10; i++ {
		prevHash := hashOnions[i]
		hashOnions = append(hashOnions, hashutil.Hash(prevHash[:]))
	}
	return hashOnions
}

// This function determines the block randao reveal assuming there are no skipped slots,
// given a list of randao hash onions such as [pre-image, 0x01, 0x02, 0x03], for the
// 0th epoch, the block randao reveal will be 0x02 and the proposer commitment 0x03.
// The next epoch, the block randao reveal will be 0x01 and the commitment 0x02,
// so on and so forth until all randao layers are peeled off.
func determineSimulatedBlockRandaoReveal(layersPeeled int, hashOnions [][32]byte) [32]byte {
	if layersPeeled == 0 {
		return hashOnions[len(hashOnions)-2]
	}
	return hashOnions[len(hashOnions)-layersPeeled-2]
}

// Generates initial deposits for creating a beacon state in the simulated
// backend based on the yaml configuration.
func generateInitialSimulatedDeposits(randaoCommit [32]byte) ([]*pb.Deposit, error) {
	genesisTime := params.BeaconConfig().GenesisTime.Unix()
	deposits := make([]*pb.Deposit, params.BeaconConfig().DepositsForChainStart)
	for i := 0; i < len(deposits); i++ {
		depositInput := &pb.DepositInput{
			Pubkey:                 []byte(strconv.Itoa(i)),
			RandaoCommitmentHash32: randaoCommit[:],
		}
		depositData, err := b.EncodeDepositData(
			depositInput,
			params.BeaconConfig().MaxDeposit,
			genesisTime,
		)
		if err != nil {
			return nil, fmt.Errorf("could not encode initial block deposits: %v", err)
		}
		deposits[i] = &pb.Deposit{DepositData: depositData}
	}
	return deposits, nil
}
