package slasher

import (
	"context"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	slashpb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsSlashableAttestation returns an attester slashing if an input
// attestation is found to be slashable.
func (s *Server) IsSlashableAttestation(
	ctx context.Context, req *ethpb.IndexedAttestation,
) (*ethpb.AttesterSlashing, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

// HighestAttestations returns the highest source and target epochs attested for
// validator indices that have been observed by slasher.
func (s *Server) HighestAttestations(
	ctx context.Context, req *slashpb.HighestAttestationRequest,
) (*slashpb.HighestAttestationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}