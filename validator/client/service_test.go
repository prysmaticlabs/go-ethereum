package client

import (
	"context"
	"testing"
	"time"

	"github.com/prysmaticlabs/prysm/shared/testutil"

	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"github.com/prysmaticlabs/prysm/shared"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

var _ = shared.Service(&ValidatorService{})

func TestStop_cancelsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	vs := &ValidatorService{
		ctx:    ctx,
		cancel: cancel,
	}

	if err := vs.Stop(); err != nil {
		t.Error(err)
	}

	select {
	case <-time.After(1 * time.Second):
		t.Error("ctx not cancelled within 1s")
	case <-vs.ctx.Done():
	}
}

func TestLifecycle(t *testing.T) {
	hook := logTest.NewGlobal()
	validatorService := NewValidatorService(
		context.Background(),
		&Config{
			Endpoint: "merkle tries",
			CertFlag: "alice.crt",
		},
	)
	validatorService.Start()
	if err := validatorService.Stop(); err != nil {
		t.Fatalf("Could not stop service: %v", err)
	}
	testutil.AssertLogsContain(t, hook, "Stopping service")
}

func TestInsecure(t *testing.T) {
	hook := logTest.NewGlobal()
	validatorService := NewValidatorService(
		context.Background(),
		&Config{
			Endpoint: "merkle tries",
		},
	)
	validatorService.Start()
	testutil.AssertLogsContain(t, hook, "You are using an insecure gRPC connection")
	if err := validatorService.Stop(); err != nil {
		t.Fatalf("Could not stop service: %v", err)
	}
	testutil.AssertLogsContain(t, hook, "Stopping service")
}

func TestBeaconServiceClient(t *testing.T) {
	validatorService := NewValidatorService(
		context.Background(),
		&Config{
			Endpoint: "merkle tries",
		},
	)
	validatorService.conn = nil
	client := validatorService.BeaconServiceClient()
	if _, ok := client.(pb.BeaconServiceClient); !ok {
		t.Error("Beacon service client function does not implement interface")
	}
}

func TestProposerServiceClient(t *testing.T) {
	validatorService := NewValidatorService(
		context.Background(),
		&Config{
			Endpoint: "merkle tries",
		},
	)
	validatorService.conn = nil
	client := validatorService.ProposerServiceClient()
	if _, ok := client.(pb.ProposerServiceClient); !ok {
		t.Error("Beacon service client function does not implement interface")
	}
}

func TestAttesterServiceClient(t *testing.T) {
	validatorService := NewValidatorService(
		context.Background(),
		&Config{
			Endpoint: "merkle tries",
		},
	)
	validatorService.conn = nil
	client := validatorService.AttesterServiceClient()
	if _, ok := client.(pb.AttesterServiceClient); !ok {
		t.Error("Beacon service client function does not implement interface")
	}
}
