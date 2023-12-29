package engineio

import "github.com/sarvalabs/go-moi-identifiers"

// StateDriver represents an interface for accessing and manipulating state information of an account.
// It is bounded to a particular account and can only mutate within applicable portions
// of the state within the bounds of the logic's namespace
type StateDriver interface {
	Address() identifiers.Address
	LogicID() identifiers.LogicID

	GetStorageEntry([]byte) ([]byte, bool)
	SetStorageEntry([]byte, []byte) bool
}

// StateMatrix is matrix indicating the use of different
// types of state and the element pointers for them
type StateMatrix map[StateKind]ElementPtr

// Persistent indicates if the StateMatrix has an entry for PersistentState
func (matrix StateMatrix) Persistent() bool {
	_, exists := matrix[StatePersistent]

	return exists
}

// Ephemeral indicates if the StateMatrix has an entry for EphemeralState
func (matrix StateMatrix) Ephemeral() bool {
	_, exists := matrix[StateEphemeral]

	return exists
}

// StateKind represents the kind of stateful data
type StateKind int

const (
	StateEphemeral StateKind = iota
	StatePersistent
)
