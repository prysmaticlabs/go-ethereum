package keystore_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/gogo/protobuf/proto"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/keystore"
	"github.com/prysmaticlabs/prysm/shared/params"
)

func TestDepositInput_GeneratesPb(t *testing.T) {
	k1, err := keystore.NewKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	k2, err := keystore.NewKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	result := keystore.DepositInput(k1, k2)
	if !bytes.Equal(result.Pubkey, k1.PublicKey.Marshal()) {
		t.Errorf("Mismatched pubkeys in deposit input. Want = %x, got = %x", result.Pubkey, k1.PublicKey.Marshal())
	}

	sig, err := bls.SignatureFromBytes(result.ProofOfPossession)
	if err != nil {
		t.Fatal(err)
	}

	proofOfPossessionInputPb := proto.Clone(result).(*pb.DepositInput)
	proofOfPossessionInputPb.ProofOfPossession = nil
	proofOfPossessionInput, err := proto.Marshal(proofOfPossessionInputPb)
	if err != nil {
		t.Fatal(err)
	}

	if !sig.Verify(proofOfPossessionInput, k1.PublicKey, params.BeaconConfig().DomainDeposit) {
		t.Error("Invalid proof of proofOfPossession signature")
	}
}
