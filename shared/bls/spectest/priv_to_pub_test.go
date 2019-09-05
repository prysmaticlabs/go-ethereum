package spectest

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"path"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/testutil"
)

func TestPrivToPubYaml(t *testing.T) {
	testFolders, testFolderPath := testutil.TestFolders(t, "general", "bls/priv_to_pub/small")

	for _, folder := range testFolders {
		t.Run(folder.Name(), func(t *testing.T) {
			file, err := loadBlsYaml(path.Join(testFolderPath, folder.Name(), "data.yaml"))
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			test := &PrivToPubTest{}
			if err := yaml.Unmarshal(file, test); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			pkBytes, err := hexutil.Decode(test.Input)
			if err != nil {
				t.Fatalf("Cannot decode string to bytes: %v", err)
			}
			sk, err := bls.SecretKeyFromBytes(pkBytes)
			if err != nil {
				t.Fatalf("Cannot unmarshal input to secret key: %v", err)
			}

			outputBytes, err := hexutil.Decode(test.Output)
			if err != nil {
				t.Fatalf("Cannot decode string to bytes: %v", err)
			}
			if !bytes.Equal(outputBytes, sk.PublicKey().Marshal()) {
				t.Fatal("Output does not marshaled public key bytes")
			}
		})
	}
}
