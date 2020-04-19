package beacon

import (
	"context"
	"sort"
	"strconv"
	"time"

	stateTrie "github.com/prysmaticlabs/prysm/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/shared/featureconfig"

	ptypes "github.com/gogo/protobuf/types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/feed"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/feed/operation"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/db/filters"
	"github.com/prysmaticlabs/prysm/beacon-chain/flags"
	"github.com/prysmaticlabs/prysm/shared/attestationutil"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/pagination"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sortableAttestations implements the Sort interface to sort attestations
// by slot as the canonical sorting attribute.
type sortableAttestations []*ethpb.Attestation

func (s sortableAttestations) Len() int      { return len(s) }
func (s sortableAttestations) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortableAttestations) Less(i, j int) bool {
	return s[i].Data.Slot < s[j].Data.Slot
}

func mapAttestationToBlockRoot(ctx context.Context, atts []*ethpb.Attestation) map[[32]byte][]*ethpb.Attestation {
	attsMap := make(map[[32]byte][]*ethpb.Attestation)
	if len(atts) == 0 {
		return attsMap
	}
	for _, att := range atts {
		attsMap[bytesutil.ToBytes32(att.Data.BeaconBlockRoot)] = append(attsMap[bytesutil.ToBytes32(att.Data.BeaconBlockRoot)], att)
	}
	return attsMap
}

// ListAttestations retrieves attestations by block root, slot, or epoch.
// Attestations are sorted by data slot by default.
//
// The server may return an empty list when no attestations match the given
// filter criteria. This RPC should not return NOT_FOUND. Only one filter
// criteria should be used.
func (bs *Server) ListAttestations(
	ctx context.Context, req *ethpb.ListAttestationsRequest,
) (*ethpb.ListAttestationsResponse, error) {
	if int(req.PageSize) > flags.Get().MaxPageSize {
		return nil, status.Errorf(codes.InvalidArgument, "Requested page size %d can not be greater than max size %d",
			req.PageSize, flags.Get().MaxPageSize)
	}
	var blocks []*ethpb.SignedBeaconBlock
	var err error
	switch q := req.QueryFilter.(type) {
	case *ethpb.ListAttestationsRequest_GenesisEpoch:
		blocks, err = bs.BeaconDB.Blocks(ctx, filters.NewFilter().SetStartEpoch(0).SetEndEpoch(0))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not fetch attestations: %v", err)
		}
	case *ethpb.ListAttestationsRequest_Epoch:
		blocks, err = bs.BeaconDB.Blocks(ctx, filters.NewFilter().SetStartEpoch(q.Epoch).SetEndEpoch(q.Epoch))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not fetch attestations: %v", err)
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "Must specify a filter criteria for fetching attestations")
	}
	atts := make([]*ethpb.Attestation, 0, params.BeaconConfig().MaxAttestations*uint64(len(blocks)))
	for _, block := range blocks {
		atts = append(atts, block.Block.Body.Attestations...)
	}
	// We sort attestations according to the Sortable interface.
	sort.Sort(sortableAttestations(atts))
	numAttestations := len(atts)

	// If there are no attestations, we simply return a response specifying this.
	// Otherwise, attempting to paginate 0 attestations below would result in an error.
	if numAttestations == 0 {
		return &ethpb.ListAttestationsResponse{
			Attestations:  make([]*ethpb.Attestation, 0),
			TotalSize:     int32(0),
			NextPageToken: strconv.Itoa(0),
		}, nil
	}

	start, end, nextPageToken, err := pagination.StartAndEndPage(req.PageToken, int(req.PageSize), numAttestations)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not paginate attestations: %v", err)
	}
	return &ethpb.ListAttestationsResponse{
		Attestations:  atts[start:end],
		TotalSize:     int32(numAttestations),
		NextPageToken: nextPageToken,
	}, nil
}

// ListIndexedAttestations retrieves indexed attestations by block root.
// IndexedAttestationsForEpoch are sorted by data slot by default. Either a start-end epoch
// filter is used to retrieve blocks with.
//
// The server may return an empty list when no attestations match the given
// filter criteria. This RPC should not return NOT_FOUND.
func (bs *Server) ListIndexedAttestations(
	ctx context.Context, req *ethpb.ListIndexedAttestationsRequest,
) (*ethpb.ListIndexedAttestationsResponse, error) {
	blocks := make([]*ethpb.SignedBeaconBlock, 0)
	var err error
	switch q := req.QueryFilter.(type) {
	case *ethpb.ListIndexedAttestationsRequest_GenesisEpoch:
		blocks, err = bs.BeaconDB.Blocks(ctx, filters.NewFilter().SetStartEpoch(0).SetEndEpoch(0))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not fetch attestations: %v", err)
		}
	case *ethpb.ListIndexedAttestationsRequest_Epoch:
		blocks, err = bs.BeaconDB.Blocks(ctx, filters.NewFilter().SetStartEpoch(q.Epoch).SetEndEpoch(q.Epoch))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not fetch attestations: %v", err)
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "Must specify a filter criteria for fetching attestations")
	}
	attsArray := make([]*ethpb.Attestation, 0, params.BeaconConfig().MaxAttestations*uint64(len(blocks)))
	for _, block := range blocks {
		attsArray = append(attsArray, block.Block.Body.Attestations...)
	}
	// We sort attestations according to the Sortable interface.
	sort.Sort(sortableAttestations(attsArray))
	numAttestations := len(attsArray)

	// If there are no attestations, we simply return a response specifying this.
	// Otherwise, attempting to paginate 0 attestations below would result in an error.
	if numAttestations == 0 {
		return &ethpb.ListIndexedAttestationsResponse{
			IndexedAttestations: make([]*ethpb.IndexedAttestation, 0),
			TotalSize:           int32(0),
			NextPageToken:       strconv.Itoa(0),
		}, nil
	}
	// We use the retrieved committees for the block root to convert all attestations
	// into indexed form effectively.
	mappedAttestations := mapAttestationToBlockRoot(ctx, attsArray)
	indexedAtts := make([]*ethpb.IndexedAttestation, numAttestations, numAttestations)
	for atts := range mappedAttestations {
		var attState *stateTrie.BeaconState
		if !featureconfig.Get().DisableNewStateMgmt {
			attState, err = bs.StateGen.StateByRoot(ctx, bytesutil.ToBytes32(att.Data.BeaconBlockRoot))
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					"Could not retrieve state for attestation data block root %v: %v",
					att.Data.BeaconBlockRoot,
					err,
				)
			}
		} else {
			attState, err = bs.BeaconDB.State(ctx, bytesutil.ToBytes32(att.Data.BeaconBlockRoot))
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					"Could not retrieve state for attestation data block root %v: %v",
					att.Data.BeaconBlockRoot,
					err,
				)
			}
		}
		for i := 0; i < len(atts); i++ {
			att := atts[i]

			committee, err := helpers.BeaconCommitteeFromState(attState, att.Data.Slot, att.Data.CommitteeIndex)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					"Could not retrieve committees from state %v",
					err,
				)
			}
			idxAtt := attestationutil.ConvertToIndexed(ctx, atts[i], committee)
			indexedAtts[i] = idxAtt
		}
	}

	start, end, nextPageToken, err := pagination.StartAndEndPage(req.PageToken, int(req.PageSize), len(indexedAtts))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not paginate attestations: %v", err)
	}
	return &ethpb.ListIndexedAttestationsResponse{
		IndexedAttestations: indexedAtts[start:end],
		TotalSize:           int32(len(indexedAtts)),
		NextPageToken:       nextPageToken,
	}, nil
}

// StreamAttestations to clients at the end of every slot. This method retrieves the
// aggregated attestations currently in the pool at the start of a slot and sends
// them over a gRPC stream.
func (bs *Server) StreamAttestations(
	_ *ptypes.Empty, stream ethpb.BeaconChain_StreamAttestationsServer,
) error {
	attestationsChannel := make(chan *feed.Event, 1)
	attSub := bs.AttestationNotifier.OperationFeed().Subscribe(attestationsChannel)
	defer attSub.Unsubscribe()
	for {
		select {
		case event := <-attestationsChannel:
			if event.Type == operation.UnaggregatedAttReceived {
				data, ok := event.Data.(*operation.UnAggregatedAttReceivedData)
				if !ok {
					// Got bad data over the stream.
					continue
				}
				if data.Attestation == nil {
					// One nil attestation shouldn't stop the stream.
					continue
				}
				if err := stream.Send(data.Attestation); err != nil {
					return status.Errorf(codes.Unavailable, "Could not send over stream: %v", err)
				}
			}
		case <-bs.Ctx.Done():
			return status.Error(codes.Canceled, "Context canceled")
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Context canceled")
		}
	}
}

// StreamIndexedAttestations to clients at the end of every slot. This method retrieves the
// aggregated attestations currently in the pool, converts them into indexed form, and
// sends them over a gRPC stream.
func (bs *Server) StreamIndexedAttestations(
	_ *ptypes.Empty, stream ethpb.BeaconChain_StreamIndexedAttestationsServer,
) error {
	attestationsChannel := make(chan *feed.Event, 1)
	attSub := bs.AttestationNotifier.OperationFeed().Subscribe(attestationsChannel)
	defer attSub.Unsubscribe()
	go bs.collectReceivedAttestations(stream.Context())
	for {
		select {
		case event := <-attestationsChannel:
			if event.Type == operation.UnaggregatedAttReceived {
				data, ok := event.Data.(*operation.UnAggregatedAttReceivedData)
				if !ok {
					// Got bad data over the stream.
					continue
				}
				if data.Attestation == nil {
					// One nil attestation shouldn't stop the stream.
					continue
				}
				bs.ReceivedAttestationsBuffer <- data.Attestation
			}
		case atts := <-bs.CollectedAttestationsBuffer:
			// We aggregate the received attestations.
			aggAtts, err := helpers.AggregateAttestations(atts)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					"Could not aggregate attestations: %v",
					err,
				)
			}
			if len(aggAtts) == 0 {
				continue
			}
			// All attestations we receive have the same target epoch given they
			// have the same data root, so we just use the target epoch from
			// the first one to determine committees for converting into indexed
			// form.
			epoch := aggAtts[0].Data.Target.Epoch
			committeesBySlot, _, err := bs.retrieveCommitteesForEpoch(stream.Context(), epoch)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					"Could not retrieve committees for epoch %d: %v",
					epoch,
					err,
				)
			}
			// We use the retrieved committees for the epoch to convert all attestations
			// into indexed form effectively.
			startSlot := helpers.StartSlot(epoch)
			endSlot := startSlot + params.BeaconConfig().SlotsPerEpoch
			for _, att := range aggAtts {
				// Out of range check, the attestation slot cannot be greater
				// the last slot of the requested epoch or smaller than its start slot
				// given committees are accessed as a map of slot -> commitees list, where there are
				// SLOTS_PER_EPOCH keys in the map.
				if att.Data.Slot < startSlot || att.Data.Slot > endSlot {
					continue
				}
				committeesForSlot, ok := committeesBySlot[att.Data.Slot]
				if !ok || committeesForSlot.Committees == nil {
					continue
				}
				committee := committeesForSlot.Committees[att.Data.CommitteeIndex]
				idxAtt := attestationutil.ConvertToIndexed(stream.Context(), att, committee.ValidatorIndices)
				if err := stream.Send(idxAtt); err != nil {
					return status.Errorf(codes.Unavailable, "Could not send over stream: %v", err)
				}
			}
		case <-bs.Ctx.Done():
			return status.Error(codes.Canceled, "Context canceled")
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Context canceled")
		}
	}
}

// TODO(#5031): Instead of doing aggregation here, leverage the aggregation
// already being done by the attestation pool in the operations service.
func (bs *Server) collectReceivedAttestations(ctx context.Context) {
	attsByRoot := make(map[[32]byte][]*ethpb.Attestation)
	halfASlot := time.Duration(params.BeaconConfig().SecondsPerSlot / 2)
	ticker := time.NewTicker(time.Second * halfASlot)
	for {
		select {
		case <-ticker.C:
			for root, atts := range attsByRoot {
				if len(atts) > 0 {
					bs.CollectedAttestationsBuffer <- atts
					attsByRoot[root] = make([]*ethpb.Attestation, 0)
				}
			}
		case att := <-bs.ReceivedAttestationsBuffer:
			attDataRoot, err := ssz.HashTreeRoot(att.Data)
			if err != nil {
				logrus.Errorf("Could not hash tree root data: %v", err)
				continue
			}
			attsByRoot[attDataRoot] = append(attsByRoot[attDataRoot], att)
		case <-ctx.Done():
			return
		case <-bs.Ctx.Done():
			return
		}
	}
}

// AttestationPool retrieves pending attestations.
//
// The server returns a list of attestations that have been seen but not
// yet processed. Pool attestations eventually expire as the slot
// advances, so an attestation missing from this request does not imply
// that it was included in a block. The attestation may have expired.
// Refer to the ethereum 2.0 specification for more details on how
// attestations are processed and when they are no longer valid.
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/core/0_beacon-chain.md#attestations
func (bs *Server) AttestationPool(
	ctx context.Context, req *ethpb.AttestationPoolRequest,
) (*ethpb.AttestationPoolResponse, error) {
	if int(req.PageSize) > flags.Get().MaxPageSize {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Requested page size %d can not be greater than max size %d",
			req.PageSize,
			flags.Get().MaxPageSize,
		)
	}
	atts := bs.AttestationsPool.AggregatedAttestations()
	numAtts := len(atts)
	if numAtts == 0 {
		return &ethpb.AttestationPoolResponse{
			Attestations:  make([]*ethpb.Attestation, 0),
			TotalSize:     int32(0),
			NextPageToken: strconv.Itoa(0),
		}, nil
	}
	start, end, nextPageToken, err := pagination.StartAndEndPage(req.PageToken, int(req.PageSize), numAtts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not paginate attestations: %v", err)
	}
	return &ethpb.AttestationPoolResponse{
		Attestations:  atts[start:end],
		TotalSize:     int32(numAtts),
		NextPageToken: nextPageToken,
	}, nil
}
