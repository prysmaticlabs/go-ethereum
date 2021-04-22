package epoch_processing

import (
	"testing"

	"github.com/prysmaticlabs/prysm/spectest/shared/phase0/epoch_processing"
)

func TestMinimal_Phase0_EpochProcessing_ReseRegistryUpdates(t *testing.T) {
	epoch_processing.RunRegistryUpdatesTests(t, "minimal")
}
