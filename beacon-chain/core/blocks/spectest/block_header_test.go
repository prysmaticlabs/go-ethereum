package spectest

import (
	"context"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/proto"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params/spectest"
)

// Block header test is actually a full block processing test. Not sure why it
// was named "block_header". The note in the test format readme says "Note that
// block_header is not strictly an operation (and is a full Block), but
// processed in the same manner, and hence included here.". This also tests a
// state transition function, not specifically a block function, but we'll leave
// this test function here to group with the "block operations" in consistent
// manner with the upstream yaml tests.
func runBlockHeaderTest(t *testing.T, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	test := &BlockHeaderTest{}
	if err := yaml.Unmarshal(file, test); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if err := spectest.SetConfig(test.Config); err != nil {
		t.Fatal(err)
	}

	for _, tt := range test.TestCases {
		t.Run(tt.Description, func(t *testing.T) {
			pre := &pb.BeaconState{}
			err := convertToPb(tt.Pre, pre)
			if err != nil {
				t.Fatal(err)
			}

			block := &pb.BeaconBlock{}
			if err := convertToPb(tt.Block, block); err != nil {
				t.Fatal(err)
			}

			post, err := state.ProcessBlock(
				context.Background(),
				pre,
				block,
				state.DefaultConfig(),
			)


			if !reflect.ValueOf(tt.Post).IsValid() {
				// Note: This doesn't test anything worthwhile. It essentially tests
				// that *any* error has occurred, not any specific error.
				if err == nil {
					t.Fatal("did not fail when expected")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			expectedPost := &pb.BeaconState{}
			if err := convertToPb(tt.Post, expectedPost); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(post, expectedPost) {
				t.Fatal("Post state does not match expected")
			}
		})
	}
}

func TestBlockHeaderMinimal(t *testing.T) {
	runBlockHeaderTest(t,"block_header_minimal_formatted.yaml")
}

func TestBlockHeaderMainnet(t *testing.T) {
	runBlockHeaderTest(t,"block_header_mainnet_formatted.yaml")
}
