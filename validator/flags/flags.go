package flags

import (
	"github.com/prysmaticlabs/prysm/shared/cmd"
	"github.com/urfave/cli"
)

var (
	// NoCustomConfigFlag determines whether to launch a beacon chain using real parameters or demo parameters.
	NoCustomConfigFlag = cli.BoolFlag{
		Name:  "no-custom-config",
		Usage: "Run the beacon chain with the real parameters from phase 0.",
	}
	// BeaconRPCProviderFlag defines a beacon node RPC endpoint.
	BeaconRPCProviderFlag = cli.StringFlag{
		Name:  "beacon-rpc-provider",
		Usage: "Beacon node RPC provider endpoint",
		Value: "localhost:4000",
	}
	// CertFlag defines a flag for the node's TLS certificate.
	CertFlag = cli.StringFlag{
		Name:  "tls-cert",
		Usage: "Certificate for secure gRPC. Pass this and the tls-key flag in order to use gRPC securely.",
	}
	// KeystorePathFlag defines the location of the keystore directory for a validator's account.
	KeystorePathFlag = cmd.DirectoryFlag{
		Name:  "keystore-path",
		Usage: "Path to the desired keystore directory",
		Value: cmd.DirectoryString{Value: ""},
	}
	// UnencryptedKeysFlag specifies a file path of a JSON file of unencrypted validator keys as an
	// alternative from launching the validator client from decrypting a keystore directory.
	UnencryptedKeysFlag = cli.StringFlag{
		Name:  "unencrypted-keys",
		Usage: "Filepath to a JSON file of unencrypted validator keys for easier launching of the validator client",
		Value: "",
	}
	// KeyManager specifies the key manager to use.
	KeyManager = cli.StringFlag{
		Name:  "keymanager",
		Usage: "The keymanger to use (unencrypted, interop, keystore, wallet)",
		Value: "",
	}
	// KeyManagerOpts specifies the key manager options.
	KeyManagerOpts = cli.StringFlag{
		Name:  "keymanageropts",
		Usage: "The options for the keymanger, either a JSON string or path to same",
		Value: "",
	}
	// PasswordFlag defines the password value for storing and retrieving validator private keys from the keystore.
	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "String value of the password for your validator private keys",
	}
	// DisablePenaltyRewardLogFlag defines the ability to not log reward/penalty information during deployment
	DisablePenaltyRewardLogFlag = cli.BoolFlag{
		Name:  "disable-rewards-penalties-logging",
		Usage: "Disable reward/penalty logging during cluster deployment",
	}
	// GraffitiFlag defines the graffiti value included in proposed blocks
	GraffitiFlag = cli.StringFlag{
		Name:  "graffiti",
		Usage: "String to include in proposed blocks",
	}
	// GrpcMaxCallRecvMsgSizeFlag defines the max call message size for GRPC
	GrpcMaxCallRecvMsgSizeFlag = cli.IntFlag{
		Name:  "grpc-max-msg-size",
		Usage: "Integer to define max recieve message call size (default: 52428800 (for 50Mb)).",
	}
	// AccountMetricsFlag defines the graffiti value included in proposed blocks, default false.
	AccountMetricsFlag = cli.BoolFlag{
		Name:  "enable-account-metrics",
		Usage: "Enable prometheus metrics for validator accounts",
	}
)
