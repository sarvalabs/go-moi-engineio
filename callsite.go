package engineio

import (
	"encoding/json"
	"errors"

	"github.com/sarvalabs/go-polo"
	"gopkg.in/yaml.v3"
)

// Callsite represents a callable point in a LogicDriver.
// It can be resolved from a string identifier with the GetCallsite method on LogicDriver
type Callsite struct {
	Ptr  ElementPtr
	Kind CallsiteKind
}

// CallsiteKind represents the type of callable point in a LogicDriver.
type CallsiteKind int

const (
	LocalCallsite CallsiteKind = iota
	InvokableCallsite
	InteractableCallsite
	DeployerCallsite
	EnlisterCallsite
)

var callsiteKindToString = map[CallsiteKind]string{
	LocalCallsite:        "local",
	InvokableCallsite:    "invokable",
	InteractableCallsite: "interactable",
	DeployerCallsite:     "deployer",
	EnlisterCallsite:     "enlister",
}

var callsiteKindFromString = map[string]CallsiteKind{
	"local":        LocalCallsite,
	"invokable":    InvokableCallsite,
	"interactable": InteractableCallsite,
	"deployer":     DeployerCallsite,
	"enlister":     EnlisterCallsite,
}

// String implements the Stringer interface for CallsiteKind
func (callsite CallsiteKind) String() string {
	str, ok := callsiteKindToString[callsite]
	if !ok {
		panic("unknown CallsiteKind variant")
	}

	return str
}

// Polorize implements the polo.Polorizable interface for CallsiteKind
func (callsite CallsiteKind) Polorize() (*polo.Polorizer, error) {
	polorizer := polo.NewPolorizer()
	polorizer.PolorizeString(callsite.String())

	return polorizer, nil
}

// Depolorize implements the polo.Depolorizable interface for CallsiteKind
func (callsite *CallsiteKind) Depolorize(depolorizer *polo.Depolorizer) error {
	raw, err := depolorizer.DepolorizeString()
	if err != nil {
		return err
	}

	kind, ok := callsiteKindFromString[raw]
	if !ok {
		return errors.New("invalid CallsiteKind value")
	}

	*callsite = kind

	return nil
}

// MarshalJSON implements the json.Marshaller interface for CallsiteKind
func (callsite CallsiteKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(callsite.String())
}

// UnmarshalJSON implements the json.Unmarshaller interface for CallsiteKind
func (callsite *CallsiteKind) UnmarshalJSON(data []byte) error {
	raw := new(string)
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}

	kind, ok := callsiteKindFromString[*raw]
	if !ok {
		return errors.New("invalid CallsiteKind value")
	}

	*callsite = kind

	return nil
}

// MarshalYAML implements the yaml.Marshaller interface for CallsiteKind
func (callsite CallsiteKind) MarshalYAML() (interface{}, error) {
	return callsite.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaller interface for CallsiteKind
func (callsite *CallsiteKind) UnmarshalYAML(node *yaml.Node) error {
	raw := new(string)
	if err := node.Decode(raw); err != nil {
		return err
	}

	kind, ok := callsiteKindFromString[*raw]
	if !ok {
		return errors.New("invalid CallsiteKind value")
	}

	*callsite = kind

	return nil
}
