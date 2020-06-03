package endtoend

import (
	"testing"

	ev "github.com/prysmaticlabs/prysm/endtoend/evaluators"
	e2eParams "github.com/prysmaticlabs/prysm/endtoend/params"
	"github.com/prysmaticlabs/prysm/endtoend/types"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
)

func TestEndToEnd_Slashing_MinimalConfig(t *testing.T) {
	t.Skip("Skipping until eth1 changes in v0.12 can work with e2e")
	testutil.ResetCache()
	params.UseE2EConfig()

	minimalConfig := &types.E2EConfig{
		BeaconFlags:    []string{},
		ValidatorFlags: []string{},
		EpochsToRun:    3,
		TestSync:       false,
		TestSlasher:    true,
		TestDeposits:   false,
		Evaluators: []types.Evaluator{
			ev.PeersConnect,
			ev.HealthzCheck,
			ev.ValidatorsSlashed,
			ev.SlashedValidatorsLoseBalance,
			ev.InjectDoubleVote,
			ev.ProposeDoubleBlock,
		},
	}
	if err := e2eParams.Init(2); err != nil {
		t.Fatal(err)
	}

	runEndToEndTest(t, minimalConfig)
}
