package p2p

import (
	"testing"
)

func TestBuildOptions(t *testing.T) {
	opts := buildOptions(&ServerConfig{})

	_ = opts
}
