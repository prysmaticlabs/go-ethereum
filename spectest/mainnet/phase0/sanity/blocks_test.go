package sanity

import (
	"testing"

	"github.com/prysmaticlabs/prysm/spectest/shared/phase0/sanity"
)

func TestMainnet_Phase0_Sanity_Blocks(t *testing.T) {
	sanity.RunBlockProcessingTest(t, "mainnet")
}
