package main

import (
	"context"
	"fmt"
	"time"

	pbBeacon "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

var keyInterval = 3 * time.Second

type keyChecker struct {
	db     *db
	client pbBeacon.ValidatorServiceClient
}

func newkeyChecker(db *db, beaconRPCAddr string) *keyChecker {
	// connect to the beacon node
	dialOpt := grpc.WithInsecure()
	conn, err := grpc.DialContext(context.Background(), beaconRPCAddr, dialOpt, grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
	if err != nil {
		log.Errorf("Could not dial endpoint: %s, %v", beaconRPCAddr, err)
	}
	valClient := pbBeacon.NewValidatorServiceClient(conn)

	return &keyChecker{
		db:     db,
		client: valClient,
	}
}

func (k *keyChecker) checkKeys() error {
	pubkeys, keyMap, err := k.db.KeyMap()
	if err != nil {
		return err
	}

	req := &pbBeacon.ExitedValidatorsRequest{
		PublicKeys: pubkeys,
	}

	ctx, cancel := context.WithTimeout(context.Background(), keyInterval)
	defer cancel()
	resp, err := k.client.ExitedValidators(ctx, req)
	if err != nil {
		return err
	}

	log.WithField(
		"resp_keys", len(resp.PublicKeys),
	).WithField(
		"req_keys", len(req.PublicKeys),
	).Debug("Received EXITED key list")

	for _, key := range resp.PublicKeys {
		log.WithField("key", fmt.Sprintf("%#x", key)).Debug("Removing EXITED key")
		kMap := keyMap[bytesutil.ToBytes48(key)]
		if err := k.db.RemovePKFromPod(kMap.podName, kMap.privateKey); err != nil {
			return err
		}
	}
	return nil
}

func (k *keyChecker) run() {
	for {
		time.Sleep(keyInterval)
		if err := k.checkKeys(); err != nil {
			log.WithField("error", err).Error("Failed to check keys")
		}
	}
}
