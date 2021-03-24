package components

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os/exec"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	contracts "github.com/prysmaticlabs/prysm/contracts/deposit-contract"
	"github.com/prysmaticlabs/prysm/endtoend/helpers"
	e2e "github.com/prysmaticlabs/prysm/endtoend/params"
	e2etypes "github.com/prysmaticlabs/prysm/endtoend/types"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/featureconfig"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"golang.org/x/sync/errgroup"
)

const depositGasLimit = 4000000

var _ e2etypes.ComponentRunner = (*ValidatorNode)(nil)
var _ e2etypes.ComponentRunner = (*ValidatorNodes)(nil)

// ValidatorNodes represents set of validator nodes.
type ValidatorNodes struct {
	e2etypes.ComponentRunner
	config  *e2etypes.E2EConfig
	started chan struct{}
}

// NewValidatorNodes creates and returns a set of validator nodes.
func NewValidatorNodes(config *e2etypes.E2EConfig) *ValidatorNodes {
	return &ValidatorNodes{
		config:  config,
		started: make(chan struct{}, 1),
	}
}

// Start starts the configured amount of validators, also sending and mining their deposits.
func (s *ValidatorNodes) Start(ctx context.Context) error {
	// Always using genesis count since using anything else would be difficult to test for.
	validatorNum := int(params.BeaconConfig().MinGenesisActiveValidatorCount)
	beaconNodeNum := e2e.TestParams.BeaconNodeCount
	if validatorNum%beaconNodeNum != 0 {
		return errors.New("validator count is not easily divisible by beacon node count")
	}
	validatorsPerNode := validatorNum / beaconNodeNum

	// Create validator nodes.
	nodes := make([]*ValidatorNode, beaconNodeNum)
	for i := 0; i < beaconNodeNum; i++ {
		nodes[i] = NewValidatorNode(s.config, validatorsPerNode, i, validatorsPerNode*i)
	}

	// Start all created validator nodes.
	g, ctx := errgroup.WithContext(ctx)
	for _, node := range nodes {
		node := node
		g.Go(func() error {
			return node.Start(ctx)
		})
	}

	// Mark set as ready (happens when all contained nodes report as started).
	go func() {
		for _, node := range nodes {
			select {
			case <-ctx.Done():
				return
			case <-node.Started():
				continue
			}
		}
		// Only close if all nodes are started. When context is cancelled, node is not considered as started.
		// Stalled nodes do not pose problem -- after certain period of time main routine times out.
		close(s.started)
	}()

	return g.Wait()
}

// Started checks whether validator node set is started and all nodes are ready to be queried.
func (s *ValidatorNodes) Started() <-chan struct{} {
	return s.started
}

// ValidatorNode represents a validator node.
type ValidatorNode struct {
	e2etypes.ComponentRunner
	config       *e2etypes.E2EConfig
	started      chan struct{}
	validatorNum int
	index        int
	offset       int
}

// NewValidatorNode creates and returns a validator node.
func NewValidatorNode(config *e2etypes.E2EConfig, validatorNum, index, offset int) *ValidatorNode {
	return &ValidatorNode{
		config:       config,
		validatorNum: validatorNum,
		index:        index,
		offset:       offset,
		started:      make(chan struct{}, 1),
	}
}

// StartNewValidatorNode starts a validator client.
func (v *ValidatorNode) Start(ctx context.Context) error {
	binaryPath, found := bazel.FindBinary("cmd/validator", "validator")
	if !found {
		return errors.New("validator binary not found")
	}

	config, validatorNum, index, offset := v.config, v.validatorNum, v.index, v.offset
	beaconRPCPort := e2e.TestParams.BeaconNodeRPCPort + index
	if beaconRPCPort >= e2e.TestParams.BeaconNodeRPCPort+e2e.TestParams.BeaconNodeCount {
		// Point any extra validator clients to a node we know is running.
		beaconRPCPort = e2e.TestParams.BeaconNodeRPCPort
	}

	file, err := helpers.DeleteAndCreateFile(e2e.TestParams.LogPath, fmt.Sprintf(e2e.ValidatorLogFileName, index))
	if err != nil {
		return err
	}
	gFile, err := helpers.GraffitiYamlFile(e2e.TestParams.TestPath)
	if err != nil {
		return err
	}
	args := []string{
		fmt.Sprintf("--datadir=%s/eth2-val-%d", e2e.TestParams.TestPath, index),
		fmt.Sprintf("--log-file=%s", file.Name()),
		fmt.Sprintf("--graffiti-file=%s", gFile),
		fmt.Sprintf("--interop-num-validators=%d", validatorNum),
		fmt.Sprintf("--interop-start-index=%d", offset),
		fmt.Sprintf("--monitoring-port=%d", e2e.TestParams.ValidatorMetricsPort+index),
		fmt.Sprintf("--grpc-gateway-port=%d", e2e.TestParams.ValidatorGatewayPort+index),
		fmt.Sprintf("--beacon-rpc-provider=localhost:%d", beaconRPCPort),
		"--grpc-headers=dummy=value,foo=bar", // Sending random headers shouldn't break anything.
		"--force-clear-db",
		"--e2e-config",
		"--accept-terms-of-use",
		"--verbosity=debug",
	}
	args = append(args, featureconfig.E2EValidatorFlags...)
	args = append(args, config.ValidatorFlags...)

	cmd := exec.CommandContext(ctx, binaryPath, args...)
	log.Infof("Starting validator client %d with flags: %s", index, strings.Join(args[2:], " "))
	if err = cmd.Start(); err != nil {
		return err
	}

	// Mark node as ready.
	close(v.started)

	return cmd.Wait()
}

// Started checks whether validator node is started and ready to be queried.
func (v *ValidatorNode) Started() <-chan struct{} {
	return v.started
}

// StartValidatorClients starts the configured amount of validators, also sending and mining their validator deposits.
// Should only be used on initialization.
// Deprecated: this method will be removed once ValidatorNodes component is used.
func StartValidatorClients(t *testing.T, config *e2etypes.E2EConfig) {
	// Always using genesis count since using anything else would be difficult to test for.
	validatorNum := int(params.BeaconConfig().MinGenesisActiveValidatorCount)
	beaconNodeNum := e2e.TestParams.BeaconNodeCount
	if validatorNum%beaconNodeNum != 0 {
		t.Fatal("Validator count is not easily divisible by beacon node count.")
	}
	validatorsPerNode := validatorNum / beaconNodeNum
	for i := 0; i < beaconNodeNum; i++ {
		go StartNewValidatorClient(t, config, validatorsPerNode, i, validatorsPerNode*i)
	}
}

// StartNewValidatorClient starts a validator client with the passed in configuration.
// Deprecated: this method will be removed once ValidatorNode component is used.
func StartNewValidatorClient(t *testing.T, config *e2etypes.E2EConfig, validatorNum, index, offset int) {
	binaryPath, found := bazel.FindBinary("cmd/validator", "validator")
	if !found {
		t.Fatal("validator binary not found")
	}

	beaconRPCPort := e2e.TestParams.BeaconNodeRPCPort + index
	if beaconRPCPort >= e2e.TestParams.BeaconNodeRPCPort+e2e.TestParams.BeaconNodeCount {
		// Point any extra validator clients to a node we know is running.
		beaconRPCPort = e2e.TestParams.BeaconNodeRPCPort
	}

	file, err := helpers.DeleteAndCreateFile(e2e.TestParams.LogPath, fmt.Sprintf(e2e.ValidatorLogFileName, index))
	if err != nil {
		t.Fatal(err)
	}
	gFile, err := helpers.GraffitiYamlFile(e2e.TestParams.TestPath)
	if err != nil {
		t.Fatal(err)
	}
	args := []string{
		fmt.Sprintf("--datadir=%s/eth2-val-%d", e2e.TestParams.TestPath, index),
		fmt.Sprintf("--log-file=%s", file.Name()),
		fmt.Sprintf("--graffiti-file=%s", gFile),
		fmt.Sprintf("--interop-num-validators=%d", validatorNum),
		fmt.Sprintf("--interop-start-index=%d", offset),
		fmt.Sprintf("--monitoring-port=%d", e2e.TestParams.ValidatorMetricsPort+index),
		fmt.Sprintf("--grpc-gateway-port=%d", e2e.TestParams.ValidatorGatewayPort+index),
		fmt.Sprintf("--beacon-rpc-provider=localhost:%d", beaconRPCPort),
		"--grpc-headers=dummy=value,foo=bar", // Sending random headers shouldn't break anything.
		"--force-clear-db",
		"--e2e-config",
		"--accept-terms-of-use",
		"--verbosity=debug",
	}
	args = append(args, featureconfig.E2EValidatorFlags...)
	args = append(args, config.ValidatorFlags...)

	cmd := exec.Command(binaryPath, args...)
	t.Logf("Starting validator client %d with flags: %s", index, strings.Join(args[2:], " "))
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
}

// SendAndMineDeposits sends the requested amount of deposits and mines the chain after to ensure the deposits are seen.
func SendAndMineDeposits(t *testing.T, keystorePath string, validatorNum, offset int, partial bool) {
	client, err := rpc.DialHTTP(fmt.Sprintf("http://127.0.0.1:%d", e2e.TestParams.Eth1RPCPort))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	web3 := ethclient.NewClient(client)

	keystoreBytes, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		t.Fatal(err)
	}
	if err = sendDeposits(web3, keystoreBytes, validatorNum, offset, partial); err != nil {
		t.Fatal(err)
	}
	mineKey, err := keystore.DecryptKey(keystoreBytes, "" /*password*/)
	if err != nil {
		t.Fatal(err)
	}
	if err = mineBlocks(web3, mineKey, params.BeaconConfig().Eth1FollowDistance); err != nil {
		t.Fatalf("failed to mine blocks %v", err)
	}
}

// sendDeposits uses the passed in web3 and keystore bytes to send the requested deposits.
func sendDeposits(web3 *ethclient.Client, keystoreBytes []byte, num, offset int, partial bool) error {
	txOps, err := bind.NewTransactor(bytes.NewReader(keystoreBytes), "" /*password*/)
	if err != nil {
		return err
	}
	txOps.GasLimit = depositGasLimit
	nonce, err := web3.PendingNonceAt(context.Background(), txOps.From)
	if err != nil {
		return err
	}
	txOps.Nonce = big.NewInt(int64(nonce))

	contract, err := contracts.NewDepositContract(e2e.TestParams.ContractAddress, web3)
	if err != nil {
		return err
	}

	balances := make([]uint64, num+offset)
	for i := 0; i < len(balances); i++ {
		if i < len(balances)/2 && partial {
			balances[i] = params.BeaconConfig().MaxEffectiveBalance / 2
		} else {
			balances[i] = params.BeaconConfig().MaxEffectiveBalance
		}
	}
	deposits, trie, err := testutil.DepositsWithBalance(balances)
	if err != nil {
		return err
	}
	allDeposits := deposits
	allRoots := trie.Items()
	allBalances := balances
	if partial {
		deposits2, trie2, err := testutil.DepositsWithBalance(balances)
		if err != nil {
			return err
		}
		allDeposits = append(deposits, deposits2[:len(balances)/2]...)
		allRoots = append(trie.Items(), trie2.Items()[:len(balances)/2]...)
		allBalances = append(balances, balances[:len(balances)/2]...)
	}
	for index, dd := range allDeposits {
		if index < offset {
			continue
		}
		depositInGwei := big.NewInt(int64(allBalances[index]))
		txOps.Value = depositInGwei.Mul(depositInGwei, big.NewInt(int64(params.BeaconConfig().GweiPerEth)))
		_, err = contract.Deposit(txOps, dd.Data.PublicKey, dd.Data.WithdrawalCredentials, dd.Data.Signature, bytesutil.ToBytes32(allRoots[index]))
		if err != nil {
			return errors.Wrap(err, "unable to send transaction to contract")
		}
		txOps.Nonce = txOps.Nonce.Add(txOps.Nonce, big.NewInt(1))
	}
	return nil
}
