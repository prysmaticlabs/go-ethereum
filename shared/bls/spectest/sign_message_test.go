package autogenerated

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/prysmaticlabs/prysm/shared/bls"
)

func TestSignMessageYaml(t *testing.T) {
	file, err := ioutil.ReadFile("sign_msg_formatted.yaml")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	test := &SignMessageTest{}
	if err := yaml.Unmarshal(file, test); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	for i, tt := range test.TestCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			sk, err := bls.SecretKeyFromBytes(tt.Input.Privkey)
			if err != nil {
				t.Fatalf("Cannot unmarshal input to secret key: %v", err)
			}
			domain, _ := binary.Uvarint(tt.Input.Domain)
			if err != nil {
				t.Fatal(err)
			}
			sig := sk.Sign(tt.Input.Message, domain)
			if !bytes.Equal(tt.Output, sig.Marshal()) {
				t.Errorf("Signature does not match the expected output. Expected %#x but received %#x", tt.Output, sig.Marshal())
			}
		})
	}
}
