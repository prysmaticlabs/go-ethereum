package spectest

import (
	"testing"
)

func TestBlockProcessingMainnetYaml(t *testing.T) {
	t.Skip("Disabled until v0.9.0 (#3865) completes")
	runBlockProcessingTest(t, "mainnet")
}
