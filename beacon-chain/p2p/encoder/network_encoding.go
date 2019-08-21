package encoder

import (
	"io"

	"github.com/gogo/protobuf/proto"
)

const (
	SSZ = iota
	SSZ_SNAPPY
)

// NetworkEncoding represents an encoder compatible with Ethereum 2.0 p2p.
type NetworkEncoding interface {
	// Decode reads bytes from the reader and decodes it to the provided message.
	Decode(io.Reader, proto.Message) error
	// Encode an arbitrary message to the provided writer.
	Encode(io.Writer, proto.Message) (int, error)
	// ProtocolSuffix returns the last part of the protocol ID to indicate the encoding scheme.
	ProtocolSuffix() string
}

type Encoding int

func (e Encoding) String() string {
	formats := []string{"SSZ", "SSZ_SNAPPY"}
	if int(e) >= len(formats) {
		// Send ssz as default if encoding format doesn't exist.
		return formats[0]
	}
	return formats[e]
}
