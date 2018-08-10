// Package node defines a backend for a sharding-enabled, Ethereum blockchain.
// It defines a struct which handles the lifecycle of services in the
// sharding system, providing a bridge to the main Ethereum blockchain,
// as well as instantiating peer-to-peer networking for shards.
package node

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/node"
	"github.com/prysmaticlabs/prysm/client/attester"
	"github.com/prysmaticlabs/prysm/client/mainchain"
	"github.com/prysmaticlabs/prysm/client/observer"
	"github.com/prysmaticlabs/prysm/client/params"
	"github.com/prysmaticlabs/prysm/client/proposer"
	"github.com/prysmaticlabs/prysm/client/rpcclient"
	"github.com/prysmaticlabs/prysm/client/syncer"
	"github.com/prysmaticlabs/prysm/client/txpool"
	"github.com/prysmaticlabs/prysm/client/types"
	"github.com/prysmaticlabs/prysm/shared"
	"github.com/prysmaticlabs/prysm/shared/cmd"
	"github.com/prysmaticlabs/prysm/shared/database"
	"github.com/prysmaticlabs/prysm/shared/debug"
	"github.com/prysmaticlabs/prysm/shared/p2p"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var log = logrus.WithField("prefix", "node")

const shardChainDBName = "shardchaindata"

// ShardEthereum is a service that is registered and started when geth is launched.
// it contains APIs and fields that handle the different components of the sharded
// Ethereum network.
type ShardEthereum struct {
	shardConfig *params.Config // Holds necessary information to configure shards.

	// Lifecycle and service stores.
	services *shared.ServiceRegistry
	lock     sync.RWMutex
	stop     chan struct{} // Channel to wait for termination notifications.
	db       *database.DB
}

// NewShardInstance creates a new sharding-enabled Ethereum instance. This is called in the main
// geth sharding entrypoint.
func NewShardInstance(ctx *cli.Context) (*ShardEthereum, error) {
	registry := shared.NewServiceRegistry()
	shardEthereum := &ShardEthereum{
		services: registry,
		stop:     make(chan struct{}),
	}

	// Configure shardConfig by loading the default.
	shardEthereum.shardConfig = params.DefaultConfig()

	if err := shardEthereum.startDB(ctx); err != nil {
		return nil, err
	}

	if err := shardEthereum.registerP2P(); err != nil {
		return nil, err
	}

	if err := shardEthereum.registerMainchainClient(ctx); err != nil {
		return nil, err
	}

	actorFlag := ctx.GlobalString(types.ActorFlag.Name)
	if err := shardEthereum.registerTXPool(actorFlag); err != nil {
		return nil, err
	}

	shardIDFlag := ctx.GlobalInt(types.ShardIDFlag.Name)
	if err := shardEthereum.registerSyncerService(shardEthereum.shardConfig, shardIDFlag); err != nil {
		return nil, err
	}

	if err := shardEthereum.registerBeaconRPCService(ctx); err != nil {
		return nil, err
	}

	if err := shardEthereum.registerActorService(shardEthereum.shardConfig, actorFlag, shardIDFlag); err != nil {
		return nil, err
	}

	return shardEthereum, nil
}

// Start the ShardEthereum service and kicks off the p2p and actor's main loop.
func (s *ShardEthereum) Start() {
	s.lock.Lock()

	log.Info("Starting sharding node")

	s.services.StartAll()

	stop := s.stop
	s.lock.Unlock()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		log.Info("Got interrupt, shutting down...")
		go s.Close()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.Info("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
		debug.Exit() // Ensure trace and CPU profile data are flushed.
		panic("Panic closing the sharding node")
	}()

	// Wait for stop channel to be closed.
	<-stop
}

// Close handles graceful shutdown of the system.
func (s *ShardEthereum) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.db.Close()
	s.services.StopAll()
	log.Info("Stopping sharding node")

	close(s.stop)
}

// startDB attaches a LevelDB wrapped object to the shardEthereum instance.
func (s *ShardEthereum) startDB(ctx *cli.Context) error {
	path := ctx.GlobalString(cmd.DataDirFlag.Name)
	config := &database.DBConfig{DataDir: path, Name: shardChainDBName, InMemory: false}
	db, err := database.NewDB(config)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

// registerP2P attaches a p2p server to the ShardEthereum instance.
func (s *ShardEthereum) registerP2P() error {
	shardp2p, err := p2p.NewServer()
	if err != nil {
		return fmt.Errorf("could not register shardp2p service: %v", err)
	}
	return s.services.RegisterService(shardp2p)
}

// registerMainchainClient registers the mainchain client to the shard services.
func (s *ShardEthereum) registerMainchainClient(ctx *cli.Context) error {
	path := node.DefaultDataDir()
	if ctx.GlobalIsSet(cmd.DataDirFlag.Name) {
		path = ctx.GlobalString(cmd.DataDirFlag.Name)
	}

	endpoint := ctx.Args().First()
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s/%s.ipc", path, mainchain.ClientIdentifier)
	}
	if ctx.GlobalIsSet(cmd.RPCProviderFlag.Name) {
		endpoint = ctx.GlobalString(cmd.RPCProviderFlag.Name)
	} else if ctx.GlobalIsSet(cmd.IPCPathFlag.Name) {
		endpoint = ctx.GlobalString(cmd.IPCPathFlag.Name)
	}
	passwordFile := ctx.GlobalString(cmd.PasswordFileFlag.Name)
	depositFlag := ctx.GlobalBool(types.DepositFlag.Name)

	client, err := mainchain.NewSMCClient(endpoint, path, depositFlag, passwordFile)
	if err != nil {
		return fmt.Errorf("could not register smc client service: %v", err)
	}
	return s.services.RegisterService(client)
}

// registerTXPool is only relevant to proposers in the sharded system. It will
// spin up a transaction pool that will relay incoming transactions via an
// event feed. For our first releases, this can just relay test/fake transaction data
// the proposer can serialize into collation blobs.
// TODO: design this txpool system for our first release.
func (s *ShardEthereum) registerTXPool(actor string) error {
	if actor != "proposer" {
		return nil
	}
	var shardp2p *p2p.Server
	if err := s.services.FetchService(&shardp2p); err != nil {
		return err
	}
	pool, err := txpool.NewTXPool(shardp2p)
	if err != nil {
		return fmt.Errorf("could not register shard txpool service: %v", err)
	}
	return s.services.RegisterService(pool)
}

// registerActorService registers the actor according to CLI flags. Either attester/proposer/observer.
func (s *ShardEthereum) registerActorService(config *params.Config, actor string, shardID int) error {
	var shardp2p *p2p.Server
	if err := s.services.FetchService(&shardp2p); err != nil {
		return err
	}

	var client *mainchain.SMCClient
	if err := s.services.FetchService(&client); err != nil {
		return err
	}

	var sync *syncer.Syncer
	if err := s.services.FetchService(&sync); err != nil {
		return err
	}

	var rpcService *rpcclient.Service
	if err := s.services.FetchService(&rpcService); err != nil {
		return err
	}

	switch actor {
	case "attester":
		att := attester.NewAttester(context.TODO(), rpcService)
		return s.services.RegisterService(att)
	case "proposer":
		prop := proposer.NewProposer(context.TODO(), rpcService)
		return s.services.RegisterService(prop)
	default:
		obs, err := observer.NewObserver(shardp2p, s.db, shardID, sync, client)
		if err != nil {
			return fmt.Errorf("could not register observer service: %v", err)
		}
		return s.services.RegisterService(obs)
	}
}

// registerSyncerService adds the p2p and mainchain services to the syncer and
// registers the syncer to the shard.
func (s *ShardEthereum) registerSyncerService(config *params.Config, shardID int) error {
	var shardp2p *p2p.Server
	if err := s.services.FetchService(&shardp2p); err != nil {
		return err
	}
	var client *mainchain.SMCClient
	if err := s.services.FetchService(&client); err != nil {
		return err
	}

	sync, err := syncer.NewSyncer(config, client, shardp2p, s.db, shardID)
	if err != nil {
		return fmt.Errorf("could not register syncer service: %v", err)
	}
	return s.services.RegisterService(sync)
}

// registersBeaconRPCService registers a new RPC client for the shard.
func (s *ShardEthereum) registerBeaconRPCService(ctx *cli.Context) error {
	endpoint := ctx.GlobalString(types.BeaconRPCProviderFlag.Name)
	rpcService := rpcclient.NewRPCClient(context.TODO(), &rpcclient.Config{
		Endpoint: endpoint,
	})
	return s.services.RegisterService(rpcService)
}
