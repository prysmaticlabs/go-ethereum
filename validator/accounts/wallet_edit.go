package accounts

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/validator/accounts/prompt"
	"github.com/prysmaticlabs/prysm/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/validator/keymanager"
	"github.com/prysmaticlabs/prysm/validator/keymanager/remote"
	"github.com/urfave/cli/v2"
)

// EditWalletConfigurationCli for a user's on-disk wallet, being able to change
// things such as remote gRPC credentials for remote signing, derivation paths
// for HD wallets, and more.
func EditWalletConfigurationCli(cliCtx *cli.Context) error {
	w, err := wallet.OpenWalletOrElseCli(cliCtx, func(cliCtx *cli.Context) (*wallet.Wallet, error) {
		return nil, errors.New(
			"no wallet found, no configuration to edit",
		)
	})
	if err != nil {
		return errors.Wrap(err, "could not open wallet")
	}
	switch w.KeymanagerKind() {
	case keymanager.Imported:
		return errors.New("not possible to edit imported keymanager configuration")
	case keymanager.Derived:
		return errors.New("derived keymanager is not yet supported")
	case keymanager.Remote:
		enc, err := w.ReadKeymanagerConfigFromDisk(cliCtx.Context)
		if err != nil {
			return errors.Wrap(err, "could not read config")
		}
		opts, err := remote.UnmarshalOptionsFile(enc)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal config")
		}
		log.Info("Current configuration")
		// Prints the current configuration to stdout.
		fmt.Println(opts)
		newCfg, err := prompt.InputRemoteKeymanagerConfig(cliCtx)
		if err != nil {
			return errors.Wrap(err, "could not get keymanager config")
		}
		encodedCfg, err := remote.MarshalOptionsFile(cliCtx.Context, newCfg)
		if err != nil {
			return errors.Wrap(err, "could not marshal config file")
		}
		if err := w.WriteKeymanagerConfigToDisk(cliCtx.Context, encodedCfg); err != nil {
			return errors.Wrap(err, "could not write config to disk")
		}
	default:
		return fmt.Errorf("keymanager type %s is not supported", w.KeymanagerKind())
	}
	return nil
}
