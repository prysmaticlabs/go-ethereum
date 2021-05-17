package sync

import (
	"errors"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
)

func writeContextToStream(stream network.Stream, chain blockchainService) error {
	rpcCtx, err := rpcContext(stream, chain)
	if err != nil {
		return err
	}
	// Exit early if there is an empty context.
	if len(rpcCtx) == 0 {
		return nil
	}
	_, err = stream.Write(rpcCtx)
	return err
}

func readContextFromStream(stream network.Stream, chain blockchainService) ([]byte, error) {
	rpcCtx, err := rpcContext(stream, chain)
	if err != nil {
		return nil, err
	}
	if len(rpcCtx) == 0 {
		return []byte{}, nil
	}
	// Read context (fork-digest) from stream
	b := make([]byte, 4)
	if _, err := stream.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func rpcContext(stream network.Stream, chain blockchainService) ([]byte, error) {
	_, _, version, err := p2p.TopicDeconstructor(string(stream.Protocol()))
	if err != nil {
		return nil, err
	}
	switch version {
	case p2p.SchemaVersionV1:
		// Return empty context for a v1 method.
		return []byte{}, nil
	default:
		return nil, errors.New("invalid version of %s registered for topic: %s")
	}
}
