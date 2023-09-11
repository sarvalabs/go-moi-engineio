package engineio

import "context"

// EngineKind is an enum with its variants
// representing the set of valid engines
type EngineKind string

const (
	// PISA is the EngineKind for the PISA VM Runtime.
	// The canonical implementation is available at https://github.com/sarvalabs/go-pisa
	PISA EngineKind = "PISA"

	// MERU is the EngineKind for a hypothetical engine runtime that works as a
	// WASI (WebAssembly) based VM Runtime for MOI. No implementation exists yet.
	MERU EngineKind = "MERU"
)

// EngineFuel is a measure of execution effort
type EngineFuel = uint64

// Engine is an execution engine interface with a specific EngineKind.
// A new Engine instance can be spawned from an EngineRuntime with its
// SpawnEngine method an is bound to a specific Logic and EnvDriver.
//
// An Engine can be used to perform calls on its Logic with an IxnDriver
// and some optional participants with their CtxDriver objects.
type Engine interface {
	// Kind returns the kind of engine
	Kind() EngineKind

	// Call calls a logical callsite on the Engine's Logic.
	// The callsite and input calldata are provided in the given IxnDriver.
	// Optionally accepts some participant CtxDriver objects based on the interaction type.
	Call(context.Context, IxnDriver, ...CtxDriver) (CallResult, error)
}

// EngineRuntime is an interface that defines an engine runtime.

// EngineRuntime is the base definition for execution engine runtime. It is
// used for runtime level behavioural capabilities rather for logic execution.
// This can include:
// - Compiling Manifest objects for the runtime into a Logic
// - Spawning execution Engine instances for a specific Logic
// - Validating input calldata for a specific callsite on a Logic
// - Obtaining a calldata encoder for a specific callsite on a Logic
// - Decoding DependencyDriver and ErrorResult objects for the runtime
type EngineRuntime interface {
	// Kind returns the kind of engine that the factory can produce
	Kind() EngineKind
	// Version returns the semver version string of the engine runtime
	Version() string

	// SpawnEngine returns a new Engine instance and initializes it with some
	// Fuel, a LogicDriver, the CtxDriver associated with the logic and an EnvDriver.
	// Will return an error if the LogicDriver and its CtxDriver do not match.
	SpawnEngine(EngineFuel, Logic, CtxDriver, EnvDriver) (Engine, error)

	// CompileManifest generates a LogicDescriptor from a Manifest, which can then be used to generate
	// a LogicDriver object. The fuel spent during compile is returned with any potential error.
	CompileManifest(EngineFuel, *Manifest) (*LogicDescriptor, EngineFuel, error)

	// ValidateCalldata verifies the calldata and callsite in an IxnObject.
	// The LogicDriver must describe a callsite which accepts the calldata.
	ValidateCalldata(Logic, IxnDriver) error

	// GetElementGenerator returns a generator function for an element schema with the
	// given ElementKind. Returns false, if no such element is defined by the runtime
	GetElementGenerator(ElementKind) (ManifestElementGenerator, bool)

	// GetCallEncoder returns a CallEncoder object for a given
	// callsite element pointer from a LogicDriver object
	GetCallEncoder(*Callsite, Logic) (CallEncoder, error)

	// DecodeDependencyDriver decodes the given bytes of the given
	// encoding into a DepDriver that is supported by the engine runtime
	DecodeDependencyDriver([]byte, Encoding) (DependencyDriver, error)

	// DecodeErrorResult decodes the given bytes into an
	// ErrorResult that is used by the engine runtime
	DecodeErrorResult([]byte) (ErrorResult, error)
}
