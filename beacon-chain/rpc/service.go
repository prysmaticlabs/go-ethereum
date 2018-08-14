// Package rpc defines the services that the beacon-chain uses to communicate via gRPC.
package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/prysmaticlabs/prysm/beacon-chain/types"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logrus.WithField("prefix", "rpc")

// Service defining an RPC server for a beacon node.
type Service struct {
	ctx       context.Context
	cancel    context.CancelFunc
	announcer types.CanonicalEventAnnouncer
	port      string
	listener  net.Listener
	withCert  string
	withKey   string
}

// Config options for the beacon node RPC server.
type Config struct {
	Port     string
	CertFlag string
	KeyFlag  string
}

// NewRPCService creates a new instance of a struct implementing the BeaconServiceServer
// interface.
func NewRPCService(ctx context.Context, cfg *Config, announcer types.CanonicalEventAnnouncer) *Service {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		ctx:       ctx,
		cancel:    cancel,
		announcer: announcer,
		port:      cfg.Port,
		withCert:  cfg.CertFlag,
		withKey:   cfg.KeyFlag,
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

	var grpcServer *grpc.Server
	if s.withCert != "" && s.withKey != "" {
		creds, err := credentials.NewServerTLSFromFile(s.withCert, s.withKey)
		if err != nil {
			log.Errorf("could not load TLS keys: %s", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
	} else {
		log.Warn("You are using an insecure gRPC connection! Please provide a certificate and key to use a secure connection")
		grpcServer = grpc.NewServer()
	}

	pb.RegisterBeaconServiceServer(grpcServer, s)
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Debugf("Could not serve gRPC: %v", err)
		}
	}()
}

// Stop the service.
func (s *Service) Stop() error {
	log.Info("Stopping service")
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// FetchShuffledValidatorIndices retrieves the shuffled validator indices, cutoffs, and
// assigned attestation heights at a given crystallized state hash.
// This function can be called by clients to fetch a historical list of shuffled
// validators ata point in time corresponding to a certain crystallized state.
func (s *Service) FetchShuffledValidatorIndices(ctx context.Context, req *pb.ShuffleRequest) (*pb.ShuffleResponse, error) {
	shuffleRes := &pb.ShuffleResponse{
		ShuffledValidatorIndices:   []uint64{2, 1, 0},
		CutoffIndices:              []uint64{},
		AssignedAttestationHeights: []uint64{},
	}
	return shuffleRes, nil
}

// ProposeBlock is called by a proposer in a sharding client and a full beacon node
// sends the request into a beacon block that can then be included in a canonical chain.
//
// TODO: needs implementation.
func (s *Service) ProposeBlock(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeResponse, error) {
	// TODO: implement.
	return nil, errors.New("unimplemented")
}

// SignBlock is a function called by an attester in a sharding client to sign off
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
	for {
		select {
		case block := <-s.announcer.CanonicalBlockAnnouncement():
			if err := stream.Send(block.Proto()); err != nil {
				return err
			}
		case <-s.ctx.Done():
			return nil
		}
	}
}

// LatestCrystallizedState streams the latest beacon crystallized state.
func (s *Service) LatestCrystallizedState(req *empty.Empty, stream pb.BeaconService_LatestCrystallizedStateServer) error {
	// Right now, this streams every newly created crystallized state but should only
	// stream canonical states.
	for {
		select {
		case state := <-s.announcer.CanonicalCrystallizedStateAnnouncement():
			if err := stream.Send(state.Proto()); err != nil {
				return err
			}
		case <-s.ctx.Done():
			return nil
		}
	}
}
