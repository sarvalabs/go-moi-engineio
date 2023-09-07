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

// ReferenceVal is a reference identifier
// that resolves to an encodable value
type ReferenceVal string

func (ref ReferenceVal) String() string {
	return "ref<" + string(ref) + ">"
}

// ReferenceProvider resolves a ReferenceVal into
// an encodable value confirming the resolution
type ReferenceProvider interface {
	GetReference(ReferenceVal) (any, bool)
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
		sort.Slice(keys, MapSorter(keys))

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
	case ReferenceVal:
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

// MapSorter is used by the sort package to sort a slice of reflect.Value objects.
// Assumes that the reflect.Value objects can only be types which are comparable
// i.e, can be used as a map key. (will panic otherwise)
//
// todo: remove this when go-polo has a MapPolorizer implementation
func MapSorter(keys []reflect.Value) func(int, int) bool {
	return func(i int, j int) bool {
		a, b := keys[i], keys[j]
		if a.Kind() == reflect.Interface {
			a, b = a.Elem(), b.Elem()
		}

		switch a.Kind() {
		case reflect.Bool:
			return b.Bool()

		case reflect.String:
			return a.String() < b.String()

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() < b.Int()

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return a.Uint() < b.Uint()

		case reflect.Float32, reflect.Float64:
			return a.Float() < b.Float()

		case reflect.Array:
			if a.Len() != b.Len() {
				panic("array length must equal")
			}

			for i := 0; i < a.Len(); i++ {
				result := MapCompare(a.Index(i), b.Index(i))
				if result == 0 {
					continue
				}

				return result < 0
			}

			return false
		}

		panic("unsupported key compare")
	}
}

// MapCompare returns an integer representing the comparison between two reflect.Value objects.
// Assumes that a and b can only have a type that is comparable. (will panic otherwise).
// Returns 1 (a > b); 0 (a == b); -1 (a < b)
//
// todo: remove this when go-polo has a MapPolorizer implementation
func MapCompare(a, b reflect.Value) int {
	if a.Kind() == reflect.Interface {
		a, b = a.Elem(), b.Elem()
	}

	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		av, bv := a.Int(), b.Int()

		switch {
		case av < bv:
			return -1
		case av == bv:
			return 0
		case av > bv:
			return 1
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		av, bv := a.Uint(), b.Uint()

		switch {
		case av < bv:
			return -1
		case av == bv:
			return 0
		case av > bv:
			return 1
		}

	case reflect.Float32, reflect.Float64:
		av, bv := a.Float(), b.Float()

		switch {
		case av < bv:
			return -1
		case av == bv:
			return 0
		case av > bv:
			return 1
		}

	case reflect.String:
		av, bv := a.String(), b.String()

		switch {
		case av < bv:
			return -1
		case av == bv:
			return 0
		case av > bv:
			return 1
		}

	case reflect.Array:
		if a.Len() != b.Len() {
			panic("array length must equal")
		}

		for i := 0; i < a.Len(); i++ {
			result := MapCompare(a.Index(i), b.Index(i))
			if result == 0 {
				continue
			}

			return result
		}

		return 0
	}

	panic("unsupported key compare")
}
