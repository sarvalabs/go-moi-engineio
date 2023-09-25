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
// SpawnEngine method and is bound to a specific LogicDriver and EnvironmentDriver.
//
// An Engine can be used to perform calls on its LogicDriver with an InteractionDriver
// and some optional participants with their StateDriver objects.
type Engine interface {
	// Kind returns the kind of engine
	Kind() EngineKind

	// Call calls a logical callsite on the Engine's LogicDriver.
	// The callsite and input calldata are provided in the given InteractionDriver.
	// Optionally accepts some participant StateDriver objects based on the interaction type.
	Call(context.Context, InteractionDriver, ...StateDriver) (CallResult, error)
}

// EngineRuntime is the base definition for execution engine runtime. It is
// used for runtime level behavioural capabilities rather for logic execution.
//
// This can include:
//   - Compiling Manifest objects for the runtime into a LogicDriver
//   - Spawning execution Engine instances for a specific LogicDriver
//   - Validating input calldata for a specific callsite on a LogicDriver
//   - Obtaining a calldata encoder for a specific callsite on a LogicDriver
//   - Decoding DependencyDriver and ErrorResult objects for the runtime
type EngineRuntime interface {
	// Kind returns the kind of engine that the factory can produce
	Kind() EngineKind
	// Version returns the semver version string of the engine runtime
	Version() string

	// SpawnEngine returns a new Engine instance and initializes it with some EngineFuel,
	// a LogicDriver, the StateDriver associated with the logic and an EnvironmentDriver.
	// Will return an error if the LogicDriver and its CtxDriver do not match.
	SpawnEngine(EngineFuel, LogicDriver, StateDriver, EnvironmentDriver) (Engine, error)

	// CompileManifest generates a LogicDescriptor from a Manifest, which can then be used to generate
	// a LogicDriver object. The fuel spent during compile is returned with any potential error.
	CompileManifest(EngineFuel, *Manifest) (*LogicDescriptor, EngineFuel, error)

	// ValidateCalldata verifies the calldata and callsite in an InteractionDriver.
	// The LogicDriver must describe a callsite which accepts the calldata.
	ValidateCalldata(LogicDriver, InteractionDriver) error

	// GetElementGenerator returns a generator function for an element schema with the
	// given ElementKind. Returns false, if no such element is defined by the runtime
	GetElementGenerator(ElementKind) (ManifestElementGenerator, bool)

	// GetCallEncoder returns a CallEncoder object for a given
	// callsite element pointer from a LogicDriver object
	GetCallEncoder(*Callsite, LogicDriver) (CallEncoder, error)

	// DecodeDependencyDriver decodes the given bytes of the given
	// encoding into a DepDriver that is supported by the engine runtime
	DecodeDependencyDriver([]byte, Encoding) (DependencyDriver, error)

	// DecodeErrorResult decodes the given bytes into an
	// ErrorResult that is used by the engine runtime
	DecodeErrorResult([]byte) (ErrorResult, error)
}

// registry is an in-memory registry of supported EngineRuntime instances.
// Support for different engine runtimes is only available if they are first registered with this package.
var registry = map[EngineKind]EngineRuntime{}

// RegisterRuntime registers an EngineRuntime with the package.
// If a runtime instance already exists for the EngineKind, it is overwritten.
func RegisterRuntime(runtime EngineRuntime) {
	registry[runtime.Kind()] = runtime
}

// FetchRuntime retrieves an EngineRuntime for a given EngineKind.
// If the runtime for the engine kind is not registered, returns false.
func FetchRuntime(kind EngineKind) (EngineRuntime, bool) {
	runtime, exists := registry[kind]
	if !exists {
		return nil, false
	}

	return runtime, true
}
