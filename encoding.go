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
	// Kind returns the kind of engine that the encoder works for
	Kind() EngineKind

	// EncodeValues encodes any given value into a byte representation suitable for the runtime.
	// It also supports reference based encoding by accepting a ReferenceProvider
	EncodeValues(any, ReferenceProvider) ([]byte, error)

	// EncodeManifest encodes a Manifest with the given Encoding
	EncodeManifest(Manifest, Encoding) ([]byte, error)
	// DecodeManifest decodes a Manifest from some raw data of the given Encoding
	DecodeManifest([]byte, Encoding) (Manifest, error)

	// EncodeErrorResult encodes an ErrorResult into some
	// byte representation suitable for the runtime
	EncodeErrorResult(ErrorResult) ([]byte, error)
	// DecodeErrorResult decodes the given bytes into an
	// ErrorResult that is suitable for the engine runtime
	DecodeErrorResult([]byte) (ErrorResult, error)

	// EncodeDependencies encodes a DependencyDriver into the given encoding
	EncodeDependencies(DependencyDriver, Encoding) ([]byte, error)
	// DecodeDependencies decodes the given bytes of the given encoding
	// into a DependencyDriver that is supported by the engine runtime
	DecodeDependencies([]byte, Encoding) (DependencyDriver, error)

	// EncodeCallsite encodes a Callsite into the given encoding
	EncodeCallsite(Callsite, Encoding) ([]byte, error)
	// DecodeCallsite decodes a Callsite from the given encoding
	DecodeCallsite([]byte, Encoding) (Callsite, error)

	// EncodeClassdef encodes a Classdef into the given encoding
	EncodeClassdef(Classdef, Encoding) ([]byte, error)
	// DecodeClassdef decodes a Classdef from the given encoding
	DecodeClassdef([]byte, Encoding) (Classdef, error)
}

// Reference is a reference identifier
// that resolves to an encodable value
type Reference interface {
	Identifier() string
}

// ReferenceProvider resolves a Reference into
// an encodable value confirming the resolution
type ReferenceProvider interface {
	GetReference(Reference) (any, bool)
}
