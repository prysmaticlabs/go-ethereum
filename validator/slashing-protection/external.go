package slashingprotection

import (
	"context"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	log "github.com/sirupsen/logrus"
)

func (s *Service) VerifyBlock(ctx context.Context, blockHeader *ethpb.SignedBeaconBlockHeader) bool {
	ps, err := s.slasherClient.IsSlashableBlock(ctx, blockHeader)
	if err != nil {
		log.Warnf("External slashing block protection returned an error: %v", err)
	}
	if ps != nil {
		return false
	}
	return true
}

func (s *Service) VerifyAttestation(ctx context.Context, attestation *ethpb.IndexedAttestation) bool {
	as, err := s.slasherClient.IsSlashableAttestation(ctx, attestation)
	if err != nil {
		log.Warnf("External slashing attestation protection returned an error: %v", err)
	}
	if as != nil {
		return false
	}
	return true
}
