package kv

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sort"

	ssz "github.com/ferranbt/fastssz"
	"github.com/golang/snappy"
	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
	slashpb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	bolt "go.etcd.io/bbolt"
	"go.opencensus.io/trace"
)

// LastEpochWrittenForValidator given a validator index returns the latest
// epoch we have recorded the validator attested for.
func (s *Store) LastEpochWrittenForValidators(
	ctx context.Context, validatorIndices []types.ValidatorIndex,
) ([]*slashertypes.AttestedEpochForValidator, error) {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.LastEpochWrittenForValidators")
	defer span.End()
	attestedEpochs := make([]*slashertypes.AttestedEpochForValidator, 0)
	encodedIndices := make([][]byte, len(validatorIndices))
	for i, valIdx := range validatorIndices {
		encodedIndices[i] = encodeValidatorIndex(valIdx)
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(attestedEpochsByValidator)
		for i, encodedIndex := range encodedIndices {
			epochBytes := bkt.Get(encodedIndex)
			if epochBytes != nil {
				var epoch types.Epoch
				if err := epoch.UnmarshalSSZ(epochBytes); err != nil {
					return err
				}
				attestedEpochs = append(attestedEpochs, &slashertypes.AttestedEpochForValidator{
					ValidatorIndex: validatorIndices[i],
					Epoch:          epoch,
				})
			}
		}
		return nil
	})
	return attestedEpochs, err
}

// SaveLastEpochWrittenForValidators updates the latest epoch a slice
// of validator indices has attested to.
func (s *Store) SaveLastEpochWrittenForValidators(
	ctx context.Context, validatorIndices []types.ValidatorIndex, epoch types.Epoch,
) error {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.SaveLastEpochWrittenForValidators")
	defer span.End()
	encodedIndices := make([][]byte, len(validatorIndices))
	for i, valIdx := range validatorIndices {
		encodedIndices[i] = encodeValidatorIndex(valIdx)
	}
	encodedEpoch, err := epoch.MarshalSSZ()
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(attestedEpochsByValidator)
		for _, encodedIndex := range encodedIndices {
			if err = bkt.Put(encodedIndex, encodedEpoch); err != nil {
				return err
			}
		}
		return nil
	})
}

// CheckDoubleAttesterVotes retries any slashable double votes that exist
// for a series of input attestations.
func (s *Store) CheckAttesterDoubleVotes(
	ctx context.Context, attestations []*slashertypes.IndexedAttestationWrapper, historyLength types.Epoch,
) ([]*slashertypes.AttesterDoubleVote, error) {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.CheckAttesterDoubleVotes")
	defer span.End()
	doubleVotes := make([]*slashertypes.AttesterDoubleVote, 0)
	err := s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(attestationRecordsBucket)
		for _, att := range attestations {
			encEpoch := encodeTargetEpoch(att.IndexedAttestation.Data.Target.Epoch, historyLength)
			for _, valIdx := range att.IndexedAttestation.AttestingIndices {
				encIdx := encodeValidatorIndex(types.ValidatorIndex(valIdx))
				key := append(encEpoch, encIdx...)
				encExistingAttRecord := bkt.Get(key)
				if len(encExistingAttRecord) < 32 {
					continue
				}
				existingSigningRoot := bytesutil.ToBytes32(encExistingAttRecord[:32])
				if existingSigningRoot != att.SigningRoot {
					existingAttRecord, err := decodeAttestationRecord(encExistingAttRecord)
					if err != nil {
						return err
					}
					doubleVotes = append(doubleVotes, &slashertypes.AttesterDoubleVote{
						ValidatorIndex:         types.ValidatorIndex(valIdx),
						Target:                 att.IndexedAttestation.Data.Target.Epoch,
						PrevAttestationWrapper: existingAttRecord,
						AttestationWrapper:     att,
					})
				}
			}
		}
		return nil
	})
	return doubleVotes, err
}

// AttestationRecordForValidator given a validator index and a target epoch,
// retrieves an existing attestation record we have stored in the database.
func (s *Store) AttestationRecordForValidator(
	ctx context.Context, validatorIdx types.ValidatorIndex, targetEpoch types.Epoch, historyLength types.Epoch,
) (*slashertypes.IndexedAttestationWrapper, error) {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.AttestationRecordForValidator")
	defer span.End()
	var record *slashertypes.IndexedAttestationWrapper
	encIdx := encodeValidatorIndex(validatorIdx)
	encEpoch := encodeTargetEpoch(targetEpoch, historyLength)
	key := append(encEpoch, encIdx...)
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(attestationRecordsBucket)
		value := bkt.Get(key)
		if value == nil {
			return nil
		}
		decoded, err := decodeAttestationRecord(value)
		if err != nil {
			return err
		}
		record = decoded
		return nil
	})
	return record, err
}

// SaveAttestationRecordsForValidators saves an attestation records for the
// specified validator indices.
func (s *Store) SaveAttestationRecordsForValidators(
	ctx context.Context,
	attestations []*slashertypes.IndexedAttestationWrapper,
	historyLength types.Epoch,
) error {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.SaveAttestationRecordsForValidators")
	defer span.End()
	encodedTargetEpoch := make([][]byte, len(attestations))
	encodedRecords := make([][]byte, len(attestations))
	for i, att := range attestations {
		encEpoch := encodeTargetEpoch(att.IndexedAttestation.Data.Target.Epoch, historyLength)
		value, err := encodeAttestationRecord(att)
		if err != nil {
			return err
		}
		encodedTargetEpoch[i] = encEpoch
		encodedRecords[i] = value
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(attestationRecordsBucket)
		for i, att := range attestations {
			for _, valIdx := range att.IndexedAttestation.AttestingIndices {
				encIdx := encodeValidatorIndex(types.ValidatorIndex(valIdx))
				key := append(encodedTargetEpoch[i], encIdx...)
				if err := bkt.Put(key, encodedRecords[i]); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// LoadSlasherChunks given a chunk kind and a disk keys, retrieves chunks for a validator
// min or max span used by slasher from our database.
func (s *Store) LoadSlasherChunks(
	ctx context.Context, kind slashertypes.ChunkKind, diskKeys [][]byte,
) ([][]uint16, []bool, error) {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.LoadSlasherChunk")
	defer span.End()
	chunks := make([][]uint16, 0)
	var exists []bool
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(slasherChunksBucket)
		for _, diskKey := range diskKeys {
			key := append(ssz.MarshalUint8(make([]byte, 0), uint8(kind)), diskKey...)
			chunkBytes := bkt.Get(key)
			if chunkBytes == nil {
				chunks = append(chunks, []uint16{})
				exists = append(exists, false)
				continue
			}
			chunk, err := decodeSlasherChunk(chunkBytes)
			if err != nil {
				return err
			}
			chunks = append(chunks, chunk)
			exists = append(exists, true)
		}
		return nil
	})
	return chunks, exists, err
}

// SaveSlasherChunk given a chunk kind, list of disk keys, and list of chunks,
// saves the chunks to our database for use by slasher in slashing detection.
func (s *Store) SaveSlasherChunks(
	ctx context.Context, kind slashertypes.ChunkKind, chunkKeys [][]byte, chunks [][]uint16,
) error {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.SaveSlasherChunks")
	defer span.End()
	encodedKeys := make([][]byte, len(chunkKeys))
	encodedChunks := make([][]byte, len(chunkKeys))
	for i := 0; i < len(chunkKeys); i++ {
		encodedKeys[i] = append(ssz.MarshalUint8(make([]byte, 0), uint8(kind)), chunkKeys[i]...)
		encodedChunks[i] = encodeSlasherChunk(chunks[i])
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(slasherChunksBucket)
		for i := 0; i < len(chunkKeys); i++ {
			if err := bkt.Put(encodedKeys[i], encodedChunks[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

// CheckDoubleBlockProposals takes in a list of proposals and for each,
// checks if there already exists a proposal at the same slot+validatorIndex combination. If so,
// We check if the existing signing root is not-empty and is different than the incoming
// proposal signing root. If so, we return a double block proposal object.
func (s *Store) CheckDoubleBlockProposals(
	ctx context.Context, proposals []*slashertypes.SignedBlockHeaderWrapper,
) ([]*ethpb.ProposerSlashing, error) {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.CheckDoubleBlockProposals")
	defer span.End()
	proposerSlashings := make([]*ethpb.ProposerSlashing, 0, len(proposals))
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(proposalRecordsBucket)
		for _, proposal := range proposals {
			key, err := keyForValidatorProposal(proposal)
			if err != nil {
				return err
			}
			encExistingProposalWrapper := bkt.Get(key)
			if len(encExistingProposalWrapper) < 32 {
				continue
			}
			existingSigningRoot := bytesutil.ToBytes32(encExistingProposalWrapper[:32])
			if existingSigningRoot != proposal.SigningRoot {
				existingProposalWrapper, err := decodeProposalRecord(encExistingProposalWrapper)
				if err != nil {
					return err
				}
				proposerSlashings = append(proposerSlashings, &ethpb.ProposerSlashing{
					Header_1: existingProposalWrapper.SignedBeaconBlockHeader,
					Header_2: proposal.SignedBeaconBlockHeader,
				})
			}
		}
		return nil
	})
	return proposerSlashings, err
}

// SaveBlockProposals takes in a list of block proposals and saves them to our
// proposal records bucket in the database.
func (s *Store) SaveBlockProposals(
	ctx context.Context, proposals []*slashertypes.SignedBlockHeaderWrapper,
) error {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.SaveBlockProposals")
	defer span.End()
	encodedKeys := make([][]byte, len(proposals))
	encodedProposals := make([][]byte, len(proposals))
	for i, proposal := range proposals {
		key, err := keyForValidatorProposal(proposal)
		if err != nil {
			return err
		}
		enc, err := encodeProposalRecord(proposal)
		if err != nil {
			return err
		}
		encodedKeys[i] = key
		encodedProposals[i] = enc
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(proposalRecordsBucket)
		for i := range proposals {
			if err := bkt.Put(encodedKeys[i], encodedProposals[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

// PruneProposals prunes all proposal data older than historyLength.
func (s *Store) PruneProposals(ctx context.Context, currentEpoch types.Epoch, historyLength types.Epoch) error {
	if currentEpoch < historyLength {
		return nil
	}
	// + 1 here so we can prune everything less than this, but not equal.
	endPruneSlot, err := helpers.StartSlot(currentEpoch - historyLength)
	if err != nil {
		return err
	}
	endEnc, err := endPruneSlot.MarshalSSZ()
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		proposalBkt := tx.Bucket(proposalRecordsBucket)
		c := proposalBkt.Cursor()
		for k, _ := c.Seek(endEnc); k != nil; k, _ = c.Prev() {
			if !slotPrefixLessThan(k, endEnc) {
				continue
			}
			slasherProposalsPrunedTotal.Inc()
			if err := proposalBkt.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}

// PruneAttestations prunes all proposal data older than historyLength.
func (s *Store) PruneAttestations(ctx context.Context, currentEpoch types.Epoch, historyLength types.Epoch) error {
	if currentEpoch < historyLength {
		return nil
	}
	// + 1 here so we can prune everything less than this, but not equal.
	endPruneEpoch := currentEpoch - types.Epoch(historyLength)
	epochEnc := encodeTargetEpoch(endPruneEpoch, historyLength)
	return s.db.Update(func(tx *bolt.Tx) error {
		attBkt := tx.Bucket(attestationRecordsBucket)
		c := attBkt.Cursor()
		for k, _ := c.Seek(epochEnc); k != nil; k, _ = c.Prev() {
			if !epochPrefixLessThan(k, epochEnc) {
				continue
			}
			slasherAttestationsPrunedTotal.Inc()
			if err := attBkt.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}

// HighestAttestations retrieves the last attestation data from the database for all indices.
func (s *Store) HighestAttestations(
	ctx context.Context,
	indices []types.ValidatorIndex,
) ([]*slashpb.HighestAttestation, error) {
	if len(indices) == 0 {
		return nil, nil
	}
	// Sort indices to keep DB interactions short.
	sort.SliceStable(indices, func(i, j int) bool {
		return uint64(indices[i]) < uint64(indices[j])
	})

	var err error
	encodedIndices := make([][]byte, len(indices))
	for i, valIdx := range indices {
		encodedIndices[i] = encodeValidatorIndex(valIdx)
	}

	history := make([]*slashpb.HighestAttestation, 0, len(encodedIndices))
	err = s.db.View(func(tx *bolt.Tx) error {
		attBkt := tx.Bucket(attestationRecordsBucket)
		for i := 0; i < len(encodedIndices); i++ {
			c := attBkt.Cursor()
			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				if suffixForAttestationRecordsKey(k, encodedIndices[i]) {
					attWrapper, err := decodeAttestationRecord(v)
					if err != nil {
						return err
					}
					highestAtt := &slashpb.HighestAttestation{
						ValidatorIndex:     uint64(indices[i]),
						HighestSourceEpoch: attWrapper.IndexedAttestation.Data.Source.Epoch,
						HighestTargetEpoch: attWrapper.IndexedAttestation.Data.Target.Epoch,
					}
					history = append(history, highestAtt)
					break
				}
			}
		}
		return nil
	})
	return history, err
}

func slotPrefixLessThan(key, lessThan []byte) bool {
	encSlot := key[:8]
	return bytes.Compare(encSlot, lessThan) < 0
}

func epochPrefixLessThan(key, lessThan []byte) bool {
	encSlot := key[:2]
	return bytes.Compare(encSlot, lessThan) < 0
}

func suffixForAttestationRecordsKey(key, encodedValidatorIndex []byte) bool {
	encIdx := key[2:]
	return bytes.Equal(encIdx, encodedValidatorIndex)
}

// Disk key for a validator proposal, including a slot+validatorIndex as a byte slice.
func keyForValidatorProposal(proposal *slashertypes.SignedBlockHeaderWrapper) ([]byte, error) {
	encSlot, err := proposal.SignedBeaconBlockHeader.Header.Slot.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	encValidatorIdx := encodeValidatorIndex(proposal.SignedBeaconBlockHeader.Header.ProposerIndex)
	return append(encSlot, encValidatorIdx...), nil
}

func encodeSlasherChunk(chunk []uint16) []byte {
	val := make([]byte, 0)
	for i := 0; i < len(chunk); i++ {
		val = append(val, ssz.MarshalUint16(make([]byte, 0), chunk[i])...)
	}
	return snappy.Encode(nil, val)
}

func decodeSlasherChunk(enc []byte) ([]uint16, error) {
	chunkBytes, err := snappy.Decode(nil, enc)
	if err != nil {
		return nil, err
	}
	chunk := make([]uint16, 0)
	for i := 0; i < len(chunkBytes); i += 2 {
		distance := ssz.UnmarshallUint16(chunkBytes[i : i+2])
		chunk = append(chunk, distance)
	}
	return chunk, nil
}

// Decode attestation record from bytes.
func encodeAttestationRecord(att *slashertypes.IndexedAttestationWrapper) ([]byte, error) {
	if att == nil || att.IndexedAttestation == nil {
		return []byte{}, errors.New("nil proposal record")
	}
	encodedAtt, err := att.IndexedAttestation.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	compressedAtt := snappy.Encode(nil, encodedAtt)
	return append(att.SigningRoot[:], compressedAtt...), nil
}

// Decode attestation record from bytes.
func decodeAttestationRecord(encoded []byte) (*slashertypes.IndexedAttestationWrapper, error) {
	if len(encoded) < 32 {
		return nil, fmt.Errorf("wrong length for encoded attestation record, want 32, got %d", len(encoded))
	}
	signingRoot := encoded[:32]
	decodedAtt := &ethpb.IndexedAttestation{}
	decodedAttBytes, err := snappy.Decode(nil, encoded[32:])
	if err != nil {
		return nil, err
	}
	if err := decodedAtt.UnmarshalSSZ(decodedAttBytes); err != nil {
		return nil, err
	}
	return &slashertypes.IndexedAttestationWrapper{
		IndexedAttestation: decodedAtt,
		SigningRoot:        bytesutil.ToBytes32(signingRoot),
	}, nil
}

func encodeProposalRecord(blkHdr *slashertypes.SignedBlockHeaderWrapper) ([]byte, error) {
	if blkHdr == nil || blkHdr.SignedBeaconBlockHeader == nil {
		return []byte{}, errors.New("nil proposal record")
	}
	encodedHdr, err := blkHdr.SignedBeaconBlockHeader.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	compressedHdr := snappy.Encode(nil, encodedHdr)
	return append(blkHdr.SigningRoot[:], compressedHdr...), nil
}

func decodeProposalRecord(encoded []byte) (*slashertypes.SignedBlockHeaderWrapper, error) {
	if len(encoded) < 32 {
		return nil, fmt.Errorf("wrong length for encoded proposal record, want 32, got %d", len(encoded))
	}
	signingRoot := encoded[:32]
	decodedBlkHdr := &ethpb.SignedBeaconBlockHeader{}
	decodedHdrBytes, err := snappy.Decode(nil, encoded[32:])
	if err != nil {
		return nil, err
	}
	if err := decodedBlkHdr.UnmarshalSSZ(decodedHdrBytes); err != nil {
		return nil, err
	}
	return &slashertypes.SignedBlockHeaderWrapper{
		SignedBeaconBlockHeader: decodedBlkHdr,
		SigningRoot:             bytesutil.ToBytes32(signingRoot),
	}, nil
}

// Encodes an epoch by performing modulo HISTORY_SIZE from slasher using 2 bytes instead of 8 as a
// client optimization to save space in the database.
func encodeTargetEpoch(epoch types.Epoch, historyLength types.Epoch) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(epoch%historyLength))
	return buf
}

// Encodes a validator index using 5 bytes instead of 8 as a
// client optimization to save space in the database. Because the max validator
// registry size is 2**40, this is a safe optimization.
func encodeValidatorIndex(index types.ValidatorIndex) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(index))
	return buf[:5]
}
