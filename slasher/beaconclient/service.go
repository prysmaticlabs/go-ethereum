/*
Package beaconclient defines a service that interacts with a beacon
node via a gRPC client to listen for streamed blocks, attestations, and to
submit proposer/attester slashings to the node in case they are detected.
*/
package beaconclient

import (
	"context"
	"errors"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/event"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logrus.WithField("prefix", "beaconclient")

// Notifier defines a struct which exposes event feeds regarding beacon blocks,
// attestations, and more information received from a beacon node.
type Notifier interface {
	BlockFeed() *event.Feed
	AttestationFeed() *event.Feed
	ClientReadyFeed() *event.Feed
}

// HistoricalFetcher defines a struct which can retrieve historical
// block and indexed attestation data from the beacon chain.
type HistoricalFetcher interface {
	RequestHistoricalAttestations(ctx context.Context, epoch uint64) ([]*ethpb.IndexedAttestation, error)
}

// ChainFetcher defines a struct which can retrieve
// chain information from a beacon node such as the latest chain head.
type ChainFetcher interface {
	ChainHead(ctx context.Context) (*ethpb.ChainHead, error)
}

// Service struct for the beaconclient service of the slasher.
type Service struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	cert                  string
	conn                  *grpc.ClientConn
	provider              string
	beaconClient          ethpb.BeaconChainClient
	nodeClient            ethpb.NodeClient
	clientFeed            *event.Feed
	blockFeed             *event.Feed
	attestationFeed       *event.Feed
	proposerSlashingsChan chan *ethpb.ProposerSlashing
	attesterSlashingsChan chan *ethpb.AttesterSlashing
	attesterSlashingsFeed *event.Feed
	proposerSlashingsFeed *event.Feed
}

// Config options for the beaconclient service.
type Config struct {
	BeaconProvider        string
	BeaconCert            string
	ProposerSlashingsFeed *event.Feed
	AttesterSlashingsFeed *event.Feed
}

// NewBeaconClientService instantiation.
func NewBeaconClientService(ctx context.Context, cfg *Config) *Service {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		cert:                  cfg.BeaconCert,
		ctx:                   ctx,
		cancel:                cancel,
		provider:              cfg.BeaconProvider,
		blockFeed:             new(event.Feed),
		clientFeed:            new(event.Feed),
		attestationFeed:       new(event.Feed),
		proposerSlashingsChan: make(chan *ethpb.ProposerSlashing, 1),
		attesterSlashingsChan: make(chan *ethpb.AttesterSlashing, 1),
		attesterSlashingsFeed: cfg.AttesterSlashingsFeed,
		proposerSlashingsFeed: cfg.ProposerSlashingsFeed,
	}
}

// BlockFeed returns a feed other services in slasher can subscribe to
// blocks received via the beacon node through gRPC.
func (bs *Service) BlockFeed() *event.Feed {
	return bs.blockFeed
}

// AttestationFeed returns a feed other services in slasher can subscribe to
// attestations received via the beacon node through gRPC.
func (bs *Service) AttestationFeed() *event.Feed {
	return bs.attestationFeed
}

// ClientReadyFeed returns a feed other services in slasher can subscribe to
// to indicate when the gRPC connection is ready.
func (bs *Service) ClientReadyFeed() *event.Feed {
	return bs.clientFeed
}

// Stop the beacon client service by closing the gRPC connection.
func (bs *Service) Stop() error {
	bs.cancel()
	log.Info("Stopping service")
	if bs.conn != nil {
		return bs.conn.Close()
	}
	return nil
}

// Status returns an error if there exists a gRPC connection error
// in the service.
func (bs *Service) Status() error {
	if bs.conn == nil {
		return errors.New("no connection to beacon RPC")
	}
	return nil
}

// Start the main runtime of the beaconclient service, initializing
// a gRPC client connection with a beacon node, listening for
// streamed blocks/attestations, and submitting slashing operations
// after they are detected by other services in the slasher.
func (bs *Service) Start() {
	var dialOpt grpc.DialOption
	if bs.cert != "" {
		creds, err := credentials.NewClientTLSFromFile(bs.cert, "")
		if err != nil {
			log.Errorf("Could not get valid credentials: %v", err)
		}
		dialOpt = grpc.WithTransportCredentials(creds)
	} else {
		dialOpt = grpc.WithInsecure()
		log.Warn(
			"You are using an insecure gRPC connection to beacon chain! Please provide a certificate and key to use a secure connection",
		)
	}
	beaconOpts := []grpc.DialOption{
		dialOpt,
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithStreamInterceptor(middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(),
			grpc_prometheus.StreamClientInterceptor,
		)),
		grpc.WithUnaryInterceptor(middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
			grpc_prometheus.UnaryClientInterceptor,
		)),
	}
	conn, err := grpc.DialContext(bs.ctx, bs.provider, beaconOpts...)
	if err != nil {
		log.Fatalf("Could not dial endpoint: %s, %v", bs.provider, err)
	}
	log.Info("Successfully started gRPC connection")
	bs.conn = conn
	bs.beaconClient = ethpb.NewBeaconChainClient(bs.conn)
	bs.nodeClient = ethpb.NewNodeClient(bs.conn)

	// We poll for the sync status of the beacon node until it is fully synced.
	bs.querySyncStatus(bs.ctx)

	// We notify other services in slasher that the beacon client is ready
	// and the connection is active.
	bs.clientFeed.Send(true)

	// We register subscribers for any detected proposer/attester slashings
	// in the slasher services that we can submit to the beacon node
	// as they are found.
	go bs.subscribeDetectedProposerSlashings(bs.ctx, bs.proposerSlashingsChan)
	go bs.subscribeDetectedAttesterSlashings(bs.ctx, bs.attesterSlashingsChan)

	// We listen to a stream of blocks and attestations from the beacon node.
	go bs.receiveBlocks(bs.ctx)
	go bs.receiveAttestations(bs.ctx)
}
