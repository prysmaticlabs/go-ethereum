package stateV0

import (
	"fmt"

	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// RotateAttestations sets the previous epoch attestations to the current epoch attestations and
// then clears the current epoch attestations.
func (b *BeaconState) RotateAttestations() error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.setPreviousEpochAttestations(b.currentEpochAttestations())
	b.setCurrentEpochAttestations([]*pbp2p.PendingAttestation{})
	return nil
}

func (b *BeaconState) setPreviousEpochAttestations(val []*pbp2p.PendingAttestation) {
	b.sharedFieldReferences[previousEpochAttestations].MinusRef()
	b.sharedFieldReferences[previousEpochAttestations] = stateutil.NewRef(1)

	b.state.PreviousEpochAttestations = val
	b.markFieldAsDirty(previousEpochAttestations)
	b.rebuildTrie[previousEpochAttestations] = true
}

func (b *BeaconState) setCurrentEpochAttestations(val []*pbp2p.PendingAttestation) {
	b.sharedFieldReferences[currentEpochAttestations].MinusRef()
	b.sharedFieldReferences[currentEpochAttestations] = stateutil.NewRef(1)

	b.state.CurrentEpochAttestations = val
	b.markFieldAsDirty(currentEpochAttestations)
	b.rebuildTrie[currentEpochAttestations] = true
}

// AppendCurrentEpochAttestations for the beacon state. Appends the new value
// to the the end of list.
func (b *BeaconState) AppendCurrentEpochAttestations(val *pbp2p.PendingAttestation) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	atts := b.state.CurrentEpochAttestations
	max := uint64(params.BeaconConfig().SlotsPerEpoch) * params.BeaconConfig().MaxAttestations
	if uint64(len(atts)) >= max {
		return fmt.Errorf("current pending attestation exceeds max length %d", max)
	}

	if b.sharedFieldReferences[currentEpochAttestations].Refs() > 1 {
		// Copy elements in underlying array by reference.
		atts = make([]*pbp2p.PendingAttestation, len(b.state.CurrentEpochAttestations))
		copy(atts, b.state.CurrentEpochAttestations)
		b.sharedFieldReferences[currentEpochAttestations].MinusRef()
		b.sharedFieldReferences[currentEpochAttestations] = stateutil.NewRef(1)
	}

	b.state.CurrentEpochAttestations = append(atts, val)
	b.markFieldAsDirty(currentEpochAttestations)
	b.addDirtyIndices(currentEpochAttestations, []uint64{uint64(len(b.state.CurrentEpochAttestations) - 1)})
	return nil
}

// AppendPreviousEpochAttestations for the beacon state. Appends the new value
// to the the end of list.
func (b *BeaconState) AppendPreviousEpochAttestations(val *pbp2p.PendingAttestation) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	atts := b.state.PreviousEpochAttestations
	max := uint64(params.BeaconConfig().SlotsPerEpoch) * params.BeaconConfig().MaxAttestations
	if uint64(len(atts)) >= max {
		return fmt.Errorf("previous pending attestation exceeds max length %d", max)
	}

	if b.sharedFieldReferences[previousEpochAttestations].Refs() > 1 {
		atts = make([]*pbp2p.PendingAttestation, len(b.state.PreviousEpochAttestations))
		copy(atts, b.state.PreviousEpochAttestations)
		b.sharedFieldReferences[previousEpochAttestations].MinusRef()
		b.sharedFieldReferences[previousEpochAttestations] = stateutil.NewRef(1)
	}

	b.state.PreviousEpochAttestations = append(atts, val)
	b.markFieldAsDirty(previousEpochAttestations)
	b.addDirtyIndices(previousEpochAttestations, []uint64{uint64(len(b.state.PreviousEpochAttestations) - 1)})
	return nil
}
