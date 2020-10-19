package kv

import (
	"context"
	"fmt"
	"testing"

	ethereum_slashing "github.com/prysmaticlabs/prysm/proto/slashing"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestSaveHighestAttestation(t *testing.T) {
	db := setupDB(t)
	ctx := context.Background()

	tests := []struct {
		name         string
		toSave       []*ethereum_slashing.HighestAttestation
		cacheEnabled bool
	}{
		{
			name: "save to cache",
			toSave: []*ethereum_slashing.HighestAttestation{
				&ethereum_slashing.HighestAttestation{
					HighestTargetEpoch: 1,
					HighestSourceEpoch: 0,
					ValidatorId:        1,
				},
			},
			cacheEnabled: true,
		},
		{
			name: "save to db",
			toSave: []*ethereum_slashing.HighestAttestation{
				&ethereum_slashing.HighestAttestation{
					HighestTargetEpoch: 1,
					HighestSourceEpoch: 0,
					ValidatorId:        2,
				},
			},
			cacheEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, att := range tt.toSave {
				db.highestAttCacheEnabled = tt.cacheEnabled

				require.NoError(t, db.SaveHighestAttestation(ctx, att), "Save highest attestation failed")

				found, err := db.HighestAttestation(ctx, att.ValidatorId)
				require.NoError(t, err)
				require.NotNil(t, found)
				require.Equal(t, att.ValidatorId, found.ValidatorId)
				require.Equal(t, att.HighestSourceEpoch, found.HighestSourceEpoch)
				require.Equal(t, att.HighestTargetEpoch, found.HighestTargetEpoch)
			}
		})
	}
}

func TestFetchNonExistingHighestAttestation(t *testing.T) {
	db := setupDB(t)
	ctx := context.Background()

	t.Run("cached", func(t *testing.T) {
		db.highestAttCacheEnabled = true
		found, err := db.HighestAttestation(ctx, 1)
		require.NoError(t, err)
		if found != nil {
			require.NoError(t, fmt.Errorf("should not find HighestAttestation"))
		}
	})

	t.Run("disk", func(t *testing.T) {
		db.highestAttCacheEnabled = false
		found, err := db.HighestAttestation(ctx, 1)
		require.NoError(t, err)
		if found != nil {
			require.NoError(t, fmt.Errorf("should not find HighestAttestation"))
		}
	})

}
