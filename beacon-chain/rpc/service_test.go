package rpc

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	mock "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	mockPOW "github.com/prysmaticlabs/prysm/beacon-chain/powchain/testing"
	mockSync "github.com/prysmaticlabs/prysm/beacon-chain/sync/initial-sync/testing"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(ioutil.Discard)
}

func TestLifecycle_OK(t *testing.T) {
	hook := logTest.NewGlobal()
	chainService := &mock.ChainService{
		Genesis: time.Now(),
	}
	rpcService, serviceCtx := NewService(&Config{
		Port:                "7348",
		SyncService:         &mockSync.Sync{IsSyncing: false},
		BlockReceiver:       chainService,
		AttestationReceiver: chainService,
		HeadFetcher:         chainService,
		GenesisTimeFetcher:  chainService,
		POWChainService:     &mockPOW.POWChain{},
		StateNotifier:       chainService.StateNotifier(),
	})

	rpcService.Start(serviceCtx.Ctx)

	require.LogsContain(t, hook, "listening on port")
	assert.NoError(t, rpcService.Stop(serviceCtx.Ctx))
}

func TestStatus_CredentialError(t *testing.T) {
	credentialErr := errors.New("credentialError")
	s := &Service{credentialError: credentialErr}

	assert.ErrorContains(t, s.credentialError.Error(), s.Status())
}

func TestRPC_InsecureEndpoint(t *testing.T) {
	hook := logTest.NewGlobal()
	chainService := &mock.ChainService{Genesis: time.Now()}
	rpcService, serviceCtx := NewService(&Config{
		Port:                "7777",
		SyncService:         &mockSync.Sync{IsSyncing: false},
		BlockReceiver:       chainService,
		GenesisTimeFetcher:  chainService,
		AttestationReceiver: chainService,
		HeadFetcher:         chainService,
		POWChainService:     &mockPOW.POWChain{},
		StateNotifier:       chainService.StateNotifier(),
	})

	rpcService.Start(serviceCtx.Ctx)

	require.LogsContain(t, hook, "listening on port")
	require.LogsContain(t, hook, "You are using an insecure gRPC server")
	assert.NoError(t, rpcService.Stop(serviceCtx.Ctx))
}
