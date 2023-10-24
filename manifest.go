package engineio

// Manifest is the canonical deployment artifact for logics in MOI.
// It is a composite artifact that describes the bytecode, the binary interface (ABI) and
// other parameters for the runtime of choice. The spec for LogicManifests is available at
// https://sarvalabs.notion.site/Logic-Manifest-Standard-93f5fee1af8d4c3cad155b9827b97930?pvs=4
type Manifest interface {
	// Kind returns the kind of engine of the Manifest
	Kind() EngineKind
	// Hash returns the 256-bit hash of the Manifest.
	Hash() [32]byte
	// Size returns the number of elements in the Manifest
	Size() uint64

	// Syntax returns the syntax version of the Manifest
	Syntax() uint64
	// Engine returns the engine information of the Manifest as a ManifestEngine
	Engine() ManifestEngine
	// Header returns the header information of the Manifest as a ManifestHeader
	Header() ManifestHeader

	// Elements returns all the elements in the Manifest as an array of ManifestElement
	Elements() []ManifestElement
	// GetElement returns the ManifestElement from the Manifest with the given ElementPtr.
	// The boolean indicated is such an element exists in the Manifest.
	GetElement(ElementPtr) (ManifestElement, bool)

	// Encode returns the encoded bytes form of the Manifest for the specified encoding.
	Encode(Encoding) ([]byte, error)
	// GenerateCallEncoder returns a CallEncoder for a given Callsite on the Manifest.
	// An error occurs if no callable element with an ElementPtr corresponding to the given Callsite exists.
	GenerateCallEncoder(Callsite) (CallEncoder, error)
}

// ManifestEngine describes the engine specification in the Manifest
type ManifestEngine struct {
	Kind  EngineKind
	Flags []string
}

// ManifestElement describes a single element in the Manifest.
// It is converted into a LogicElement after compilation.
//
// Each element is of a particular type (described by the engine runtime) and is identified by
// a unique 64-bit pointer  and describes its dependencies with other elements in the manifest.
// The data of the manifest element must be resolved into the format specific for the runtime based on its
// kind. The raw object to decode into can be accessed with the GetElementGenerator method of EngineRuntime.
type ManifestElement struct {
	Ptr  ElementPtr
	Deps []ElementPtr
	Kind ElementKind
	Data any
}

// ManifestHeader represents the header for a Manifest and describes its syntax form
// and engine specification. Useful for determining which engine to use to handle the
// Manifest. Every engine's manifest implementation must be able to decode into this header.
type ManifestHeader struct {
	Syntax uint64         `yaml:"syntax" json:"syntax"`
	Engine ManifestEngine `yaml:"engine" json:"engine"`
}
