package sync

import (
	"bytes"
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/p2p/types"
	p2ptest "github.com/prysmaticlabs/prysm/beacon-chain/p2p/testing"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestRegularSync_generateErrorResponse(t *testing.T) {
	r := &Service{
		p2p: p2ptest.NewTestP2P(t),
	}
	data, err := r.generateErrorResponse(responseCodeServerError, "something bad happened")
	require.NoError(t, err)

	buf := bytes.NewBuffer(data)
	b := make([]byte, 1)
	_, err = buf.Read(b)
	require.NoError(t, err)
	assert.Equal(t, responseCodeServerError, b[0], "The first byte was not the status code")
	msg := &types.ErrorMessage{}
	require.NoError(t, r.p2p.Encoding().DecodeWithMaxLength(buf, msg))
	assert.Equal(t, "something bad happened", string(*msg), "Received the wrong message")
}
