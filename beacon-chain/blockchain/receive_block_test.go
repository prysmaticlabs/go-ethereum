package blockchain

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"path"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	blockchainTesting "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/cache"
	"github.com/prysmaticlabs/prysm/beacon-chain/db/kv"
	"github.com/prysmaticlabs/prysm/beacon-chain/forkchoice/protoarray"
	"github.com/prysmaticlabs/prysm/beacon-chain/operations/attestations"
	"github.com/prysmaticlabs/prysm/beacon-chain/operations/voluntaryexits"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stategen"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
)

func TestService_ReceiveBlockNoPubsub(t *testing.T) {
	ctx := context.Background()

	genesis, keys := testutil.DeterministicGenesisState(t, 64)
	genFullBlock := func(t *testing.T, conf *testutil.BlockGenConfig, slot uint64) *ethpb.SignedBeaconBlock {
		blk, err := testutil.GenerateFullBlock(genesis, keys, conf, slot)
		if err != nil {
			t.Error(err)
		}
		return blk
	}
	bc := params.BeaconConfig()
	bc.ShardCommitteePeriod = 0 // Required for voluntary exits test in reasonable time.
	params.OverrideBeaconConfig(bc)

	type args struct {
		block *ethpb.SignedBeaconBlock
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		check   func(*testing.T, *Service)
	}{
		{
			name: "applies block with state transition",
			args: args{
				block: genFullBlock(t, testutil.DefaultBlockGenConfig(), 2 /*slot*/),
			},
			check: func(t *testing.T, s *Service) {
				if hs := s.head.state.Slot(); hs != 2 {
					t.Errorf("Unexpected state slot. Got %d but wanted %d", hs, 2)
				}
				if bs := s.head.block.Block.Slot; bs != 2 {
					t.Errorf("Unexpected head block slot. Got %d but wanted %d", bs, 2)
				}
			},
		},
		{
			name: "saves attestations to pool",
			args: args{
				block: genFullBlock(t,
					&testutil.BlockGenConfig{
						NumProposerSlashings: 0,
						NumAttesterSlashings: 0,
						NumAttestations:      2,
						NumDeposits:          0,
						NumVoluntaryExits:    0,
					},
					1, /*slot*/
				),
			},
			check: func(t *testing.T, s *Service) {
				if baCount := len(s.attPool.BlockAttestations()); baCount != 2 {
					t.Errorf("Did not get the correct number of block attestations saved to the pool. "+
						"Got %d but wanted %d", baCount, 2)
				}
			},
		},

		{
			name: "updates exit pool",
			args: args{
				block: genFullBlock(t, &testutil.BlockGenConfig{
					NumProposerSlashings: 0,
					NumAttesterSlashings: 0,
					NumAttestations:      0,
					NumDeposits:          0,
					NumVoluntaryExits:    3,
				},
					1, /*slot*/
				),
			},
			check: func(t *testing.T, s *Service) {
				var n int
				for i := uint64(0); int(i) < genesis.NumValidators(); i++ {
					if s.exitPool.HasBeenIncluded(i) {
						n++
					}
				}
				if n != 3 {
					t.Errorf("Did not mark the correct number of exits. Got %d but wanted %d", n, 3)
				}
			},
		},
		{
			name: "sets epochParticipation",
			args: args{
				block: genFullBlock(t, testutil.DefaultBlockGenConfig(), 1 /*slot*/),
			},
			check: func(t *testing.T, s *Service) {
				e := s.epochParticipation[0]
				if e == nil {
					t.Error("epoch participation is nil")
				}
			},
		},
		{
			name: "notifies block processed on state feed",
			args: args{
				block: genFullBlock(t, testutil.DefaultBlockGenConfig(), 1 /*slot*/),
			},
			check: func(t *testing.T, s *Service) {
				t.Skip("TODO: Need to add a better mock for state notifier.")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randPath, err := rand.Int(rand.Reader, big.NewInt(1000000))
			if err != nil {
				t.Fatalf("could not generate random file path: %v", err)
			}
			p := path.Join(testutil.TempDir(), fmt.Sprintf("/%d", randPath))
			stateSummaryCache := cache.NewStateSummaryCache()
			db, err := kv.NewKVStore(p, stateSummaryCache)
			if err != nil {
				t.Fatal(err)
			}
			genesisBlockRoot := bytesutil.ToBytes32(nil)
			if err := db.SaveState(ctx, genesis, genesisBlockRoot); err != nil {
				t.Fatal(err)
			}

			cfg := &Config{
				BeaconDB: db,
				ForkChoiceStore: protoarray.New(
					0, // justifiedEpoch
					0, // finalizedEpoch
					genesisBlockRoot,
				),
				AttPool:       attestations.NewPool(),
				ExitPool:      voluntaryexits.NewPool(),
				StateNotifier: &blockchainTesting.MockStateNotifier{},
				StateGen:      stategen.New(db, stateSummaryCache),
			}
			s, err := NewService(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}
			if err := s.saveGenesisData(ctx, genesis); err != nil {
				t.Fatal(err)
			}
			root, err := stateutil.BlockRoot(tt.args.block.Block)
			if err != nil {
				t.Error(err)
			}
			if err := s.ReceiveBlockNoPubsub(ctx, tt.args.block, root); (err != nil) != tt.wantErr {
				t.Errorf("ReceiveBlockNoPubsub() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				tt.check(t, s)
			}
		})
	}
}
