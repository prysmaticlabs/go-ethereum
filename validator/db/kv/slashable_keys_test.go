package kv

import (
	"context"
	"fmt"
	"testing"

	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestStore_SlashablePublicKeys(t *testing.T) {
	ctx := context.Background()
	numValidators := 100
	publicKeys := make([][48]byte, numValidators)
	for i := 0; i < numValidators; i++ {
		key := [48]byte{}
		copy(key[:], fmt.Sprintf("%d", i))
		publicKeys[i] = key
	}

	// No slashable keys returns empty.
	validatorDB := setupDB(t, publicKeys)
	received, err := validatorDB.SlashablePublicKeys(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, len(received))

	// Save half of the public keys as as slashable and attempt to retrieve.
	err = validatorDB.SaveSlashablePublicKeys(ctx, publicKeys[:50])
	require.NoError(t, err)
	received, err = validatorDB.SlashablePublicKeys(ctx)
	require.NoError(t, err)

	// Keys are not guaranteed to be ordered, so we create a map for comparisons.
	want := make(map[[48]byte]bool)
	for _, pubKey := range publicKeys[:50] {
		want[pubKey] = true
	}
	for _, pubKey := range received {
		ok := want[pubKey]
		require.Equal(t, true, ok)
	}
}
