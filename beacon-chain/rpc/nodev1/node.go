package nodev1

import (
	"context"
	"fmt"
	"runtime"

	ptypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"github.com/prysmaticlabs/prysm/shared/version"
	"go.opencensus.io/trace"
)

// GetIdentity retrieves data about the node's network presence.
func (ns *Server) GetIdentity(ctx context.Context, _ *ptypes.Empty) (*ethpb.IdentityResponse, error) {
	ctx, span := trace.StartSpan(ctx, "nodev1.GetIdentity")
	defer span.End()

	peerId := ns.PeerManager.PeerID().Pretty()
	enr := ns.PeerManager.ENR().IdentityScheme()
	var p2pAddresses []string
	for _, address := range ns.PeerManager.Host().Addrs() {
		p2pAddresses = append(p2pAddresses, address.String())
	}
	discoveryAddress, err := ns.PeerManager.DiscoveryAddress()
	if err != nil {
		return nil, errors.Wrap(err, "could not obtain discovery address")
	}
	metadata := &ethpb.Metadata{
		SeqNumber: ns.MetadataProvider.MetadataSeq(),
		Attnets:   ns.MetadataProvider.Metadata().Attnets,
	}

	return &ethpb.IdentityResponse{
		Data: &ethpb.Identity{
			PeerId:             peerId,
			Enr:                enr,
			P2PAddresses:       p2pAddresses,
			DiscoveryAddresses: []string{discoveryAddress.String()},
			Metadata:           metadata,
		},
	}, nil
}

// GetPeer retrieves data about the given peer.
func (ns *Server) GetPeer(ctx context.Context, req *ethpb.PeerRequest) (*ethpb.PeerResponse, error) {
	return nil, errors.New("unimplemented")
}

// ListPeers retrieves data about the node's network peers.
func (ns *Server) ListPeers(ctx context.Context, _ *ptypes.Empty) (*ethpb.PeersResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetVersion requests that the beacon node identify information about its implementation in a
// format similar to a HTTP User-Agent field.
func (ns *Server) GetVersion(_ context.Context, _ *ptypes.Empty) (*ethpb.VersionResponse, error) {
	v := fmt.Sprintf("Prysm/%s (%s %s)", version.GetSemanticVersion(), runtime.GOOS, runtime.GOARCH)
	return &ethpb.VersionResponse{
		Data: &ethpb.Version{
			Version: v,
		},
	}, nil
}

// GetSyncStatus requests the beacon node to describe if it's currently syncing or not, and
// if it is, what block it is up to.
func (ns *Server) GetSyncStatus(ctx context.Context, _ *ptypes.Empty) (*ethpb.SyncingResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetHealth returns node health status in http status codes. Useful for load balancers.
// Response Usage:
//    "200":
//      description: Node is ready
//    "206":
//      description: Node is syncing but can serve incomplete data
//    "503":
//      description: Node not initialized or having issues
func (ns *Server) GetHealth(ctx context.Context, _ *ptypes.Empty) (*ptypes.Empty, error) {
	return nil, errors.New("unimplemented")
}
