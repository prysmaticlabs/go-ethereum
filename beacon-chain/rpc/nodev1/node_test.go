package nodev1

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	ptypes "github.com/gogo/protobuf/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	mockp2p "github.com/prysmaticlabs/prysm/beacon-chain/p2p/testing"
	syncmock "github.com/prysmaticlabs/prysm/beacon-chain/sync/initial-sync/testing"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/shared/version"
)

type dummyIdentity enode.ID

func (id dummyIdentity) Verify(_ *enr.Record, _ []byte) error { return nil }
func (id dummyIdentity) NodeAddr(_ *enr.Record) []byte        { return id[:] }

func TestGetVersion(t *testing.T) {
	semVer := version.GetSemanticVersion()
	os := runtime.GOOS
	arch := runtime.GOARCH
	res, err := (&Server{}).GetVersion(context.Background(), &ptypes.Empty{})
	require.NoError(t, err)
	v := res.Data.Version
	assert.Equal(t, true, strings.Contains(v, semVer))
	assert.Equal(t, true, strings.Contains(v, os))
	assert.Equal(t, true, strings.Contains(v, arch))
}

func TestGetHealth(t *testing.T) {
	ctx := context.Background()
	checker := &syncmock.Sync{}
	s := &Server{
		SyncChecker: checker,
	}

	_, err := s.GetHealth(ctx, &ptypes.Empty{})
	require.ErrorContains(t, "node not initialized or having issues", err)
	checker.IsInitialized = true
	_, err = s.GetHealth(ctx, &ptypes.Empty{})
	require.NoError(t, err)
	checker.IsInitialized = false
	checker.IsSyncing = true
	require.NoError(t, err)
}

func TestGetIdentity(t *testing.T) {
	ctx := context.Background()
	p2pAddr, err := ma.NewMultiaddr("/ip4/7.7.7.7/udp/30303")
	require.NoError(t, err)
	discAddr1, err := ma.NewMultiaddr("/ip4/7.7.7.7/udp/30303/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	require.NoError(t, err)
	discAddr2, err := ma.NewMultiaddr("/ip6/1:2:3:4:5:6:7:8/udp/20202/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	require.NoError(t, err)
	enrRecord := &enr.Record{}
	err = enrRecord.SetSig(dummyIdentity{1}, []byte{42})
	require.NoError(t, err)
	enrRecord.Set(enr.IPv4{7, 7, 7, 7})
	err = enrRecord.SetSig(dummyIdentity{}, []byte{})
	require.NoError(t, err)
	attnets := bitfield.NewBitvector64()
	attnets.SetBitAt(1, true)
	metadataProvider := &mockp2p.MockMetadataProvider{Data: &pb.MetaData{SeqNumber: 1, Attnets: attnets}}

	t.Run("OK", func(t *testing.T) {
		peerManager := &mockp2p.MockPeerManager{
			Enr:           enrRecord,
			PID:           "foo",
			BHost:         &mockp2p.MockHost{Addresses: []ma.Multiaddr{p2pAddr}},
			DiscoveryAddr: []multiaddr.Multiaddr{discAddr1, discAddr2},
		}
		s := &Server{
			PeerManager:      peerManager,
			MetadataProvider: metadataProvider,
		}

		resp, err := s.GetIdentity(ctx, &ptypes.Empty{})
		require.NoError(t, err)
		expectedID := peer.ID("foo").Pretty()
		assert.Equal(t, expectedID, resp.Data.PeerId)
		expectedEnr, err := p2p.SerializeENR(enrRecord)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprint("enr:", expectedEnr), resp.Data.Enr)
		require.Equal(t, 1, len(resp.Data.P2PAddresses))
		assert.Equal(t, fmt.Sprint(p2pAddr.String(), "/p2p/", expectedID), resp.Data.P2PAddresses[0])
		require.Equal(t, 2, len(resp.Data.DiscoveryAddresses))
		ipv4Found, ipv6Found := false, false
		for _, address := range resp.Data.DiscoveryAddresses {
			if address == discAddr1.String() {
				ipv4Found = true
			} else if address == discAddr2.String() {
				ipv6Found = true
			}
		}
		assert.Equal(t, true, ipv4Found, "IPv4 discovery address not found")
		assert.Equal(t, true, ipv6Found, "IPv6 discovery address not found")
		assert.Equal(t, discAddr1.String(), resp.Data.DiscoveryAddresses[0])
		assert.Equal(t, discAddr2.String(), resp.Data.DiscoveryAddresses[1])
	})

	t.Run("ENR failure", func(t *testing.T) {
		peerManager := &mockp2p.MockPeerManager{
			Enr:           &enr.Record{},
			PID:           "foo",
			BHost:         &mockp2p.MockHost{Addresses: []ma.Multiaddr{p2pAddr}},
			DiscoveryAddr: []multiaddr.Multiaddr{discAddr1, discAddr2},
		}
		s := &Server{
			PeerManager:      peerManager,
			MetadataProvider: metadataProvider,
		}

		_, err = s.GetIdentity(ctx, &ptypes.Empty{})
		assert.ErrorContains(t, "could not obtain enr", err)
	})

	t.Run("Discovery addresses failure", func(t *testing.T) {
		peerManager := &mockp2p.MockPeerManager{
			Enr:               enrRecord,
			PID:               "foo",
			BHost:             &mockp2p.MockHost{Addresses: []ma.Multiaddr{p2pAddr}},
			DiscoveryAddr:     []multiaddr.Multiaddr{discAddr1, discAddr2},
			FailDiscoveryAddr: true,
		}
		s := &Server{
			PeerManager:      peerManager,
			MetadataProvider: metadataProvider,
		}

		_, err = s.GetIdentity(ctx, &ptypes.Empty{})
		assert.ErrorContains(t, "could not obtain discovery address", err)
	})
}
