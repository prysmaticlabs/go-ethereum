package p2p

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	"github.com/karlseguin/ccache"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p/encoder"
	"github.com/prysmaticlabs/prysm/shared"
)

var _ = shared.Service(&Service{})

var pollingPeriod = 1 * time.Second
var bootnodePingPeriod = 5 * time.Minute
var ttl = 1 * time.Hour

const prysmProtocolPrefix = "/prysm/0.0.0"

// Service for managing peer to peer (p2p) networking.
type Service struct {
	ctx           context.Context
	cancel        context.CancelFunc
	started       bool
	cfg           *Config
	startupErr    error
	dv5Listener   Listener
	host          host.Host
	pubsub        *pubsub.PubSub
	exclusionList *ccache.Cache
	privKey       *ecdsa.PrivateKey
}

// NewService initializes a new p2p service compatible with shared.Service interface. No
// connections are made until the Start function is called during the service registry startup.
func NewService(cfg *Config) (*Service, error) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())

	s := &Service{
		ctx:           ctx,
		cancel:        cancel,
		cfg:           cfg,
		exclusionList: ccache.New(ccache.Configure()),
	}

	ipAddr := ipAddr(s.cfg)
	s.privKey, err = privKey(s.cfg)
	if err != nil {
		log.WithError(err).Error("Failed to generate p2p private key")
		return nil, err
	}

	// TODO(3147): Add host options
	opts := buildOptions(s.cfg, ipAddr, s.privKey)
	h, err := libp2p.New(s.ctx, opts...)
	if err != nil {
		log.WithError(err).Error("Failed to create p2p host")
		return nil, err
	}

	dopts := []dhtopts.Option{
		dhtopts.Datastore(dsync.MutexWrap(ds.NewMapDatastore())),
		dhtopts.Protocols(
			protocol.ID(prysmProtocolPrefix + "/dht"),
		),
	}

	dht, err := kaddht.New(ctx, h, dopts...)
	if err != nil {
		return nil, err
	}
	// Wrap host with a routed host so that peers can be looked up in the
	// distributed hash table by their peer ID.
	h = rhost.Wrap(h, dht)
	s.host = h

	// TODO(3147): Add gossip sub options
	// Gossipsub registration is done before we add in any new peers
	// due to libp2p's gossipsub implementation not taking into
	// account previously added peers when creating the gossipsub
	// object.
	psOpts := []pubsub.Option{
		pubsub.WithMessageSigning(false),
		pubsub.WithStrictSignatureVerification(false),
	}
	gs, err := pubsub.NewGossipSub(s.ctx, s.host, psOpts...)
	if err != nil {
		log.WithError(err).Error("Failed to start pubsub")
		return nil, err
	}
	s.pubsub = gs

	return s, nil
}

// Start the p2p service.
func (s *Service) Start() {
	if s.started {
		log.Error("Attempted to start p2p service when it was already started")
		return
	}

	if s.cfg.BootstrapNodeAddr != "" && !s.cfg.NoDiscovery {
		ipAddr := ipAddr(s.cfg)
		listener, err := startDiscoveryV5(ipAddr, s.privKey, s.cfg)
		if err != nil {
			log.WithError(err).Error("Failed to start discovery")
			s.startupErr = err
			return
		}
		err = s.addBootNodeToExclusionList()
		if err != nil {
			log.WithError(err).Error("Could not add bootnode to the exclusion list")
			s.startupErr = err
			return
		}
		s.dv5Listener = listener

		go s.listenForNewNodes()
		go s.maintainBootnode()
	}

	s.started = true

	if len(s.cfg.StaticPeers) > 0 {
		addrs, err := manyMultiAddrsFromString(s.cfg.StaticPeers)
		if err != nil {
			log.Errorf("Could not connect to static peer: %v", err)
		}
		s.connectWithAllPeers(addrs)
	}

	registerMetrics(s)
	multiAddrs := s.host.Network().ListenAddresses()
	logIP4Addr(s.host.ID(), multiAddrs...)
}

// Stop the p2p service and terminate all peer connections.
func (s *Service) Stop() error {
	s.started = false
	if s.dv5Listener != nil {
		s.dv5Listener.Close()
	}
	return nil
}

// Status of the p2p service. Will return an error if the service is considered unhealthy to
// indicate that this node should not serve traffic until the issue has been resolved.
func (s *Service) Status() error {
	if !s.started {
		return errors.New("not running")
	}
	return nil
}

// Started returns true if the p2p service has successfully started.
func (s *Service) Started() bool {
	return s.started
}

// Encoding returns the configured networking encoding.
func (s *Service) Encoding() encoder.NetworkEncoding {
	encoding := s.cfg.Encoding
	switch encoding {
	case encoder.SSZ:
		return &encoder.SszNetworkEncoder{}
	case encoder.SSZSnappy:
		return &encoder.SszNetworkEncoder{UseSnappyCompression: true}
	default:
		panic("Invalid Network Encoding Flag Provided")
	}
}

// PubSub returns the p2p pubsub framework.
func (s *Service) PubSub() *pubsub.PubSub {
	return s.pubsub
}

// SetStreamHandler sets the protocol handler on the p2p host multiplexer.
// This method is a pass through to libp2pcore.Host.SetStreamHandler.
func (s *Service) SetStreamHandler(topic string, handler network.StreamHandler) {
	s.host.SetStreamHandler(protocol.ID(topic), handler)
}

// PeerID returns the Peer ID of the local peer.
func (s *Service) PeerID() peer.ID {
	return s.host.ID()
}

// Disconnect from a peer.
func (s *Service) Disconnect(pid peer.ID) error {
	return s.host.Network().ClosePeer(pid)
}

// listen for new nodes watches for new nodes in the network and adds them to the peerstore.
func (s *Service) listenForNewNodes() {
	ticker := time.NewTicker(pollingPeriod)
	bootNode, err := enode.Parse(enode.ValidSchemes, s.cfg.BootstrapNodeAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ticker.C:
			nodes := s.dv5Listener.Lookup(bootNode.ID())
			multiAddresses := convertToMultiAddr(nodes)
			s.connectWithAllPeers(multiAddresses)
		case <-s.ctx.Done():
			log.Debug("p2p context is closed, exiting routine")
			break
		}
	}
}

// maintainBootnode connection by infrequently pinging the bootnode ENR. If the bootnode server has
// restarted and this client pruned them from the local table, then a ping will reinsert this
// clients' ENR into the table of the bootnode again.
func (s *Service) maintainBootnode() {
	if s.cfg.BootstrapNodeAddr == "" {
		return
	}
	ticker := time.NewTicker(bootnodePingPeriod)
	bootNode := enode.MustParse(s.cfg.BootstrapNodeAddr)
	for {
		select {
		case <-ticker.C:
			log.Debug("Pinging bootnode")
			if err := s.dv5Listener.Ping(bootNode); err != nil {
				log.WithError(err).Error("Failed to ping bootnode")
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Service) connectWithAllPeers(multiAddrs []ma.Multiaddr) {
	addrInfos, err := peer.AddrInfosFromP2pAddrs(multiAddrs...)
	if err != nil {
		log.Errorf("Could not convert to peer address info's from multiaddresses: %v", err)
		return
	}
	for _, info := range addrInfos {
		if info.ID == s.host.ID() {
			continue
		}
		if s.exclusionList.Get(info.ID.String()) != nil {
			continue
		}
		if err := s.host.Connect(s.ctx, info); err != nil {
			log.Errorf("Could not connect with peer %s: %v", info.String(), err)
			s.exclusionList.Set(info.ID.String(), true, ttl)
		}
	}
}

func (s *Service) addBootNodeToExclusionList() error {
	bootNode, err := enode.Parse(enode.ValidSchemes, s.cfg.BootstrapNodeAddr)
	if err != nil {
		return err
	}
	multAddr, err := convertToSingleMultiAddr(bootNode)
	if err != nil {
		return err
	}
	addrInfo, err := peer.AddrInfoFromP2pAddr(multAddr)
	if err != nil {
		return err
	}
	// bootnode is never dialled, so ttl is tentatively 1 year
	s.exclusionList.Set(addrInfo.ID.String(), true, 365*24*time.Hour)

	return nil
}

func logIP4Addr(id peer.ID, addrs ...ma.Multiaddr) {
	var correctAddr ma.Multiaddr
	for _, addr := range addrs {
		if strings.Contains(addr.String(), "/ip4/") {
			correctAddr = addr
			break
		}
	}
	log.Infof("Node's listening multiaddr is %s", correctAddr.String()+"/p2p/"+id.String())
}
