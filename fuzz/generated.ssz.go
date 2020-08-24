// Code generated by fastssz. DO NOT EDIT.
package fuzz

import (
	ssz "github.com/ferranbt/fastssz"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// MarshalSSZ ssz marshals the InputBlockHeader object
func (i *InputBlockHeader) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(i)
}

// MarshalSSZTo ssz marshals the InputBlockHeader object to a target array
func (i *InputBlockHeader) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(12)

	// Field (0) 'StateID'
	dst = ssz.MarshalUint64(dst, i.StateID)

	// Offset (1) 'Block'
	dst = ssz.WriteOffset(dst, offset)
	if i.Block == nil {
		i.Block = new(ethpb.BeaconBlock)
	}
	offset += i.Block.SizeSSZ()

	// Field (1) 'Block'
	if dst, err = i.Block.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the InputBlockHeader object
func (i *InputBlockHeader) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 12 {
		return ssz.ErrSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'StateID'
	i.StateID = ssz.UnmarshallUint64(buf[0:8])

	// Offset (1) 'Block'
	if o1 = ssz.ReadOffset(buf[8:12]); o1 > size {
		return ssz.ErrOffset
	}

	// Field (1) 'Block'
	{
		buf = tail[o1:]
		if i.Block == nil {
			i.Block = new(ethpb.BeaconBlock)
		}
		if err = i.Block.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the InputBlockHeader object
func (i *InputBlockHeader) SizeSSZ() (size int) {
	size = 12

	// Field (1) 'Block'
	if i.Block == nil {
		i.Block = new(ethpb.BeaconBlock)
	}
	size += i.Block.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the InputBlockHeader object
func (i *InputBlockHeader) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(i)
}

// HashTreeRootWith ssz hashes the InputBlockHeader object with a hasher
func (i *InputBlockHeader) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'StateID'
	hh.PutUint64(i.StateID)

	// Field (1) 'Block'
	if err = i.Block.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the InputAttesterSlashingWrapper object
func (i *InputAttesterSlashingWrapper) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(i)
}

// MarshalSSZTo ssz marshals the InputAttesterSlashingWrapper object to a target array
func (i *InputAttesterSlashingWrapper) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(12)

	// Field (0) 'StateID'
	dst = ssz.MarshalUint64(dst, i.StateID)

	// Offset (1) 'AttesterSlashing'
	dst = ssz.WriteOffset(dst, offset)
	if i.AttesterSlashing == nil {
		i.AttesterSlashing = new(ethpb.AttesterSlashing)
	}
	offset += i.AttesterSlashing.SizeSSZ()

	// Field (1) 'AttesterSlashing'
	if dst, err = i.AttesterSlashing.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the InputAttesterSlashingWrapper object
func (i *InputAttesterSlashingWrapper) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 12 {
		return ssz.ErrSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'StateID'
	i.StateID = ssz.UnmarshallUint64(buf[0:8])

	// Offset (1) 'AttesterSlashing'
	if o1 = ssz.ReadOffset(buf[8:12]); o1 > size {
		return ssz.ErrOffset
	}

	// Field (1) 'AttesterSlashing'
	{
		buf = tail[o1:]
		if i.AttesterSlashing == nil {
			i.AttesterSlashing = new(ethpb.AttesterSlashing)
		}
		if err = i.AttesterSlashing.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the InputAttesterSlashingWrapper object
func (i *InputAttesterSlashingWrapper) SizeSSZ() (size int) {
	size = 12

	// Field (1) 'AttesterSlashing'
	if i.AttesterSlashing == nil {
		i.AttesterSlashing = new(ethpb.AttesterSlashing)
	}
	size += i.AttesterSlashing.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the InputAttesterSlashingWrapper object
func (i *InputAttesterSlashingWrapper) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(i)
}

// HashTreeRootWith ssz hashes the InputAttesterSlashingWrapper object with a hasher
func (i *InputAttesterSlashingWrapper) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'StateID'
	hh.PutUint64(i.StateID)

	// Field (1) 'AttesterSlashing'
	if err = i.AttesterSlashing.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the InputAttestationWrapper object
func (i *InputAttestationWrapper) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(i)
}

// MarshalSSZTo ssz marshals the InputAttestationWrapper object to a target array
func (i *InputAttestationWrapper) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(12)

	// Field (0) 'StateID'
	dst = ssz.MarshalUint64(dst, i.StateID)

	// Offset (1) 'Attestation'
	dst = ssz.WriteOffset(dst, offset)
	if i.Attestation == nil {
		i.Attestation = new(ethpb.Attestation)
	}
	offset += i.Attestation.SizeSSZ()

	// Field (1) 'Attestation'
	if dst, err = i.Attestation.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the InputAttestationWrapper object
func (i *InputAttestationWrapper) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 12 {
		return ssz.ErrSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'StateID'
	i.StateID = ssz.UnmarshallUint64(buf[0:8])

	// Offset (1) 'Attestation'
	if o1 = ssz.ReadOffset(buf[8:12]); o1 > size {
		return ssz.ErrOffset
	}

	// Field (1) 'Attestation'
	{
		buf = tail[o1:]
		if i.Attestation == nil {
			i.Attestation = new(ethpb.Attestation)
		}
		if err = i.Attestation.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the InputAttestationWrapper object
func (i *InputAttestationWrapper) SizeSSZ() (size int) {
	size = 12

	// Field (1) 'Attestation'
	if i.Attestation == nil {
		i.Attestation = new(ethpb.Attestation)
	}
	size += i.Attestation.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the InputAttestationWrapper object
func (i *InputAttestationWrapper) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(i)
}

// HashTreeRootWith ssz hashes the InputAttestationWrapper object with a hasher
func (i *InputAttestationWrapper) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'StateID'
	hh.PutUint64(i.StateID)

	// Field (1) 'Attestation'
	if err = i.Attestation.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the InputDepositWrapper object
func (i *InputDepositWrapper) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(i)
}

// MarshalSSZTo ssz marshals the InputDepositWrapper object to a target array
func (i *InputDepositWrapper) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'StateID'
	dst = ssz.MarshalUint64(dst, i.StateID)

	// Field (1) 'Deposit'
	if i.Deposit == nil {
		i.Deposit = new(ethpb.Deposit)
	}
	if dst, err = i.Deposit.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the InputDepositWrapper object
func (i *InputDepositWrapper) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 1248 {
		return ssz.ErrSize
	}

	// Field (0) 'StateID'
	i.StateID = ssz.UnmarshallUint64(buf[0:8])

	// Field (1) 'Deposit'
	if i.Deposit == nil {
		i.Deposit = new(ethpb.Deposit)
	}
	if err = i.Deposit.UnmarshalSSZ(buf[8:1248]); err != nil {
		return err
	}

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the InputDepositWrapper object
func (i *InputDepositWrapper) SizeSSZ() (size int) {
	size = 1248
	return
}

// HashTreeRoot ssz hashes the InputDepositWrapper object
func (i *InputDepositWrapper) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(i)
}

// HashTreeRootWith ssz hashes the InputDepositWrapper object with a hasher
func (i *InputDepositWrapper) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'StateID'
	hh.PutUint64(i.StateID)

	// Field (1) 'Deposit'
	if err = i.Deposit.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the InputVoluntaryExitWrapper object
func (i *InputVoluntaryExitWrapper) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(i)
}

// MarshalSSZTo ssz marshals the InputVoluntaryExitWrapper object to a target array
func (i *InputVoluntaryExitWrapper) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'StateID'
	dst = ssz.MarshalUint64(dst, i.StateID)

	// Field (1) 'VoluntaryExit'
	if i.VoluntaryExit == nil {
		i.VoluntaryExit = new(ethpb.VoluntaryExit)
	}
	if dst, err = i.VoluntaryExit.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the InputVoluntaryExitWrapper object
func (i *InputVoluntaryExitWrapper) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 24 {
		return ssz.ErrSize
	}

	// Field (0) 'StateID'
	i.StateID = ssz.UnmarshallUint64(buf[0:8])

	// Field (1) 'VoluntaryExit'
	if i.VoluntaryExit == nil {
		i.VoluntaryExit = new(ethpb.VoluntaryExit)
	}
	if err = i.VoluntaryExit.UnmarshalSSZ(buf[8:24]); err != nil {
		return err
	}

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the InputVoluntaryExitWrapper object
func (i *InputVoluntaryExitWrapper) SizeSSZ() (size int) {
	size = 24
	return
}

// HashTreeRoot ssz hashes the InputVoluntaryExitWrapper object
func (i *InputVoluntaryExitWrapper) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(i)
}

// HashTreeRootWith ssz hashes the InputVoluntaryExitWrapper object with a hasher
func (i *InputVoluntaryExitWrapper) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'StateID'
	hh.PutUint64(i.StateID)

	// Field (1) 'VoluntaryExit'
	if err = i.VoluntaryExit.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the InputProposerSlashingWrapper object
func (i *InputProposerSlashingWrapper) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(i)
}

// MarshalSSZTo ssz marshals the InputProposerSlashingWrapper object to a target array
func (i *InputProposerSlashingWrapper) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'StateID'
	dst = ssz.MarshalUint64(dst, i.StateID)

	// Field (1) 'ProposerSlashing'
	if i.ProposerSlashing == nil {
		i.ProposerSlashing = new(ethpb.ProposerSlashing)
	}
	if dst, err = i.ProposerSlashing.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the InputProposerSlashingWrapper object
func (i *InputProposerSlashingWrapper) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 424 {
		return ssz.ErrSize
	}

	// Field (0) 'StateID'
	i.StateID = ssz.UnmarshallUint64(buf[0:8])

	// Field (1) 'ProposerSlashing'
	if i.ProposerSlashing == nil {
		i.ProposerSlashing = new(ethpb.ProposerSlashing)
	}
	if err = i.ProposerSlashing.UnmarshalSSZ(buf[8:424]); err != nil {
		return err
	}

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the InputProposerSlashingWrapper object
func (i *InputProposerSlashingWrapper) SizeSSZ() (size int) {
	size = 424
	return
}

// HashTreeRoot ssz hashes the InputProposerSlashingWrapper object
func (i *InputProposerSlashingWrapper) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(i)
}

// HashTreeRootWith ssz hashes the InputProposerSlashingWrapper object with a hasher
func (i *InputProposerSlashingWrapper) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'StateID'
	hh.PutUint64(i.StateID)

	// Field (1) 'ProposerSlashing'
	if err = i.ProposerSlashing.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}
