package engineio

type (
	// Hash is a 256-bit checksum digest
	Hash = [32]byte
	// Address is a 256-bit unique identifier for a participant.
	Address = [32]byte
)

type (
	// ElementKind is a type alias for an element kind string
	ElementKind string
	// ElementPtr is a type alias for an element pointer
	ElementPtr = uint64
)

// Classdef represents a class definition in a Logic.
// It can be resolved from a string by looking it up on the Logic
type Classdef interface {
	Name() string
	Pointer() ElementPtr
	Encode(Encoding) ([]byte, error)
}

// Callsite represents a callable point in a Logic.
// It can be resolved from a string identifier with the GetCallsite method on Logic
type Callsite interface {
	Kind() CallsiteKind
	Name() string
	Pointer() ElementPtr
	Encode(Encoding) ([]byte, error)
}

// CallsiteKind represents the type of callable point in a Logic.
type CallsiteKind int

const (
	DeployerCallsite CallsiteKind = iota
	EnlisterCallsite
	InvokableCallsite
	InteractableCallsite
)
