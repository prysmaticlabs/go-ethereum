package slasher

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
)

// IsSlashableBlock checks if an input block header is slashable
// with respect to historical block proposal data.
func (s *Service) IsSlashableBlock(
	ctx context.Context, block *ethpb.SignedBeaconBlockHeader,
) (*ethpb.ProposerSlashing, error) {
	dataRoot, err := block.Header.HashTreeRoot()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get block header hash tree root: %v", err)
	}
	signedBlockWrapper := &slashertypes.SignedBlockHeaderWrapper{
		SignedBeaconBlockHeader: block,
		SigningRoot:             dataRoot,
	}
	proposerSlashings, err := s.detectProposerSlashings(ctx, []*slashertypes.SignedBlockHeaderWrapper{signedBlockWrapper})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not check if proposal is slashable: %v", err)
	}
	if len(proposerSlashings) == 0 {
		return nil, nil
	}
	return proposerSlashings[0], nil
}

// IsSlashableAttestation checks if an input indexed attestation is slashable
// with respect to historical attestation data.
func (s *Service) IsSlashableAttestation(
	ctx context.Context, attestation *ethpb.IndexedAttestation,
) ([]*ethpb.AttesterSlashing, error) {
	dataRoot, err := attestation.Data.HashTreeRoot()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get attestation data hash tree root: %v", err)
	}
	indexedAttWrapper := &slashertypes.IndexedAttestationWrapper{
		IndexedAttestation: attestation,
		SigningRoot:        dataRoot,
	}

	// Save the attestation record to our database.
	if err := s.serviceCfg.Database.SaveAttestationRecordsForValidators(
		ctx, []*slashertypes.IndexedAttestationWrapper{indexedAttWrapper}, s.params.historyLength,
	); err != nil {
		return nil, status.Errorf(codes.Internal, "Could not save attestation records to DB: %v", err)
	}

	attesterSlashings, err := s.checkSlashableAttestations(ctx, []*slashertypes.IndexedAttestationWrapper{indexedAttWrapper})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not check if attestation is slashable: %v", err)
	}
	if len(attesterSlashings) == 0 {
		return nil, nil
	}
	return attesterSlashings, nil
}
