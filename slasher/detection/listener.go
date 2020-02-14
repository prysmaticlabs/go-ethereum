package detection

import (
	"context"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"go.opencensus.io/trace"
)

// detectIncomingBlocks subscribes to an event feed for
// block objects from a notifier interface. Upon receiving
// a signed beacon block from the feed, we run proposer slashing
// detection on the block.
func (ds *Service) detectIncomingBlocks(ctx context.Context, ch chan *ethpb.SignedBeaconBlock) {
	ctx, span := trace.StartSpan(ctx, "detection.detectIncomingBlocks")
	defer span.End()
	sub := ds.notifier.BlockFeed().Subscribe(ch)
	defer sub.Unsubscribe()
	for {
		select {
		case <-ch:
			log.Infof("Running detection on block...")
			// TODO(Raul): Run detection function for proposer slashings.
		case <-sub.Err():
			log.Error("Subscriber closed, exiting goroutine")
			return
		case <-ctx.Done():
			log.Error("Context canceled")
			return
		}
	}
}

// detectIncomingAttestations subscribes to an event feed for
// attestation objects from a notifier interface. Upon receiving
// an attestation from the feed, we run surround vote and double vote
// detection on the attestation.
func (ds *Service) detectIncomingAttestations(ctx context.Context, ch chan *ethpb.Attestation) {
	ctx, span := trace.StartSpan(ctx, "detection.detectIncomingAttestations")
	defer span.End()
	sub := ds.notifier.AttestationFeed().Subscribe(ch)
	defer sub.Unsubscribe()
	for {
		select {
		case <-ch:
			log.Infof("Running detection on attestation...")
			// TODO(Raul): Run detection function for attester double voting.
			// TODO(Raul): Run detection function for attester surround voting.
		case <-sub.Err():
			log.Error("Subscriber closed, exiting goroutine")
			return
		case <-ctx.Done():
			log.Error("Context canceled")
			return
		}
	}
}
