package spectest

import (
	"fmt"
	"path"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/epoch"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params/spectest"
	"github.com/prysmaticlabs/prysm/shared/testutil"
)

// This is a subset of state.ProcessEpoch. The spec test defines input data for
// `justification_and_finalization` only.
func processJustificationAndFinalizationWrapper(state *pb.BeaconState) (*pb.BeaconState, error) {
	helpers.ClearAllCaches()

	// This process mutates the state, so we'll make a copy in order to print debug before/after.
	state = proto.Clone(state).(*pb.BeaconState)

	prevEpochAtts, err := epoch.MatchAttestations(state, helpers.PrevEpoch(state))
	if err != nil {
		return nil, fmt.Errorf("could not get target atts prev epoch %d: %v",
			helpers.PrevEpoch(state), err)
	}
	currentEpochAtts, err := epoch.MatchAttestations(state, helpers.CurrentEpoch(state))
	if err != nil {
		return nil, fmt.Errorf("could not get target atts current epoch %d: %v",
			helpers.CurrentEpoch(state), err)
	}
	prevEpochAttestedBalance, err := epoch.AttestingBalance(state, prevEpochAtts.Target)
	if err != nil {
		return nil, fmt.Errorf("could not get attesting balance prev epoch: %v", err)
	}
	currentEpochAttestedBalance, err := epoch.AttestingBalance(state, currentEpochAtts.Target)
	if err != nil {
		return nil, fmt.Errorf("could not get attesting balance current epoch: %v", err)
	}

	state, err = epoch.ProcessJustificationAndFinalization(state, prevEpochAttestedBalance, currentEpochAttestedBalance)
	if err != nil {
		return nil, fmt.Errorf("could not process justification: %v", err)
	}

	return state, nil
}

func runJustificationAndFinalizationTests(t *testing.T, config string) {
	if err := spectest.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	testFolders, testsFolderPath := testutil.TestFolders(t, config, "phase0/epoch_processing/justification_and_finalization")

	for _, folder := range testFolders {
		t.Run(folder.Name(), func(t *testing.T) {
			folderPath := path.Join(testsFolderPath, folder.Name())
			testutil.RunEpochOperationTest(t, folderPath, processJustificationAndFinalizationWrapper)
		})
	}
}
