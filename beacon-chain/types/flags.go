package types

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/urfave/cli"
)

var (
	// Web3ProviderFlag defines a flag for a mainchain RPC endpoint.
	Web3ProviderFlag = cli.StringFlag{
		Name:  "web3provider",
		Usage: "A mainchain web3 provider string endpoint. Can either be an IPC file string or a WebSocket endpoint. Uses WebSockets by default at ws://127.0.0.1:8546. Cannot be an HTTP endpoint.",
		Value: "ws://127.0.0.1:8546",
	}
	// DataDirFlag passes in the data directory for the beacon chain.
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DirectoryString{node.DefaultDataDir()},
	}
)
