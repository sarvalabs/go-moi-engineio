package engineio

import (
	"reflect"
	"sort"

	"github.com/pkg/errors"
	"github.com/sarvalabs/go-polo"
)

// CallEncoder is an interface with capabilities to encode inputs and decode outputs for a specific callable site.
// It can be generated from either a Manifest or LogicDriver.
type CallEncoder interface {
	EncodeInputs(map[string]any, ReferenceProvider) ([]byte, error)
	DecodeOutputs([]byte) (map[string]any, error)
}

// ReferenceValue is a reference identifier
// that resolves to an encodable value
type ReferenceValue string

func (ref ReferenceValue) String() string {
	return "ref<" + string(ref) + ">"
}

// ReferenceProvider resolves a ReferenceVal into
// an encodable value confirming the resolution
type ReferenceProvider interface {
	GetReference(ReferenceValue) (any, bool)
}

// EncodeValues encodes a value into a bytes, recursively resolving any internal type data.
// Expects a ReferenceProvider for resolving reference variables (can be nil, if no references are used)
func EncodeValues(value any, references ReferenceProvider) ([]byte, error) {
	switch val := value.(type) {
	// Object Type (ClassType)
	case map[string]any:
		document := make(polo.Document)

		// For each field in the object
		for field, v := range val {
			// Encode field value
			data, err := EncodeValues(v, references)
			if err != nil {
				return nil, err
			}

			document.SetRaw(field, data)
		}

		return document.Bytes(), nil

	// Map Type (MapType)
	case map[any]any:
		// Create a new Polorizer
		polorizer := polo.NewPolorizer()

		// Reflect the value object and sort its keys
		reflected := reflect.ValueOf(val)
		keys := reflected.MapKeys()
		sort.Slice(keys, polo.MapSorter(keys))

		// Iterate over the sorted keys and encode
		// each key-value pair to the polorizer
		for _, key := range keys {
			// Encode key value
			kdata, err := EncodeValues(key.Interface(), references)
			if err != nil {
				return nil, err
			}

			// Encode val value
			vdata, err := EncodeValues(reflected.MapIndex(key).Interface(), references)
			if err != nil {
				return nil, err
			}

			// Write both key and val data into polorizer
			_ = polorizer.PolorizeAny(kdata)
			_ = polorizer.PolorizeAny(vdata)
		}

		return polorizer.Bytes(), nil

	// List Type (ArrayType & VarrayType)
	case []any:
		// Create a new Polorizer
		polorizer := polo.NewPolorizer()

		// For each element in the list
		for _, elem := range val {
			// Encode element value
			data, err := EncodeValues(elem, references)
			if err != nil {
				return nil, err
			}

			// Write element data into polorizer
			_ = polorizer.PolorizeAny(data)
		}

		return polorizer.Bytes(), nil

	// Reference Type
	case ReferenceValue:
		// If no reference provider is given, error
		if references == nil {
			return nil, errors.New("encountered reference value without a ref provider")
		}

		// Resolve the reference
		deref, ok := references.GetReference(val)
		if !ok {
			return nil, errors.Errorf("unable to resolve reference '%v'", val)
		}

		// Encode the dereferenced value
		return EncodeValues(deref, references)

	// Simple Type
	default:
		data, err := polo.Polorize(val)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
}
