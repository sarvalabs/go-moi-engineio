package engineio

// Encoding is an enum with variants that describe
// encoding schemes supported for file objects.
type Encoding int

const (
	POLO Encoding = iota
	JSON
	YAML
)

// CallEncoder is an interface with capabilities to encode inputs and decode outputs for a specific callable site.
// It can be generated from either a Manifest or Logic.
type CallEncoder interface {
	EncodeInputs(map[string]any, ReferenceProvider) ([]byte, error)
	DecodeOutputs([]byte) (map[string]any, error)
}

// RuntimeEncoder is an interface that describes all encoding capabilities for a runtime
type RuntimeEncoder interface {
	// Kind returns the kind of engine that the factory can produce
	Kind() EngineKind

	EncodeValues(any, ReferenceProvider) ([]byte, error)

	EncodeManifest(Manifest, Encoding) ([]byte, error)
	DecodeManifest([]byte, Encoding) (Manifest, error)

	// DecodeErrorResult decodes the given bytes into an
	// ErrorResult that is used by the engine runtime
	DecodeErrorResult([]byte) (ErrorResult, error)
	EncodeErrorResult(ErrorResult) ([]byte, error)

	// DecodeDependencies decodes the given bytes of the given encoding
	// into a DependencyDriver that is supported by the engine runtime
	DecodeDependencies([]byte, Encoding) (DependencyDriver, error)
	EncodeDependencies(DependencyDriver, Encoding) ([]byte, error)

	EncodeCallsite(Callsite, Encoding) ([]byte, error)
	DecodeCallsite([]byte, Encoding) (Callsite, error)

	EncodeClassdef(Classdef, Encoding) ([]byte, error)
	DecodeClassdef([]byte, Encoding) (Classdef, error)
}

// Reference is a reference identifier
// that resolves to an encodable value
type Reference interface {
	String()
}

// ReferenceProvider resolves a Reference into
// an encodable value confirming the resolution
type ReferenceProvider interface {
	GetReference(Reference) (any, bool)
}
