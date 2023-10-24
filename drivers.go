package engineio

import (
	"fmt"
)

// EnvironmentDriver represents a driver for environmental information.
// It describes information about the execution context such
// as the consensus cluster ID or execution timestamp.
type EnvironmentDriver interface {
	Timestamp() int64
	ClusterID() string
}

// DependencyDriver represents an interface for an engine's element dependency manager.
// It manages the dependency relationship between element pointers with
// the pointers being vertices and their relationship being directional edges.
type DependencyDriver interface {
	fmt.Stringer

	Insert(uint64, ...uint64)
	Remove(uint64)

	Size() uint64
	Iter() <-chan uint64
	Contains(uint64) bool
	Edges(uint64) []uint64
	Dependencies(uint64) []uint64
}

// CryptographyDriver represents an interface for cryptographic operations.
// It can be used to validate signature formats and verify them for a public key.
// This interfaces allows us to pass the capabilities of go-moi's crypto package to different engine runtimes.
type CryptographyDriver interface {
	ValidateSignature(sig []byte) bool
	VerifySignature(data, sig, pub []byte) (bool, error)
}
