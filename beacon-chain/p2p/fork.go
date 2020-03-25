package p2p

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/sirupsen/logrus"
)

// ENR key used for eth2-related fork data.
const eth2ENRKey = "eth2"

// ENRForkID represents a special value ssz-encoded into
// the local node's ENR for discovery purposes. Peers should
// only connect if their enrForkID matches.
type ENRForkID struct {
	CurrentForkDigest []byte `ssz-size:"4"`
	NextForkVersion   []byte `ssz-size:"4"`
	NextForkEpoch     uint64
}

// Compares fork ENRs between an incoming peer's record and our node's
// local record values for current and next fork version/epoch.
func (s *Service) compareForkENR(record *enr.Record) error {
	currentRecord := s.dv5Listener.LocalNode().Node().Record()
	peerForkENR, err := retrieveForkEntry(record)
	if err != nil {
		return err
	}
	currentForkENR, err := retrieveForkEntry(currentRecord)
	if err != nil {
		return err
	}
	// Clients SHOULD connect to peers with current_fork_digest, next_fork_version,
	// and next_fork_epoch that match local values.
	if !bytes.Equal(peerForkENR.CurrentForkDigest, currentForkENR.CurrentForkDigest) {
		return fmt.Errorf(
			"fork digest of peer with ENR %v: %v, does not match local value: %v",
			record,
			peerForkENR.CurrentForkDigest,
			currentForkENR.CurrentForkDigest,
		)
	}
	// Clients MAY connect to peers with the same current_fork_version but a
	// different next_fork_version/next_fork_epoch. Unless ENRForkID is manually
	// updated to matching prior to the earlier next_fork_epoch of the two clients,
	// these type of connecting clients will be unable to successfully interact
	// starting at the earlier next_fork_epoch.
	buf := bytes.NewBuffer([]byte{})
	if err := record.EncodeRLP(buf); err != nil {
		return errors.Wrap(err, "could not encode ENR record to bytes")
	}
	enrString := base64.URLEncoding.EncodeToString(buf.Bytes())
	if peerForkENR.NextForkEpoch != currentForkENR.NextForkEpoch {
		log.WithFields(logrus.Fields{
			"peerNextForkEpoch": peerForkENR.NextForkEpoch,
			"peerENR":           enrString,
		}).Debug("Peer matches fork digest but has different next fork epoch")
	}
	if !bytes.Equal(peerForkENR.NextForkVersion, currentForkENR.NextForkVersion) {
		log.WithFields(logrus.Fields{
			"peerNextForkVersion": peerForkENR.NextForkVersion,
			"peerENR":             enrString,
		}).Debug("Peer matches fork digest but has different next fork version")
	}
	return nil
}

// Adds a fork entry as an ENR record under the eth2EnrKey for
// the local node. The fork entry is an ssz-encoded enrForkID type
// which takes into account the current fork version from the current
// epoch to create a fork digest, the next fork version,
// and the next fork epoch.
func addForkEntry(
	node *enode.LocalNode,
	genesisTime time.Time,
	genesisValidatorsRoot []byte,
) (*enode.LocalNode, error) {
	currentSlot := helpers.SlotsSince(genesisTime)
	currentEpoch := helpers.SlotToEpoch(currentSlot)

	// We retrieve a list of scheduled forks by epoch.
	// We loop through the keys in this map to determine the current
	// fork version based on the current, time-based epoch number
	// since the genesis time.
	currentForkVersion := params.BeaconConfig().GenesisForkVersion
	scheduledForks := params.BeaconConfig().ForkVersionSchedule
	for epoch, forkVersion := range scheduledForks {
		if epoch <= currentEpoch {
			currentForkVersion = forkVersion
		}
	}

	digest, err := helpers.ComputeForkDigest(currentForkVersion, genesisValidatorsRoot)
	if err != nil {
		return nil, err
	}
	nextForkEpoch := params.BeaconConfig().NextForkEpoch
	enrForkID := &ENRForkID{
		CurrentForkDigest: digest[:],
		NextForkVersion:   params.BeaconConfig().NextForkVersion,
		NextForkEpoch:     nextForkEpoch,
	}
	enc, err := ssz.Marshal(enrForkID)
	if err != nil {
		return nil, err
	}
	forkEntry := enr.WithEntry(eth2ENRKey, enc)
	node.Set(forkEntry)
	return node, nil
}

// Retrieves an enrForkID from an ENR record by key lookup
// under the eth2EnrKey.
func retrieveForkEntry(record *enr.Record) (*ENRForkID, error) {
	sszEncodedForkEntry := make([]byte, 16)
	entry := enr.WithEntry(eth2ENRKey, &sszEncodedForkEntry)
	err := record.Load(entry)
	if err != nil {
		return nil, err
	}
	forkEntry := &ENRForkID{}
	if err := ssz.Unmarshal(sszEncodedForkEntry, forkEntry); err != nil {
		return nil, err
	}
	return forkEntry, nil
}
