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

// StateDriver represents an interface for accessing and manipulating state information of an account.
// It is bounded to a particular account and can only mutate within applicable portions
// of the state within the bounds of the logic's namespace
type StateDriver interface {
	Address() Address
	LogicID() LogicID

	GetStorageEntry([]byte) ([]byte, bool)
	SetStorageEntry([]byte, []byte) bool
}

// StateMatrix is matrix indicating the use of different
// types of state and the element pointers for them
type StateMatrix map[StateKind]ElementPtr

// Persistent indicates if the StateMatrix has an entry for PersistentState
func (matrix StateMatrix) Persistent() bool {
	_, exists := matrix[PersistentState]

	return exists
}

// Ephemeral indicates if the StateMatrix has an entry for EphemeralState
func (matrix StateMatrix) Ephemeral() bool {
	_, exists := matrix[EphemeralState]

	return exists
}

// StateKind represents the kind of stateful data
type StateKind int

const (
	PersistentState StateKind = iota
	EphemeralState
)

var stateKindToString = map[StateKind]string{
	PersistentState: "persistent",
	EphemeralState:  "ephemeral",
}

var stateKindFromString = map[string]StateKind{
	"persistent": PersistentState,
	"ephemeral":  EphemeralState,
}

// String implements the Stringer interface for StateKind
func (state StateKind) String() string {
	str, ok := stateKindToString[state]
	if !ok {
		panic("unknown StateKind variant")
	}

	return str
}

// Polorize implements the polo.Polorizable interface for StateKind
func (state StateKind) Polorize() (*polo.Polorizer, error) {
	polorizer := polo.NewPolorizer()
	polorizer.PolorizeString(state.String())

	return polorizer, nil
}

// Depolorize implements the polo.Depolorizable interface for StateKind
func (state *StateKind) Depolorize(depolorizer *polo.Depolorizer) error {
	raw, err := depolorizer.DepolorizeString()
	if err != nil {
		return err
	}

	kind, ok := stateKindFromString[raw]
	if !ok {
		return errors.New("invalid StateKind value")
	}

	*state = kind

	return nil
}

// MarshalJSON implements the json.Marshaller interface for StateKind
func (state StateKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

// UnmarshalJSON implements the json.Unmarshaller interface for StateKind
func (state *StateKind) UnmarshalJSON(data []byte) error {
	raw := new(string)
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}

	kind, ok := stateKindFromString[*raw]
	if !ok {
		return errors.New("invalid StateKind value")
	}

	*state = kind

	return nil
}

// MarshalYAML implements the yaml.Marshaller interface for StateKind
func (state StateKind) MarshalYAML() (interface{}, error) {
	return state.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaller interface for StateKind
func (state *StateKind) UnmarshalYAML(node *yaml.Node) error {
	raw := new(string)
	if err := node.Decode(raw); err != nil {
		return err
	}

	kind, ok := stateKindFromString[*raw]
	if !ok {
		return errors.New("invalid StateKind value")
	}

	*state = kind

	return nil
}
