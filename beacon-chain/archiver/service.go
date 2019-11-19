package archiver

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/epoch"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/statefeed"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/validators"
	"github.com/prysmaticlabs/prysm/beacon-chain/db"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "archiver")

// Service defining archiver functionality for persisting checkpointed
// beacon chain information to a database backend for historical purposes.
type Service struct {
	ctx           context.Context
	cancel        context.CancelFunc
	beaconDB      db.Database
	headFetcher   blockchain.HeadFetcher
	stateNotifier blockchain.StateNotifier
}

// Config options for the archiver service.
type Config struct {
	BeaconDB      db.Database
	HeadFetcher   blockchain.HeadFetcher
	StateNotifier blockchain.StateNotifier
}

// NewArchiverService initializes the service from configuration options.
func NewArchiverService(ctx context.Context, cfg *Config) *Service {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		ctx:           ctx,
		cancel:        cancel,
		beaconDB:      cfg.BeaconDB,
		headFetcher:   cfg.HeadFetcher,
		stateNotifier: cfg.StateNotifier,
	}
}

// Start the archiver service event loop.
func (s *Service) Start() {
	go s.run(s.ctx)
}

// Stop the archiver service event loop.
func (s *Service) Stop() error {
	defer s.cancel()
	return nil
}

// Status reports the healthy status of the archiver. Returning nil means service
// is correctly running without error.
func (s *Service) Status() error {
	return nil
}

// We archive committee information pertaining to the head state's epoch.
func (s *Service) archiveCommitteeInfo(ctx context.Context, headState *pb.BeaconState) error {
	currentEpoch := helpers.SlotToEpoch(headState.Slot)
	proposerSeed, err := helpers.Seed(headState, currentEpoch, params.BeaconConfig().DomainBeaconProposer)
	if err != nil {
		return errors.Wrap(err, "could not generate seed")
	}
	attesterSeed, err := helpers.Seed(headState, currentEpoch, params.BeaconConfig().DomainBeaconAttester)
	if err != nil {
		return errors.Wrap(err, "could not generate seed")
	}

	info := &ethpb.ArchivedCommitteeInfo{
		ProposerSeed: proposerSeed[:],
		AttesterSeed: attesterSeed[:],
	}
	if err := s.beaconDB.SaveArchivedCommitteeInfo(ctx, currentEpoch, info); err != nil {
		return errors.Wrap(err, "could not archive committee info")
	}
	return nil
}

// We archive active validator set changes that happened during the previous epoch.
func (s *Service) archiveActiveSetChanges(ctx context.Context, headState *pb.BeaconState) error {
	activations := validators.ActivatedValidatorIndices(helpers.PrevEpoch(headState), headState.Validators)
	slashings := validators.SlashedValidatorIndices(helpers.PrevEpoch(headState), headState.Validators)
	activeValidatorCount, err := helpers.ActiveValidatorCount(headState, helpers.PrevEpoch(headState))
	if err != nil {
		return errors.Wrap(err, "could not get active validator count")
	}
	exited, err := validators.ExitedValidatorIndices(headState.Validators, activeValidatorCount)
	if err != nil {
		return errors.Wrap(err, "could not determine exited validator indices")
	}
	activeSetChanges := &ethpb.ArchivedActiveSetChanges{
		Activated: activations,
		Exited:    exited,
		Slashed:   slashings,
	}
	if err := s.beaconDB.SaveArchivedActiveValidatorChanges(ctx, helpers.PrevEpoch(headState), activeSetChanges); err != nil {
		return errors.Wrap(err, "could not archive active validator set changes")
	}
	return nil
}

// We compute participation metrics by first retrieving the head state and
// matching validator attestations during the epoch.
func (s *Service) archiveParticipation(ctx context.Context, headState *pb.BeaconState) error {
	participation, err := epoch.ComputeValidatorParticipation(headState, helpers.SlotToEpoch(headState.Slot))
	if err != nil {
		return errors.Wrap(err, "could not compute participation")
	}
	return s.beaconDB.SaveArchivedValidatorParticipation(ctx, helpers.SlotToEpoch(headState.Slot), participation)
}

// We archive validator balances and active indices.
func (s *Service) archiveBalances(ctx context.Context, headState *pb.BeaconState) error {
	balances := headState.Balances
	currentEpoch := helpers.CurrentEpoch(headState)
	if err := s.beaconDB.SaveArchivedBalances(ctx, currentEpoch, balances); err != nil {
		return errors.Wrap(err, "could not archive balances")
	}
	return nil
}

func (s *Service) run(ctx context.Context) {
	subChannel := make(chan *statefeed.Event, 1)
	sub := s.stateNotifier.StateFeed().Subscribe(subChannel)
	defer sub.Unsubscribe()
	for {
		select {
		case event := <-subChannel:
			switch event.Type {
			case statefeed.BlockProcessed:
				data := event.Data.(*statefeed.BlockProcessedData)
				log.WithField("headRoot", fmt.Sprintf("%#x", data.BlockRoot)).Debug("Received block processed event")
				headState, err := s.headFetcher.HeadState(ctx)
				if err != nil {
					log.WithError(err).Error("Head state is not available")
					continue
				}
				if !helpers.IsEpochEnd(headState.Slot) {
					continue
				}
				if err := s.archiveCommitteeInfo(ctx, headState); err != nil {
					log.WithError(err).Error("Could not archive committee info")
					continue
				}
				if err := s.archiveActiveSetChanges(ctx, headState); err != nil {
					log.WithError(err).Error("Could not archive active validator set changes")
					continue
				}
				if err := s.archiveParticipation(ctx, headState); err != nil {
					log.WithError(err).Error("Could not archive validator participation")
					continue
				}
				if err := s.archiveBalances(ctx, headState); err != nil {
					log.WithError(err).Error("Could not archive validator balances and active indices")
					continue
				}
				log.WithField(
					"epoch",
					helpers.CurrentEpoch(headState),
				).Debug("Successfully archived beacon chain data during epoch")
			}
		case <-s.ctx.Done():
			log.Debug("Context closed, exiting goroutine")
			return
		case err := <-sub.Err():
			log.WithError(err).Error("Subscription to new chain head notifier failed")
			return
		}
	}
}
