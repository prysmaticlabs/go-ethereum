// Code generated by yaml_to_go. DO NOT EDIT.
// source: sanity_blocks_mainnet.yaml

package spectest

import ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
import pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"

type BlocksMainnet struct {
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	ForksTimeline string   `json:"forks_timeline"`
	Forks         []string `json:"forks"`
	Config        string   `json:"config"`
	Runner        string   `json:"runner"`
	Handler       string   `json:"handler"`
	TestCases     []struct {
		Description string               `json:"description"`
		Pre         *pb.BeaconState      `json:"pre"`
		Blocks      []*ethpb.BeaconBlock `json:"blocks"`
		Post        *pb.BeaconState      `json:"post"`
	} `json:"test_cases"`
}
