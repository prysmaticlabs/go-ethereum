package slasher

import (
	"context"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/params"
)

type MockSlashingChecker struct {
	AttesterSlashingFound bool
	ProposerSlashingFound bool
}

func (s *MockSlashingChecker) IsSlashableBlock(ctx context.Context, proposal *ethpb.SignedBeaconBlockHeader) (*ethpb.ProposerSlashing, error) {
	if s.ProposerSlashingFound {
		return &ethpb.ProposerSlashing{
			Header_1: &ethpb.SignedBeaconBlockHeader{
				Header: &ethpb.BeaconBlockHeader{
					Slot:          0,
					ProposerIndex: 0,
					ParentRoot:    params.BeaconConfig().ZeroHash[:],
					StateRoot:     params.BeaconConfig().ZeroHash[:],
					BodyRoot:      params.BeaconConfig().ZeroHash[:],
				},
				Signature: params.BeaconConfig().EmptySignature[:],
			},
			Header_2: &ethpb.SignedBeaconBlockHeader{
				Header: &ethpb.BeaconBlockHeader{
					Slot:          0,
					ProposerIndex: 0,
					ParentRoot:    params.BeaconConfig().ZeroHash[:],
					StateRoot:     params.BeaconConfig().ZeroHash[:],
					BodyRoot:      params.BeaconConfig().ZeroHash[:],
				},
				Signature: params.BeaconConfig().EmptySignature[:],
			},
		}, nil
	}
	return nil, nil
}

func (s *MockSlashingChecker) IsSlashableAttestation(ctx context.Context, attestation *ethpb.IndexedAttestation) ([]*ethpb.AttesterSlashing, error) {
	if s.AttesterSlashingFound {
		return []*ethpb.AttesterSlashing{
			{
				Attestation_1: &ethpb.IndexedAttestation{
					Data: &ethpb.AttestationData{},
				},
				Attestation_2: &ethpb.IndexedAttestation{
					Data: &ethpb.AttestationData{},
				},
			},
		}, nil
	}
	return nil, nil
}
