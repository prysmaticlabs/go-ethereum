package sync

import (
	"context"
	"time"

	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-ssz"
	eth "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
)

// sendRecentBeaconBlocksRequest sends a recent beacon blocks request to a peer to get
// those corresponding blocks from that peer.
func (r *RegularSync) sendRecentBeaconBlocksRequest(ctx context.Context, blockRoots [][32]byte, id peer.ID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stream, err := r.p2p.Send(ctx, blockRoots, id)
	if err != nil {
		return err
	}
	for i := 0; i < len(blockRoots); i++ {
		blk, err := ReadChunkedBlock(stream, r.p2p)
		if err != nil {
			log.WithError(err).Error("Unable to retrieve block from stream")
			return err
		}
		if err := r.chain.ReceiveBlock(ctx, blk); err != nil {
			log.WithError(err).Error("Unable to process block")
			return nil

	code, errMsg, err := ReadStatusCode(stream, r.p2p.Encoding())
	if err != nil {
		return err
	}

	if code != 0 {
		return errors.New(errMsg)
	}

	resp := make([]*eth.BeaconBlock, 0)
	if err := r.p2p.Encoding().DecodeWithLength(stream, &resp); err != nil {
		return err
	}

	r.pendingQueueLock.Lock()
	defer r.pendingQueueLock.Unlock()
	for _, blk := range resp {
		r.slotToPendingBlocks[blk.Slot] = blk
		blkRoot, err := ssz.SigningRoot(blk)
		if err != nil {
			return err
		}
		r.seenPendingBlocks[blkRoot] = true
	}
	return nil
}

// beaconBlocksRootRPCHandler looks up the request blocks from the database from the given block roots.
func (r *RegularSync) beaconBlocksRootRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	defer stream.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	setRPCStreamDeadlines(stream)
	log := log.WithField("handler", "beacon_blocks_by_root")

	blockRoots := msg.([][32]byte)
	if len(blockRoots) == 0 {
		resp, err := r.generateErrorResponse(responseCodeInvalidRequest, "no block roots provided in request")
		if err != nil {
			log.WithError(err).Error("Failed to generate a response error")
		} else {
			if _, err := stream.Write(resp); err != nil {
				log.WithError(err).Errorf("Failed to write to stream")
			}
		}
		return errors.New("no block roots provided")
	}

	for _, root := range blockRoots {
		blk, err := r.db.Block(ctx, root)
		if err != nil {
			log.WithError(err).Error("Failed to fetch block")
			resp, err := r.generateErrorResponse(responseCodeServerError, genericError)
			if err != nil {
				log.WithError(err).Error("Failed to generate a response error")
			} else {
				if _, err := stream.Write(resp); err != nil {
					log.WithError(err).Errorf("Failed to write to stream")
				}
			}
			return err
		}
		if err := r.chunkWriter(stream, blk); err != nil {
			return err
		}
	}
	return nil
}
