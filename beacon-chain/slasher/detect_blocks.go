package slasher

import (
	"context"

	types "github.com/prysmaticlabs/eth2-types"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
	"go.opencensus.io/trace"
)

// Given a list of blocks, check if they are slashable for the validators involved.
func (s *Service) detectSlashableBlocks(
	ctx context.Context,
	proposedBlocks []*slashertypes.SignedBlockHeaderWrapper,
) error {
	ctx, span := trace.StartSpan(ctx, "Slasher.detectSlashableBlocks")
	defer span.End()
	// We check if there are any slashable double proposals in the input list
	// of proposals with respect to each other.
	existingProposals := make(map[string]*slashertypes.SignedBlockHeaderWrapper)
	for i, proposal := range proposedBlocks {
		key := proposalKey(proposal)
		existingProposal, ok := existingProposals[key]
		if !ok {
			existingProposals[key] = proposal
			continue
		}
		if isDoubleProposal(proposedBlocks[i].SigningRoot, existingProposal.SigningRoot) {
			logDoubleProposal(proposedBlocks[i], existingProposal)
		}
	}
	// We check if there are any slashable double proposals in the input list
	// of proposals with respect to our database.
	return s.checkDoubleProposalsOnDisk(ctx, proposedBlocks)
}

// Check for double proposals in our database given a list of incoming block proposals.
// For the proposals that were not slashable, we save them to the database.
func (s *Service) checkDoubleProposalsOnDisk(
	ctx context.Context, proposedBlocks []*slashertypes.SignedBlockHeaderWrapper,
) error {
	ctx, span := trace.StartSpan(ctx, "Slasher.checkDoubleProposalsOnDisk")
	defer span.End()
	doubleProposals, err := s.serviceCfg.Database.CheckDoubleBlockProposals(ctx, proposedBlocks)
	if err != nil {
		return err
	}
	// We initialize a map of proposers that are safe from slashing.
	safeProposers := make(map[types.ValidatorIndex]*slashertypes.SignedBlockHeaderWrapper, len(proposedBlocks))
	for _, proposal := range proposedBlocks {
		safeProposers[proposal.SignedBeaconBlockHeader.Header.ProposerIndex] = proposal
	}
	for i, doubleProposal := range doubleProposals {
		logDoubleProposal(proposedBlocks[i], doubleProposal.PrevBeaconBlockWrapper)
		// If a proposer is found to have committed a slashable offense, we delete
		// them from the safe proposers map.
		delete(safeProposers, doubleProposal.ValidatorIndex)
	}
	// We save all the proposals that are determined "safe" and not-slashable to our database.
	safeProposals := make([]*slashertypes.SignedBlockHeaderWrapper, 0, len(safeProposers))
	for _, proposal := range safeProposers {
		safeProposals = append(safeProposals, proposal)
	}
	return s.serviceCfg.Database.SaveBlockProposals(ctx, safeProposals)
}

func proposalKey(proposal *slashertypes.SignedBlockHeaderWrapper) string {
	header := proposal.SignedBeaconBlockHeader.Header
	return uintToString(uint64(header.Slot)) + ":" + uintToString(uint64(header.ProposerIndex))
}
