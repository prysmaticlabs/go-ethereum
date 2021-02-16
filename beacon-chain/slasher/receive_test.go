package slasher

import (
	"context"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	logTest "github.com/sirupsen/logrus/hooks/test"

	dbtest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
	"github.com/prysmaticlabs/prysm/shared/event"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestSlasher_receiveAttestations_OK(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		serviceCfg: &ServiceConfig{
			IndexedAttsFeed: new(event.Feed),
		},
		indexedAttsChan: make(chan *ethpb.IndexedAttestation),
	}
	exitChan := make(chan struct{})
	go func() {
		s.receiveAttestations(ctx)
		exitChan <- struct{}{}
	}()
	firstIndices := []uint64{1, 2, 3}
	secondIndices := []uint64{4, 5, 6}
	att1 := &ethpb.IndexedAttestation{
		AttestingIndices: firstIndices,
		Data: &ethpb.AttestationData{
			Source: &ethpb.Checkpoint{
				Epoch: 1,
			},
			Target: &ethpb.Checkpoint{
				Epoch: 2,
			},
		},
	}
	att2 := &ethpb.IndexedAttestation{
		AttestingIndices: secondIndices,
		Data: &ethpb.AttestationData{
			Source: &ethpb.Checkpoint{
				Epoch: 1,
			},
			Target: &ethpb.Checkpoint{
				Epoch: 2,
			},
		},
	}
	s.indexedAttsChan <- att1
	s.indexedAttsChan <- att2
	cancel()
	<-exitChan
	wanted := []*slashertypes.CompactAttestation{
		{
			AttestingIndices: att1.AttestingIndices,
			Source:           att1.Data.Source.Epoch,
			Target:           att1.Data.Target.Epoch,
		},
		{
			AttestingIndices: att2.AttestingIndices,
			Source:           att2.Data.Source.Epoch,
			Target:           att2.Data.Target.Epoch,
		},
	}
	require.DeepEqual(t, wanted, s.attestationQueue)
}

func TestSlasher_receiveAttestations_OnlyValidAttestations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		serviceCfg: &ServiceConfig{
			IndexedAttsFeed: new(event.Feed),
		},
		indexedAttsChan: make(chan *ethpb.IndexedAttestation),
	}
	exitChan := make(chan struct{})
	go func() {
		s.receiveAttestations(ctx)
		exitChan <- struct{}{}
	}()
	firstIndices := []uint64{1, 2, 3}
	secondIndices := []uint64{4, 5, 6}
	// Add a valid attestation.
	validAtt := &ethpb.IndexedAttestation{
		AttestingIndices: firstIndices,
		Data: &ethpb.AttestationData{
			Source: &ethpb.Checkpoint{
				Epoch: 1,
			},
			Target: &ethpb.Checkpoint{
				Epoch: 2,
			},
		},
	}
	s.indexedAttsChan <- validAtt
	// Send an invalid, bad attestation which will not
	// pass integrity checks at it has invalid attestation data.
	s.indexedAttsChan <- &ethpb.IndexedAttestation{
		AttestingIndices: secondIndices,
	}
	cancel()
	<-exitChan
	// Expect only a single, valid attestation was added to the queue.
	require.Equal(t, 1, len(s.attestationQueue))
	wanted := []*slashertypes.CompactAttestation{
		{
			AttestingIndices: validAtt.AttestingIndices,
			Source:           validAtt.Data.Source.Epoch,
			Target:           validAtt.Data.Target.Epoch,
		},
	}
	require.DeepEqual(t, wanted, s.attestationQueue)
}

func TestSlasher_receiveBlocks_OK(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		serviceCfg: &ServiceConfig{
			BeaconBlocksFeed: new(event.Feed),
		},
		beaconBlocksChan: make(chan *ethpb.BeaconBlockHeader),
	}
	exitChan := make(chan struct{})
	go func() {
		s.receiveBlocks(ctx)
		exitChan <- struct{}{}
	}()
	block1 := &ethpb.BeaconBlockHeader{
		ProposerIndex: 1,
	}
	block2 := &ethpb.BeaconBlockHeader{
		ProposerIndex: 2,
	}
	s.beaconBlocksChan <- block1
	s.beaconBlocksChan <- block2
	cancel()
	<-exitChan
	wanted := []*slashertypes.CompactBeaconBlock{
		{
			ProposerIndex: block1.ProposerIndex,
		},
		{
			ProposerIndex: block2.ProposerIndex,
		},
	}
	require.DeepEqual(t, wanted, s.beaconBlocksQueue)
}

func TestService_processQueuedBlocks(t *testing.T) {
	hook := logTest.NewGlobal()
	beaconDB := dbtest.SetupDB(t)
	s := &Service{
		params: DefaultParams(),
		serviceCfg: &ServiceConfig{
			Database: beaconDB,
		},
		beaconBlocksQueue: []*slashertypes.CompactBeaconBlock{
			{
				ProposerIndex: 1,
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	tickerChan := make(chan uint64)
	exitChan := make(chan struct{})
	go func() {
		s.processQueuedBlocks(ctx, tickerChan)
		exitChan <- struct{}{}
	}()

	// Send a value over the ticker.
	tickerChan <- 0
	cancel()
	<-exitChan
	assert.LogsContain(t, hook, "Epoch reached, processing queued")
}