package bls_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
)

func TestMarshalUnmarshal(t *testing.T) {
	b := bls.RandKey().Marshal()
	b32 := bytesutil.ToBytes32(b)
	pk, err := bls.SecretKeyFromBytes(b32[:])
	if err != nil {
		t.Fatal(err)
	}
	pk2, err := bls.SecretKeyFromBytes(b32[:])
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pk.Marshal(), pk2.Marshal()) {
		t.Errorf("Keys not equal, received %#x == %#x", pk.Marshal(), pk2.Marshal())
	}
}

func TestSignVerify(t *testing.T) {
	priv := bls.RandKey()
	pub := priv.PublicKey()
	msg := []byte("hello")
	sig := priv.Sign(msg, 0)
	if !sig.Verify(msg, pub, 0) {
		t.Error("Signature did not verify")
	}
}

func TestVerifyAggregate(t *testing.T) {
	pubkeys := make([]*bls.PublicKey, 0, 100)
	sigs := make([]*bls.Signature, 0, 100)
	var msgs [][32]byte
	for i := 0; i < 100; i++ {
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		priv := bls.RandKey()
		pub := priv.PublicKey()
		sig := priv.Sign(msg[:], 0)
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	aggSig := bls.AggregateSignatures(sigs)
	if !aggSig.VerifyAggregate(pubkeys, msgs, 0) {
		t.Error("Signature did not verify")
	}
}

func TestVerifyAggregateCommon(t *testing.T) {
	pubkeys := make([]*bls.PublicKey, 0, 100)
	sigs := make([]*bls.Signature, 0, 100)
	msg := [32]byte{'h', 'e', 'l', 'l', 'o'}
	for i := 0; i < 100; i++ {
		priv := bls.RandKey()
		pub := priv.PublicKey()
		sig := priv.Sign(msg[:], 0)
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
	}
	aggSig := bls.AggregateSignatures(sigs)
	if !aggSig.VerifyAggregateCommon(pubkeys, msg, 0) {
		t.Error("Signature did not verify")
	}
}

func TestVerifyAggregate_ReturnsFalseOnEmptyPubKeyList(t *testing.T) {
	var pubkeys []*bls.PublicKey
	sigs := make([]*bls.Signature, 0, 100)
	msg := [32]byte{'h', 'e', 'l', 'l', 'o'}

	aggSig := bls.AggregateSignatures(sigs)
	if aggSig.VerifyAggregateCommon(pubkeys, msg, 0 /*domain*/) != false {
		t.Error("Expected VerifyAggregate to return false with empty input " +
			"of public keys.")
	}
}

func TestComputeDomain_OK(t *testing.T) {
	tests := []struct {
		epoch      uint64
		domainType uint64
		domain     uint64
	}{
		{epoch: 1, domainType: 4, domain: 4},
		{epoch: 2, domainType: 4, domain: 4},
		{epoch: 2, domainType: 5, domain: 5},
		{epoch: 3, domainType: 4, domain: 4},
		{epoch: 3, domainType: 5, domain: 5},
	}
	for _, tt := range tests {
		if bls.ComputeDomain(bytesutil.Bytes4(tt.domainType)) != tt.domain {
			t.Errorf("wanted domain version: %d, got: %d", tt.domain, bls.ComputeDomain(bytesutil.Bytes4(tt.domainType)))
		}
	}
}

func TestSecretKeyFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{
		{
			name: "Nil",
			err:  errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Short",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Long",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Bad",
			input: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:   errors.New("could not unmarshal bytes into secret key: err blsSecretKeyDeserialize ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
		{
			name:  "Good",
			input: []byte{0x25, 0x29, 0x5f, 0x0d, 0x1d, 0x59, 0x2a, 0x90, 0xb3, 0x33, 0xe2, 0x6e, 0x85, 0x14, 0x97, 0x08, 0x20, 0x8e, 0x9f, 0x8e, 0x8b, 0xc1, 0x8f, 0x6c, 0x77, 0xbd, 0x62, 0xf8, 0xad, 0x7a, 0x68, 0x66},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := bls.SecretKeyFromBytes(test.input)
			if test.err != nil {
				if err == nil {
					t.Errorf("No error returned: expected %v", test.err)
				} else if test.err.Error() != err.Error() {
					t.Errorf("Unexpected error returned: expected %v, received %v", test.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error returned: %v", err)
				} else {
					if bytes.Compare(res.Marshal(), test.input) != 0 {
						t.Errorf("Unexpected result: expected %x, received %x", test.input, res.Marshal())
					}
				}
			}

		})
	}
}

func TestPublicKeyFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{
		{
			name: "Nil",
			err:  errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Short",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Long",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Bad",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("could not unmarshal bytes into public key: err blsPublicKeyDeserialize 000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			name:  "Good",
			input: []byte{0xa9, 0x9a, 0x76, 0xed, 0x77, 0x96, 0xf7, 0xbe, 0x22, 0xd5, 0xb7, 0xe8, 0x5d, 0xee, 0xb7, 0xc5, 0x67, 0x7e, 0x88, 0xe5, 0x11, 0xe0, 0xb3, 0x37, 0x61, 0x8f, 0x8c, 0x4e, 0xb6, 0x13, 0x49, 0xb4, 0xbf, 0x2d, 0x15, 0x3f, 0x64, 0x9f, 0x7b, 0x53, 0x35, 0x9f, 0xe8, 0xb9, 0x4a, 0x38, 0xe4, 0x4c},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := bls.PublicKeyFromBytes(test.input)
			if test.err != nil {
				if err == nil {
					t.Errorf("No error returned: expected %v", test.err)
				} else if test.err.Error() != err.Error() {
					t.Errorf("Unexpected error returned: expected %v, received %v", test.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error returned: %v", err)
				} else {
					if bytes.Compare(res.Marshal(), test.input) != 0 {
						t.Errorf("Unexpected result: expected %x, received %x", test.input, res.Marshal())
					}
				}
			}

		})
	}
}

func TestSignatureFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{
		{
			name: "Nil",
			err:  errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Short",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Long",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Bad",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("could not unmarshal bytes into signature: err blsSignatureDeserialize 000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			name:  "Good",
			input: []byte{0xab, 0xb0, 0x12, 0x4c, 0x75, 0x74, 0xf2, 0x81, 0xa2, 0x93, 0xf4, 0x18, 0x5c, 0xad, 0x3c, 0xb2, 0x26, 0x81, 0xd5, 0x20, 0x91, 0x7c, 0xe4, 0x66, 0x65, 0x24, 0x3e, 0xac, 0xb0, 0x51, 0x00, 0x0d, 0x8b, 0xac, 0xf7, 0x5e, 0x14, 0x51, 0x87, 0x0c, 0xa6, 0xb3, 0xb9, 0xe6, 0xc9, 0xd4, 0x1a, 0x7b, 0x02, 0xea, 0xd2, 0x68, 0x5a, 0x84, 0x18, 0x8a, 0x4f, 0xaf, 0xd3, 0x82, 0x5d, 0xaf, 0x6a, 0x98, 0x96, 0x25, 0xd7, 0x19, 0xcc, 0xd2, 0xd8, 0x3a, 0x40, 0x10, 0x1f, 0x4a, 0x45, 0x3f, 0xca, 0x62, 0x87, 0x8c, 0x89, 0x0e, 0xca, 0x62, 0x23, 0x63, 0xf9, 0xdd, 0xb8, 0xf3, 0x67, 0xa9, 0x1e, 0x84},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := bls.SignatureFromBytes(test.input)
			if test.err != nil {
				if err == nil {
					t.Errorf("No error returned: expected %v", test.err)
				} else if test.err.Error() != err.Error() {
					t.Errorf("Unexpected error returned: expected %v, received %v", test.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error returned: %v", err)
				} else {
					if bytes.Compare(res.Marshal(), test.input) != 0 {
						t.Errorf("Unexpected result: expected %x, received %x", test.input, res.Marshal())
					}
				}
			}

		})
	}
}

func TestPublicKey_Copy(t *testing.T) {
	pubkeyA := bls.RandKey().PublicKey()
	pubkeyBytes := pubkeyA.Marshal()

	pubkeyB, _ := pubkeyA.Copy()
	pubkeyB.Aggregate(bls.RandKey().PublicKey())

	if !bytes.Equal(pubkeyA.Marshal(), pubkeyBytes) {
		t.Fatal("Pubkey was mutated after copy")
	}
}