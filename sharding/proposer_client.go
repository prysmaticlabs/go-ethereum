package sharding

//go:generate abigen --sol contracts/sharding_manager.sol --pkg contracts --out contracts/sharding_manager.go

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/sharding/contracts"
	cli "gopkg.in/urfave/cli.v1"
)

// ProposerClient communicates to Geth node via JSON RPC.
type ProposerClient struct {
	endpoint string             // Endpoint to JSON RPC
	client   *ethclient.Client  // Ethereum RPC client.
	keystore *keystore.KeyStore // Keystore containing the single signer
	ctx      *cli.Context       // Command line context
	smc      *contracts.SMC     // The deployed sharding management contract
}

// MakeProposerClient for interfacing with Geth full node.
func MakeProposerClient(ctx *cli.Context) *ProposerClient {
	path := node.DefaultDataDir()
	if ctx.GlobalIsSet(utils.DataDirFlag.Name) {
		path = ctx.GlobalString(utils.DataDirFlag.Name)
	}

	endpoint := ctx.Args().First()
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s/%s.ipc", path, clientIdentifier)
	}
	if ctx.GlobalIsSet(utils.IPCPathFlag.Name) {
		endpoint = ctx.GlobalString(utils.IPCPathFlag.Name)
	}

	config := &node.Config{
		DataDir: path,
	}

	scryptN, scryptP, keydir, err := config.AccountConfig()
	if err != nil {
		panic(err) // TODO(prestonvanloon): handle this
	}
	ks := keystore.NewKeyStore(keydir, scryptN, scryptP)

	return &ProposerClient{
		endpoint: endpoint,
		keystore: ks,
		ctx:      ctx,
	}
}

// Start the proposer client
func (c *ProposerClient) Start() error {
	log.Info("Starting proposer client")
	rpcClient, err := dialRPC(c.endpoint)
	if err != nil {
		return err
	}
	c.client = ethclient.NewClient(rpcClient)
	defer rpcClient.Close()

	// Check account existence and unlock account before starting collator client
	accounts := c.keystore.Accounts()
	if len(accounts) == 0 {
		return fmt.Errorf("no accounts found")
	}

	if err := c.unlockAccount(accounts[0]); err != nil {
		return fmt.Errorf("cannot unlock account. %v", err)
	}

	if err := subscribePendingTransactions(c); err != nil {
		return err
	}

	return nil
}

// Wait until proposer client is shutdown.
func (c *ProposerClient) Wait() {
	log.Info("Proposer client has been shutdown...")
}

// UnlockAccount will unlock the specified account using utils.PasswordFileFlag or empty string if unset.
func (c *ProposerClient) unlockAccount(account accounts.Account) error {
	pass := ""

	if c.ctx.GlobalIsSet(utils.PasswordFileFlag.Name) {
		file, err := os.Open(c.ctx.GlobalString(utils.PasswordFileFlag.Name))
		if err != nil {
			return fmt.Errorf("unable to open file containing account password %s. %v", utils.PasswordFileFlag.Value, err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			err = scanner.Err()
			if err != nil {
				return fmt.Errorf("unable to read contents of file %v", err)
			}
			return errors.New("password not found in file")
		}

		pass = scanner.Text()
	}

	return c.keystore.Unlock(account, pass)
}

// Account to use for sharding transactions.
func (c *ProposerClient) Account() *accounts.Account {
	accounts := c.keystore.Accounts()

	return &accounts[0]
}

// ChainReader for interacting with the chain.
func (c *ProposerClient) ChainReader() ethereum.ChainReader {
	return ethereum.ChainReader(c.client)
}

// PendingStateEventer provides subscription functionality
func (c *ProposerClient) PendingStateEventer() ethereum.PendingStateEventer {
	return ethereum.PendingStateEventer(c.client)
}

// Client to interact with ethereum node.
func (c *ProposerClient) Client() *ethclient.Client {
	return c.client
}

// Context from the CLI
func (c *ProposerClient) Context() *cli.Context {
	return c.ctx
}

// SMCCaller to interact with the sharding manager contract.
func (c *ProposerClient) SMCCaller() *contracts.SMCCaller {
	return &c.smc.SMCCaller
}
