package accounts

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/keystore"
	"go.opencensus.io/trace"
)

type ValidatorStatusMetadata struct {
	PublicKey []byte
	Metadata  *ethpb.ValidatorStatusResponse
}

// Maximum grpc requests allowed to fetch account statuses.
const MaxRequestLimit = 3

// XXX: This is an arbitrary number. Should compute
// max keys allowed before exceeding GrpcMaxCallRecvMsgSizeFlag.
const MaxRequestKeys = 13

// FetchAccountStatuses fetches validator statuses from the BeaconNodeValidatorClient
// for each validator public key.
func FetchAccountStatuses(
	ctx context.Context,
	beaconNodeRPCProvider ethpb.BeaconNodeValidatorClient,
	pubkeys [][]byte) ([][]ValidatorStatusMetadata, error) {
	ctx, span := trace.StartSpan(ctx, "validator.FetchAccountStatuses")
	defer span.End()

	errorChannel := make(chan error, MaxRequestLimit)
	responseChannel := make(chan *ethpb.MultipleValidatorStatusResponse, MaxRequestLimit)
	// Fetches statuses in batches.
	i, numBatches := 0, 0
	for ; i+MaxRequestKeys < len(pubkeys); i += MaxRequestKeys {
		go fetchValidatorStatus(
			ctx, beaconNodeRPCProvider, pubkeys[i:i+MaxRequestKeys], responseChannel, errorChannel)
		numBatches++
	}
	if i < len(pubkeys) {
		go fetchValidatorStatus(
			ctx, beaconNodeRPCProvider, pubkeys[i:], responseChannel, errorChannel)
		numBatches++
	}
	// Wait from fetch routines to finish.
	var err error
	statuses := make([][]ValidatorStatusMetadata, 0, numBatches)
	for i := 0; i < numBatches; i++ {
		select {
		case resp := <-responseChannel:
			statuses = append(statuses, responseToSortedMetadata(resp))
		case e := <-errorChannel:
			err = e
		}
	}

	return statuses, err
}

func fetchValidatorStatus(
	ctx context.Context,
	rpcProvder ethpb.BeaconNodeValidatorClient,
	pubkeys [][]byte,
	responseChannel chan *ethpb.MultipleValidatorStatusResponse,
	errorChannel chan error) {
	if ctx.Err() == context.Canceled {
		errorChannel <- errors.Wrap(ctx.Err(), "context has been canceled.")
		return
	}

	req := &ethpb.MultipleValidatorStatusRequest{PublicKeys: pubkeys}
	resp, err := rpcProvder.MultipleValidatorStatus(ctx, req)
	if err != nil {
		errorChannel <- errors.Wrap(
			err,
			fmt.Sprintf("could not fetch validator statuses for %v", pubkeys))
		return
	}

	responseChannel <- resp
}

// ExtractPublicKeys extracts only the public keys from the decrypted keys from the keystore.
func ExtractPublicKeys(decryptedKeys map[string]*keystore.Key) [][]byte {
	i := 0
	pubkeys := make([][]byte, len(decryptedKeys))
	for _, key := range decryptedKeys {
		pubkeys[i] = key.PublicKey.Marshal()
		i++
	}
	return pubkeys
}

func responseToSortedMetadata(resp *ethpb.MultipleValidatorStatusResponse) []ValidatorStatusMetadata {
	pubkeys := resp.GetPublicKeys()
	validatorStatuses := make([]ValidatorStatusMetadata, len(pubkeys))
	for i, status := range resp.GetStatuses() {
		validatorStatuses[i] = ValidatorStatusMetadata{
			PublicKey: pubkeys[i],
			Metadata:  status,
		}
	}
	sort.Slice(validatorStatuses, func(i, j int) bool {
		return validatorStatuses[i].Metadata.Status < validatorStatuses[j].Metadata.Status
	})
	return validatorStatuses
}

// MergeStatuses merges k sorted ValidatorStatusMetadata slices to 1.
func MergeStatuses(allStatuses [][]ValidatorStatusMetadata) []ValidatorStatusMetadata {
	if len(allStatuses) == 0 {
		return []ValidatorStatusMetadata{}
	}
	if len(allStatuses) == 1 {
		return allStatuses[0]
	}
	leftHalf := allStatuses[:len(allStatuses)/2]
	rightHalf := allStatuses[len(allStatuses)/2:]
	return mergeTwo(MergeStatuses(leftHalf), MergeStatuses(rightHalf))
}

func mergeTwo(s1, s2 []ValidatorStatusMetadata) []ValidatorStatusMetadata {
	i, j, k := 0, 0, 0
	sortedStatuses := make([]ValidatorStatusMetadata, len(s1)+len(s2))
	for j < len(s1) && k < len(s2) {
		if s1[j].Metadata.Status < s2[k].Metadata.Status {
			sortedStatuses[i] = s1[j]
			j++
		} else {
			sortedStatuses[i] = s2[k]
			k++
		}
		i++
	}
	for j < len(s1) {
		sortedStatuses[i] = s1[j]
		i, j = i+1, j+1
	}
	for k < len(s2) {
		sortedStatuses[i] = s2[k]
		i, k = i+1, k+1
	}
	return sortedStatuses
}
