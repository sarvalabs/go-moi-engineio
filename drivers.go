package engineio

import (
	"encoding/json"
	"fmt"

	"github.com/sarvalabs/go-polo"
)

// EnvDriver represents a driver for environmental information.
// It describes information about the execution context such
// as the consensus cluster ID or execution timestamp.
type EnvDriver interface {
	Timestamp() int64
	ClusterID() string
}

// DependencyDriver represents an interface
// for an engine's element dependency manager.
//
// It must be expressible as a string and encodable with JSON and POLO.
// It manages the dependency relationship between element pointers with
// the pointers being vertices and their relationship being directional edges.
type DependencyDriver interface {
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler

	polo.Polorizable
	polo.Depolorizable

	Insert(uint64, ...uint64)
	Remove(uint64)

	Size() uint64
	Iter() <-chan uint64
	Contains(uint64) bool
	Edges(uint64) []uint64
	Dependencies(uint64) []uint64
}
