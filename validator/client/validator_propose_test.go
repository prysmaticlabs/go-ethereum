package client

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-bitfield"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	slashpb "github.com/prysmaticlabs/prysm/proto/slashing"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/validator/db"
	"github.com/prysmaticlabs/prysm/validator/internal"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

type mocks struct {
	proposerClient   *internal.MockProposerServiceClient
	validatorClient  *internal.MockValidatorServiceClient
	attesterClient   *internal.MockAttesterServiceClient
	aggregatorClient *internal.MockAggregatorServiceClient
}

func setup(t *testing.T) (*validator, *mocks, func()) {
	db := db.SetupDB(t, pubKeysFromMap(keyMap))
	ctrl := gomock.NewController(t)
	m := &mocks{
		proposerClient:   internal.NewMockProposerServiceClient(ctrl),
		validatorClient:  internal.NewMockValidatorServiceClient(ctrl),
		attesterClient:   internal.NewMockAttesterServiceClient(ctrl),
		aggregatorClient: internal.NewMockAggregatorServiceClient(ctrl),
	}
	validator := &validator{
		db:               db,
		proposerClient:   m.proposerClient,
		attesterClient:   m.attesterClient,
		validatorClient:  m.validatorClient,
		aggregatorClient: m.aggregatorClient,
		keys:             keyMap,
		graffiti:         []byte{},
		attLogs:          make(map[[32]byte]*attSubmitted),
	}

	return validator, m, ctrl.Finish
}

func TestProposeBlock_DoesNotProposeGenesisBlock(t *testing.T) {
	hook := logTest.NewGlobal()
	validator, _, finish := setup(t)
	defer finish()
	validator.ProposeBlock(context.Background(), 0, validatorPubKey)

	testutil.AssertLogsContain(t, hook, "Assigned to genesis slot, skipping proposal")
}

func TestProposeBlock_DomainDataFailed(t *testing.T) {
	hook := logTest.NewGlobal()
	validator, m, finish := setup(t)
	defer finish()

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), // epoch
	).Return(nil /*response*/, errors.New("uh oh"))

	validator.ProposeBlock(context.Background(), 1, validatorPubKey)
	testutil.AssertLogsContain(t, hook, "Failed to sign randao reveal")
}

func TestProposeBlock_RequestBlockFailed(t *testing.T) {
	hook := logTest.NewGlobal()
	validator, m, finish := setup(t)
	defer finish()

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), // epoch
	).Return(&pb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().RequestBlock(
		gomock.Any(), // ctx
		gomock.Any(), // block request
	).Return(nil /*response*/, errors.New("uh oh"))

	validator.ProposeBlock(context.Background(), 1, validatorPubKey)
	testutil.AssertLogsContain(t, hook, "Failed to request block from beacon node")
}

func TestProposeBlock_ProposeBlockFailed(t *testing.T) {
	hook := logTest.NewGlobal()
	validator, m, finish := setup(t)
	defer finish()

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Return(&pb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().RequestBlock(
		gomock.Any(), // ctx
		gomock.Any(),
	).Return(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{}}, nil /*err*/)

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Return(&pb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().ProposeBlock(
		gomock.Any(), // ctx
		gomock.AssignableToTypeOf(&ethpb.BeaconBlock{}),
	).Return(nil /*response*/, errors.New("uh oh"))

	validator.ProposeBlock(context.Background(), 1, validatorPubKey)
	testutil.AssertLogsContain(t, hook, "Failed to propose block")
}

func TestProposeBlock_BlocksDoubleProposal(t *testing.T) {
	hook := logTest.NewGlobal()
	validator, m, finish := setup(t)
	defer finish()

	history := &slashpb.ProposalHistory{
		EpochBits:          bitfield.NewBitlist(params.BeaconConfig().WeakSubjectivityPeriod),
		LatestEpochWritten: 0,
	}
	if err := validator.db.SaveProposalHistory(validatorPubKey[:], history); err != nil {
		t.Fatal(err)
	}

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Times(2).Return(&ethpb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().RequestBlock(
		gomock.Any(), // ctx
		gomock.Any(),
	).Times(2).Return(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{}}, nil /*err*/)

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Return(&ethpb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().ProposeBlock(
		gomock.Any(), // ctx
		gomock.AssignableToTypeOf(&ethpb.BeaconBlock{}),
	).Return(&pb.ProposeResponse{}, nil /*error*/)

	validator.ProposeBlock(context.Background(), params.BeaconConfig().SlotsPerEpoch*5+2, validatorPubKey)
	testutil.AssertLogsDoNotContain(t, hook, "Tried to sign a double proposal")

	validator.ProposeBlock(context.Background(), params.BeaconConfig().SlotsPerEpoch*5+2, validatorPubKey)
	testutil.AssertLogsContain(t, hook, "Tried to sign a double proposal")
}

func TestProposeBlock_BroadcastsBlock(t *testing.T) {
	validator, m, finish := setup(t)
	defer finish()

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Return(&pb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().RequestBlock(
		gomock.Any(), // ctx
		gomock.Any(),
	).Return(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{}}, nil /*err*/)

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Return(&pb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().ProposeBlock(
		gomock.Any(), // ctx
		gomock.AssignableToTypeOf(&ethpb.BeaconBlock{}),
	).Return(&pb.ProposeResponse{}, nil /*error*/)

	m.validatorClient.EXPECT().ValidatorIndex(
		gomock.Any(), // ctx
		gomock.AssignableToTypeOf(&pb.ValidatorIndexRequest{}),
	).Return(&pb.ValidatorIndexResponse{}, nil)

	validator.ProposeBlock(context.Background(), 1, validatorPubKey)
}

func TestProposeBlock_BroadcastsBlock_WithGraffiti(t *testing.T) {
	validator, m, finish := setup(t)
	defer finish()

	validator.graffiti = []byte("12345678901234567890123456789012")

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Return(&pb.DomainResponse{}, nil /*err*/)

	m.proposerClient.EXPECT().RequestBlock(
		gomock.Any(), // ctx
		gomock.Any(),
	).Return(&ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{Graffiti: validator.graffiti}}, nil /*err*/)

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		gomock.Any(), //epoch
	).Return(&pb.DomainResponse{}, nil /*err*/)

	var sentBlock *ethpb.BeaconBlock

	m.proposerClient.EXPECT().ProposeBlock(
		gomock.Any(), // ctx
		gomock.AssignableToTypeOf(&ethpb.BeaconBlock{}),
	).DoAndReturn(func(ctx context.Context, block *ethpb.BeaconBlock) (*pb.ProposeResponse, error) {
		sentBlock = block
		return &pb.ProposeResponse{}, nil
	})

	m.validatorClient.EXPECT().ValidatorIndex(
		gomock.Any(), // ctx
		gomock.AssignableToTypeOf(&pb.ValidatorIndexRequest{}),
	).Return(&pb.ValidatorIndexResponse{}, nil)

	validator.ProposeBlock(context.Background(), 1, validatorPubKey)

	if string(sentBlock.Body.Graffiti) != string(validator.graffiti) {
		t.Errorf("Block was broadcast with the wrong graffiti field, wanted \"%v\", got \"%v\"", string(validator.graffiti), string(sentBlock.Body.Graffiti))
	}
}
