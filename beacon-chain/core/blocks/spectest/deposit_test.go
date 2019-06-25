package spectest

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state/stateutils"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params/spectest"
	"github.com/prysmaticlabs/prysm/shared/testutil"
)

func runDepositTest(t *testing.T, filename string) {
	file, err := ioutil.ReadFile("deposit_minimal.yaml")
	if err != nil {
		t.Fatalf("Could not load file %v", err)
	}

	test := &DepositsMinimal{}
	if err := yaml.Unmarshal(file, test); err != nil {
		t.Fatalf("Failed to Unmarshal: %v", err)
	}

	if err := spectest.SetConfig(test.Config); err != nil {
		t.Fatal(err)
	}

	for _, tt := range test.TestCases {
		t.Run(tt.Description, func(t *testing.T) {
			preState := &pb.BeaconState{}
			if err = testutil.ConvertToPb(tt.Pre, preState); err != nil {
				t.Fatal(err)
			}

			deposit := &pb.Deposit{}
			if err = testutil.ConvertToPb(tt.Deposit, deposit); err != nil {
				t.Fatal(err)
			}

			expectedPost := &pb.BeaconState{}
			if err = testutil.ConvertToPb(tt.Post, expectedPost); err != nil {
				t.Fatal(err)
			}

			valMap := stateutils.ValidatorIndexMap(preState)
			post, err := blocks.ProcessDeposit(preState, deposit, valMap, true, true)
			// Note: This doesn't test anything worthwhile. It essentially tests
			// that *any* error has occurred, not any specific error.
			if len(expectedPost.ValidatorRegistry) == 0 {
				if err == nil {
					t.Fatal("Did not fail when expected")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(post, expectedPost) {
				t.Log(post)
				t.Log(expectedPost)
				t.Error("Post state does not match expected")
			}
		})
	}
}

func TestDepositMinimalYaml(t *testing.T) {
	runDepositTest(t, "deposit_minimum.yaml")
}

func TestDepositMainnetYaml(t *testing.T) {
	runDepositTest(t, "deposit_mainnet.yaml")
}
