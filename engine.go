package engineio

import "context"

// EngineKind is an enum with its variants
// representing the set of valid engines
type EngineKind int

const (
	// InvalidEngine is a base reserved engine variant
	InvalidEngine EngineKind = iota

	// PISA is the EngineKind for the PISA VM Runtime.
	// The canonical implementation is available at https://github.com/sarvalabs/go-pisa
	PISA

	// MERU is the EngineKind for a hypothetical engine runtime that works as a WASM
	// environment that allows custom runtime implementations to run within it
	MERU
)

// EngineFuel is a measure of execution effort
type EngineFuel = uint64

// EngineInstance is an execution engine runner with a specific EngineKind.
// A new EngineRunner instance can be spawned from an EngineRuntime with its
// Spawn method and is bound to a specific Logic and EnvironmentDriver.
//
// An Engine can be used to perform calls on its Logic with an InteractionDriver
// and some optional participants with their StateDriver objects.
type EngineInstance interface {
	// Kind returns the kind of engine
	Kind() EngineKind

	// Call calls a logical callsite on the Engine's Logic.
	// The callsite and input calldata are provided in the given InteractionDriver.
	// Optionally accepts some participant StateDriver objects based on the interaction type.
	Call(context.Context, InteractionDriver, StateDriver, ...StateDriver) (CallResult, error)
}

// EngineRuntime is the base definition for execution engine runtime. It is
// used for runtime level behavioural capabilities rather for logic execution.
//
// This can include:
//   - Compiling Manifest objects for the runtime into a LogicDriver
//   - Spawning execution EngineInstance for a specific LogicDriver
//   - Validating input calldata for a specific callsite on a LogicDriver
//   - Obtaining a calldata encoder for a specific callsite on a LogicDriver
type EngineRuntime interface {
	// RuntimeEncoder is an embedded interface for EngineRuntime, allowing
	// the runtime to specify encoding/decoding rules for runtime specific
	// types like Manifest, ErrorResult, DependencyDriver, etc.
	RuntimeEncoder

	// Kind returns the kind of engine that the factory can produce
	Kind() EngineKind
	// Version returns the semver version string of the engine runtime
	Version() string

	// GenerateRuntimeEncoder returns the RuntimeEncoder instance for the runtime
	GenerateRuntimeEncoder() RuntimeEncoder
	// GenerateCallEncoder returns a CallEncoder object for a
	// given callsite element pointer from a LogicDriver object
	GenerateCallEncoder(LogicDriver, Callsite) (CallEncoder, error)

	// SpawnInstance returns a new EngineInstance instance and initializes it with some
	// EngineFuel, a LogicDriver, the StateDriver associated with the logic and an EnvironmentDriver.
	// Will return an error if the Logic and its StateDriver do not match.
	SpawnInstance(LogicDriver, EngineFuel, StateDriver, EnvironmentDriver) (EngineInstance, error)

	// CompileManifest generates a LogicDescriptor from a Manifest, which can then be used to generate
	// a Logic object. The fuel spent during compile is returned with any potential error.
	CompileManifest(Manifest, EngineFuel) (LogicDescriptor, EngineFuel, error)

	// ValidateCalldata verifies the calldata and callsite in an InteractionDriver.
	// The Logic must describe a callsite which accepts the calldata.
	ValidateCalldata(LogicDriver, InteractionDriver) error
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

// FetchRuntimeEncoder retrieves the RuntimeEncoder for a given EngineKind.
// Returns false if the EngineRuntime for the expected kind is not already registered.
func FetchRuntimeEncoder(kind EngineKind) (RuntimeEncoder, bool) {
	runtime, exists := registry[kind]
	if !exists {
		return nil, false
	}

	return runtime.GenerateRuntimeEncoder(), true
}
