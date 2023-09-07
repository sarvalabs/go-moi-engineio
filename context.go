package engineio

import (
	"encoding/json"
	"errors"

	"github.com/sarvalabs/go-polo"
	"gopkg.in/yaml.v3"
)

type (
	// Hash is a 256-bit checksum digest
	Hash = [32]byte
	// Address is a 256-bit unique identifier for a participant.
	Address = [32]byte
)

// CtxDriver represents an interface for accessing and manipulating
// context information of an account state. It is bounded to the context
// of particular account and can only mutate within applicable portions
// of the context state within the bounds of the logic's namespace
type CtxDriver interface {
	Address() Address
	LogicID() LogicID

	GetStorageEntry([]byte) ([]byte, bool)
	SetStorageEntry([]byte, []byte) bool
}

// ContextStateMatrix is matrix indicating the use of different
// types of context state and the element pointers for them
type ContextStateMatrix map[ContextStateKind]ElementPtr

// Persistent indicates if the ContextStateMatrix has an entry for PersistentState
func (matrix ContextStateMatrix) Persistent() bool {
	_, exists := matrix[PersistentState]

	return exists
}

// Ephemeral indicates if the ContextStateMatrix has an entry for EphemeralState
func (matrix ContextStateMatrix) Ephemeral() bool {
	_, exists := matrix[EphemeralState]

	return exists
}

// ContextStateKind represents the scope of stateful data in a context
type ContextStateKind int

const (
	PersistentState ContextStateKind = iota
	EphemeralState
)

var contextStateKindToString = map[ContextStateKind]string{
	PersistentState: "persistent",
	EphemeralState:  "ephemeral",
}

var contextStateKindFromString = map[string]ContextStateKind{
	"persistent": PersistentState,
	"ephemeral":  EphemeralState,
}

// String implements the Stringer interface for ContextStateKind
func (state ContextStateKind) String() string {
	str, ok := contextStateKindToString[state]
	if !ok {
		panic("unknown ContextStateKind variant")
	}

	return str
}

// Polorize implements the polo.Polorizable interface for ContextStateKind
func (state ContextStateKind) Polorize() (*polo.Polorizer, error) {
	polorizer := polo.NewPolorizer()
	polorizer.PolorizeString(state.String())

	return polorizer, nil
}

// Depolorize implements the polo.Depolorizable interface for ContextStateKind
func (state *ContextStateKind) Depolorize(depolorizer *polo.Depolorizer) error {
	raw, err := depolorizer.DepolorizeString()
	if err != nil {
		return err
	}

	kind, ok := contextStateKindFromString[raw]
	if !ok {
		return errors.New("invalid ContextStateKind value")
	}

	*state = kind

	return nil
}

// MarshalJSON implements the json.Marshaller interface for ContextStateKind
func (state ContextStateKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

// UnmarshalJSON implements the json.Unmarshaller interface for ContextStateKind
func (state *ContextStateKind) UnmarshalJSON(data []byte) error {
	raw := new(string)
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}

	kind, ok := contextStateKindFromString[*raw]
	if !ok {
		return errors.New("invalid ContextStateKind value")
	}

	*state = kind

	return nil
}

// MarshalYAML implements the yaml.Marshaller interface for ContextStateKind
func (state ContextStateKind) MarshalYAML() (interface{}, error) {
	return state.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaller interface for ContextStateKind
func (state *ContextStateKind) UnmarshalYAML(node *yaml.Node) error {
	raw := new(string)
	if err := node.Decode(raw); err != nil {
		return err
	}

	kind, ok := contextStateKindFromString[*raw]
	if !ok {
		return errors.New("invalid ContextStateKind value")
	}

	*state = kind

	return nil
}
