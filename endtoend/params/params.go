// Package params defines all custom parameter configurations
// for running end to end tests.
package params

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/ethereum/go-ethereum/common"
)

// params struct defines the parameters needed for running E2E tests to properly handle test sharding.
type params struct {
	TestPath              string
	LogPath               string
	TestShardIndex        int
	TestID                int
	BeaconNodeCount       int
	Eth1RPCPort           int
	ContractAddress       common.Address
	BootNodePort          int
	BeaconNodeRPCPort     int
	BeaconNodeMetricsPort int
	ValidatorMetricsPort  int
	ValidatorGatewayPort  int
	SlasherRPCPort        int
	SlasherMetricsPort    int
}

// TestParams is the globally accessible var for getting config elements.
var TestParams *params

// BootNodeLogFileName is the file name used for the beacon chain node logs.
var BootNodeLogFileName = "bootnode.log"

// BeaconNodeLogFileName is the file name used for the beacon chain node logs.
var BeaconNodeLogFileName = "beacon-%d.log"

// SlasherLogFileName is the file name used for the slasher client logs.
var SlasherLogFileName = "slasher-%d.log"

// ValidatorLogFileName is the file name used for the validator client logs.
var ValidatorLogFileName = "vals-%d.log"

// StandardBeaconCount is a global constant for the count of beacon nodes of standard E2E tests.
var StandardBeaconCount = 2

// DepositCount is the amount of deposits E2E makes on a separate validator client.
var DepositCount = uint64(64)

// Init initializes the E2E config, properly handling test sharding.
// In order to isolate ports and directories on per test bases, specify unique testID.
func Init(testID, beaconNodeCount int) error {
	testPath := bazel.TestTmpDir()
	logPath, ok := os.LookupEnv("TEST_UNDECLARED_OUTPUTS_DIR")
	if !ok {
		return errors.New("expected TEST_UNDECLARED_OUTPUTS_DIR to be defined")
	}
	testIndexStr, ok := os.LookupEnv("TEST_SHARD_INDEX")
	if !ok {
		testIndexStr = "0"
	}
	testIndex, err := strconv.Atoi(testIndexStr)
	if err != nil {
		return err
	}
	testPath = filepath.Join(testPath, fmt.Sprintf("shard-%d-test-%d", testIndex, testID))

	TestParams = &params{
		TestPath:        testPath,
		LogPath:         logPath,
		TestShardIndex:  testIndex,
		TestID:          testID,
		BeaconNodeCount: beaconNodeCount,
		// Adjusting port numbers, so that test index doesn't conflict with the other node ports.
		Eth1RPCPort:           3100 + testIndex*100 + testID,
		BootNodePort:          4100 + testIndex*100 + testID,
		BeaconNodeRPCPort:     4150 + testIndex*100 + testID,
		BeaconNodeMetricsPort: 5100 + testIndex*100 + testID,
		ValidatorMetricsPort:  6100 + testIndex*100 + testID,
		ValidatorGatewayPort:  7150 + testIndex*100 + testID,
		SlasherRPCPort:        7100 + testIndex*100 + testID,
		SlasherMetricsPort:    8100 + testIndex*100 + testID,
	}
	return nil
}
