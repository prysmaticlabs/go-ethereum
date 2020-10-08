package kv

import (
	"context"
	"testing"

	slashpb "github.com/prysmaticlabs/prysm/proto/slashing"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestNewAttestationHistoryArray(t *testing.T) {
	ba := newAttestationHistoryArray(0)
	assert.Equal(t, latestEpochWrittenSize+historySize, len(ba))
	ba = newAttestationHistoryArray(params.BeaconConfig().WeakSubjectivityPeriod - 1)
	assert.Equal(t, latestEpochWrittenSize+historySize*params.BeaconConfig().WeakSubjectivityPeriod, uint64(len(ba)))
	ba = newAttestationHistoryArray(params.BeaconConfig().WeakSubjectivityPeriod)
	assert.Equal(t, latestEpochWrittenSize+historySize, len(ba))
	ba = newAttestationHistoryArray(params.BeaconConfig().WeakSubjectivityPeriod + 1)
	assert.Equal(t, latestEpochWrittenSize+historySize+historySize, len(ba))

}

func TestSizeChecks(t *testing.T) {

	require.ErrorContains(t, "is smaller then minimal size", EncHistoryData{}.assertSize())
	require.NoError(t, EncHistoryData{0, 1, 2, 3, 4, 5, 6, 7}.assertSize())
	require.ErrorContains(t, "is not a multiple of entry size", EncHistoryData{0, 1, 2, 3, 4, 5, 6, 7, 8}.assertSize())
	require.NoError(t, newAttestationHistoryArray(0).assertSize())
	require.NoError(t, newAttestationHistoryArray(1).assertSize())
	require.NoError(t, newAttestationHistoryArray(params.BeaconConfig().WeakSubjectivityPeriod).assertSize())
	require.NoError(t, newAttestationHistoryArray(params.BeaconConfig().WeakSubjectivityPeriod-1).assertSize())
}

func TestGetLatestEpochWritten(t *testing.T) {
	ctx := context.Background()
	ha := newAttestationHistoryArray(0)
	ha[0] = 28
	lew, err := ha.getLatestEpochWritten(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(28), lew)
}

func TestSetLatestEpochWritten(t *testing.T) {
	ctx := context.Background()
	ha := newAttestationHistoryArray(0)
	lew, err := ha.setLatestEpochWritten(ctx, 2828282828)
	require.NoError(t, err)
	assert.Equal(t, uint64(2828282828), bytesutil.FromBytes8(lew[:latestEpochWrittenSize]))
}

func TestGetTargetData(t *testing.T) {
	ctx := context.Background()
	ha := newAttestationHistoryArray(0)
	td, err := ha.getTargetData(ctx, 0)
	require.NoError(t, err)
	assert.DeepEqual(t, &HistoryData{
		Source:      0,
		SigningRoot: bytesutil.PadTo([]byte{}, 32),
	}, td)
	_, err = ha.getTargetData(ctx, 1)
	require.ErrorContains(t, "is smaller then the requested target location", err)
}

func TestSetTargetData(t *testing.T) {
	ctx := context.Background()
	type testStruct struct {
		name        string
		enc         EncHistoryData
		target      uint64
		source      uint64
		signingRoot []byte
		expected    EncHistoryData
		error       string
	}
	tests := []testStruct{
		{
			name:        "empty enc",
			enc:         []byte{},
			target:      0,
			source:      100,
			signingRoot: []byte{1, 2, 3},
			expected:    nil,
			error:       "encapsulated data size",
		},
		{
			name:        "new enc",
			enc:         newAttestationHistoryArray(0),
			target:      0,
			source:      100,
			signingRoot: []byte{1, 2, 3},
			expected:    []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			error:       "",
		},
		{
			name:        "higher target",
			enc:         newAttestationHistoryArray(0),
			target:      2,
			source:      100,
			signingRoot: []byte{1, 2, 3},
			expected:    []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			error:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := tt.enc.setTargetData(ctx,
				tt.target,
				&HistoryData{
					Source:      tt.source,
					SigningRoot: tt.signingRoot,
				})
			if tt.error == "" {
				require.NoError(t, err)

			} else {
				assert.ErrorContains(t, tt.error, err)
			}
			require.DeepEqual(t, tt.expected, enc)
		})
	}

}

func TestAttestationHistoryForPubKeysNew_EmptyVals(t *testing.T) {
	pubkeys := [][48]byte{{30}, {25}, {20}}
	db := setupDB(t, pubkeys)

	historyForPubKeys, err := db.AttestationHistoryNewForPubKeys(context.Background(), pubkeys)
	require.NoError(t, err)

	cleanAttHistoryForPubKeys := make(map[[48]byte]EncHistoryData)
	clean := newAttestationHistoryArray(0)
	for _, pubKey := range pubkeys {
		cleanAttHistoryForPubKeys[pubKey] = clean
	}

	require.DeepEqual(t, cleanAttHistoryForPubKeys, historyForPubKeys, "Expected attestation history epoch bits to be empty")
}

func TestAttestationHistoryForPubKeysNew_OK(t *testing.T) {
	ctx := context.Background()
	pubkeys := [][48]byte{{30}, {25}, {20}}
	db := setupDB(t, pubkeys)

	_, err := db.AttestationHistoryNewForPubKeys(context.Background(), pubkeys)
	require.NoError(t, err)

	setAttHistoryForPubKeys := make(map[[48]byte]EncHistoryData)
	clean := newAttestationHistoryArray(0)
	for i, pubKey := range pubkeys {
		enc, err := clean.setTargetData(ctx,
			10,
			&HistoryData{
				Source:      uint64(i),
				SigningRoot: []byte{1, 2, 3},
			})
		require.NoError(t, err)
		setAttHistoryForPubKeys[pubKey] = enc

	}
	err = db.SaveAttestationHistoryNewForPubKeys(context.Background(), setAttHistoryForPubKeys)
	require.NoError(t, err)
	historyForPubKeys, err := db.AttestationHistoryNewForPubKeys(context.Background(), pubkeys)
	require.NoError(t, err)
	require.DeepEqual(t, setAttHistoryForPubKeys, historyForPubKeys, "Expected attestation history epoch bits to be empty")
}

func TestStore_ImportOldAttestationFormat(t *testing.T) {
	ctx := context.Background()
	pubKeys := [][48]byte{{3}, {4}}
	db := setupDB(t, pubKeys)

	farFuture := params.BeaconConfig().FarFutureEpoch
	newMap := make(map[uint64]uint64)
	// The validator attested at target epoch 2 but had no attestations for target epochs 0 and 1.
	newMap[0] = farFuture
	newMap[1] = farFuture
	newMap[2] = 1
	history := &slashpb.AttestationHistory{
		TargetToSource:     newMap,
		LatestEpochWritten: 2,
	}

	newMap2 := make(map[uint64]uint64)
	// The validator attested at target epoch 1 and 3 but had no attestations for target epochs 0 and 2.
	newMap2[0] = farFuture
	newMap2[1] = 0
	newMap2[2] = farFuture
	newMap2[3] = 2
	history2 := &slashpb.AttestationHistory{
		TargetToSource:     newMap2,
		LatestEpochWritten: 3,
	}

	attestationHistory := make(map[[48]byte]*slashpb.AttestationHistory)
	attestationHistory[pubKeys[0]] = history
	attestationHistory[pubKeys[1]] = history2

	require.NoError(t, db.SaveAttestationHistoryForPubKeys(context.Background(), attestationHistory), "Saving attestation history failed")
	require.NoError(t, db.ImportOldAttestationFormat(ctx), "Import attestation history failed")

	attHis, err := db.AttestationHistoryNewForPubKeys(ctx, pubKeys)
	require.NoError(t, err)
	for pk, encHis := range attHis {
		his, ok := attestationHistory[pk]
		require.Equal(t, true, ok)
		lew, err := encHis.getLatestEpochWritten(ctx)
		require.NoError(t, err)
		require.Equal(t, his.LatestEpochWritten, lew)
		for target, source := range his.TargetToSource {
			hd, err := encHis.getTargetData(ctx, target)
			require.NoError(t, err)
			require.Equal(t, source, hd.Source)
			require.DeepEqual(t, bytesutil.PadTo([]byte{1}, 32), hd.SigningRoot)
		}
	}
}
