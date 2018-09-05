// Package rpc defines the services that the beacon-chain uses to communicate via gRPC.
package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/prysmaticlabs/prysm/beacon-chain/types"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logrus.WithField("prefix", "rpc")

// canonicalFetcher defines a struct with methods that can be
// called on-demand to fetch the latest canonical head
// and crystallized state as well as methods that stream
// latest canonical head events to clients
// These functions are called by a validator client upon
// establishing an initial connection to a beacon node via gRPC.
type canonicalFetcher interface {
	// These methods can be called on-demand by a validator
	// to fetch canonical head and state.
	CanonicalHead() (*types.Block, error)
	CanonicalCrystallizedState() *types.CrystallizedState
	// These methods are not called on-demand by a validator
	// but instead streamed to connected validators every
	// time the canonical head changes in the chain service.
	CanonicalBlockFeed() *event.Feed
	CanonicalCrystallizedStateFeed() *event.Feed
}

type chainService interface {
	IncomingBlockFeed() *event.Feed
}

type powChainService interface {
	LatestBlockHash() common.Hash
}

// Service defining an RPC server for a beacon node.
type Service struct {
	ctx                context.Context
	cancel             context.CancelFunc
	fetcher            canonicalFetcher
	chainService       chainService
	powChainService    powChainService
	port               string
	listener           net.Listener
	withCert           string
	withKey            string
	grpcServer         *grpc.Server
	canonicalBlockChan chan *types.Block
	canonicalStateChan chan *types.CrystallizedState
}

// Config options for the beacon node RPC server.
type Config struct {
	Port             string
	CertFlag         string
	KeyFlag          string
	SubscriptionBuf  int
	CanonicalFetcher canonicalFetcher
	ChainService     chainService
	POWChainService  powChainService
}

// NewRPCService creates a new instance of a struct implementing the BeaconServiceServer
// interface.
func NewRPCService(ctx context.Context, cfg *Config) *Service {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		ctx:                ctx,
		cancel:             cancel,
		fetcher:            cfg.CanonicalFetcher,
		chainService:       cfg.ChainService,
		powChainService:    cfg.POWChainService,
		port:               cfg.Port,
		withCert:           cfg.CertFlag,
		withKey:            cfg.KeyFlag,
		canonicalBlockChan: make(chan *types.Block, cfg.SubscriptionBuf),
		canonicalStateChan: make(chan *types.CrystallizedState, cfg.SubscriptionBuf),
	}
}

// Start the gRPC server.
func (s *Service) Start() {
	log.Info("Starting service")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		log.Errorf("Could not listen to port :%s: %v", s.port, err)
		return
	}
	s.listener = lis
	log.Infof("RPC server listening on port :%s", s.port)

	if s.withCert != "" && s.withKey != "" {
		creds, err := credentials.NewServerTLSFromFile(s.withCert, s.withKey)
		if err != nil {
			log.Errorf("could not load TLS keys: %s", err)
		}
		s.grpcServer = grpc.NewServer(grpc.Creds(creds))
	} else {
		log.Warn("You are using an insecure gRPC connection! Please provide a certificate and key to use a secure connection")
		s.grpcServer = grpc.NewServer()
	}

	pb.RegisterBeaconServiceServer(s.grpcServer, s)
	pb.RegisterProposerServiceServer(s.grpcServer, s)
	go func() {
		err = s.grpcServer.Serve(lis)
		if err != nil {
			log.Errorf("Could not serve gRPC: %v", err)
		}
	}()
}

// Stop the service.
func (s *Service) Stop() error {
	log.Info("Stopping service")
	s.cancel()
	if s.listener != nil {
		s.grpcServer.GracefulStop()
		log.Debug("Initiated graceful stop of gRPC server")
	}
	return nil
}

// GenesisTimeAndCanonicalState returns the genesis timestamp and crystallized state
// determined as canonical. Validator clients send this request
// once upon establishing a connection to the beacon node in order to determine
// their role and assigned slot initially and setup an internal ticker.
func (s *Service) GenesisTimeAndCanonicalState(ctx context.Context, req *empty.Empty) (*pb.GenesisTimeAndCanonicalState, error) {
	genesis, err := types.NewGenesisBlock()
	if err != nil {
		return nil, fmt.Errorf("could not fetch genesis block: %v", err)
	}
	crystallized := s.fetcher.CanonicalCrystallizedState()
	return &pb.CanonicalResponse{
		GenesisTimestamp:        genesis.Timestamp(),
		LatestCrystallizedState: crystallized.Proto(),
	}, nil
}

// FetchShuffledValidatorIndices retrieves the shuffled validator indices, cutoffs, and
// assigned attestation slots at a given crystallized state hash.
// This function can be called by validators to fetch a historical list of shuffled
// validators ata point in time corresponding to a certain crystallized state.
func (s *Service) FetchShuffledValidatorIndices(ctx context.Context, req *pb.ShuffleRequest) (*pb.ShuffleResponse, error) {
	var shuffledIndices []uint64
	// Simulator always pushes out a validator list of length 100. By having index 0
	// as the last index, the validator will always be a proposer in the validator code.
	// TODO: Implement the real method by fetching the crystallized state in the request
	// from persistent disk storage and shuffling the indices appropriately.
	for i := 99; i >= 0; i-- {
		shuffledIndices = append(shuffledIndices, uint64(i))
	}
	// For now, this will cause validators to always pick the validator as a proposer.
	shuffleRes := &pb.ShuffleResponse{
		ShuffledValidatorIndices: shuffledIndices,
	}
	return shuffleRes, nil
}

// ProposeBlock is called by a proposer in a sharding validator and a full beacon node
// sends the request into a beacon block that can then be included in a canonical chain.
func (s *Service) ProposeBlock(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeResponse, error) {
	// TODO: handle fields such as attestation bitmask, aggregate sig, and randao reveal.
	powChainHash := s.powChainService.LatestBlockHash()
	data := &pbp2p.BeaconBlock{
		SlotNumber:  req.GetSlotNumber(),
		PowChainRef: powChainHash[:],
		ParentHash:  req.GetParentHash(),
		Timestamp:   req.GetTimestamp(),
	}
	block := types.NewBlock(data)
	h, err := block.Hash()
	if err != nil {
		return nil, fmt.Errorf("could not hash block: %v", err)
	}
	log.WithField("blockHash", fmt.Sprintf("0x%x", h)).Debugf("Block proposal received via RPC")
	// We relay the received block from the proposer to the chain service for processing.
	s.chainService.IncomingBlockFeed().Send(block)
	return &pb.ProposeResponse{BlockHash: h[:]}, nil
}

// SignBlock is a function called by an attester in a sharding validator to sign off
// on a block.
//
// TODO: needs implementation.
func (s *Service) SignBlock(ctx context.Context, req *pb.SignRequest) (*pb.SignResponse, error) {
	// TODO: implement.
	return nil, errors.New("unimplemented")
}

// LatestBeaconBlock streams the latest beacon chain data.
func (s *Service) LatestBeaconBlock(req *empty.Empty, stream pb.BeaconService_LatestBeaconBlockServer) error {
	// Right now, this streams every announced block received via p2p. It should only stream
	// finalized blocks that are canonical in the beacon node after applying the fork choice
	// rule.
	sub := s.fetcher.CanonicalBlockFeed().Subscribe(s.canonicalBlockChan)
	defer sub.Unsubscribe()
	for {
		select {
		case block := <-s.canonicalBlockChan:
			log.Info("Sending latest canonical block to RPC clients")
			if err := stream.Send(block.Proto()); err != nil {
				return err
			}
		case <-s.ctx.Done():
			log.Debug("RPC context closed, exiting goroutine")
			return nil
		}
	}
}

// LatestCrystallizedState streams the latest beacon crystallized state.
func (s *Service) LatestCrystallizedState(req *empty.Empty, stream pb.BeaconService_LatestCrystallizedStateServer) error {
	// Right now, this streams every newly created crystallized state but should only
	// stream canonical states.
	sub := s.fetcher.CanonicalCrystallizedStateFeed().Subscribe(s.canonicalStateChan)
	defer sub.Unsubscribe()
	for {
		select {
		case state := <-s.canonicalStateChan:
			log.Info("Sending crystallized state to RPC clients")
			if err := stream.Send(state.Proto()); err != nil {
				return err
			}
		case <-s.ctx.Done():
			log.Debug("RPC context closed, exiting goroutine")
			return nil
		}
	}
}
