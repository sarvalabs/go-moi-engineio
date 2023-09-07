package engineio

// registry is an in-memory registry of EngineRuntime instances and their respective crypto drivers.
// Support for different engine runtimes is only available if they are first registered with this package.
var registry = map[EngineKind]entry{}

type entry struct {
	runtime EngineRuntime
	crypto  CryptoDriver
}

// RegisterRuntime registers an EngineRuntime with the package along with a CryptoDriver for the runtime.
// If a runtime instance already exists for the EngineKind, it is overwritten.
func RegisterRuntime(runtime EngineRuntime, crypto CryptoDriver) {
	registry[runtime.Kind()] = entry{runtime, crypto}
}

// FetchEngineRuntime retrieves an EngineRuntime for a given EngineKind.
// If the runtime for the engine kind is not registered, returns false.
func FetchEngineRuntime(kind EngineKind) (EngineRuntime, bool) {
	object, exists := registry[kind]
	if !exists {
		return nil, false
	}

	return object.runtime, true
}

// FetchCryptoDriver retrieves an CryptoDriver for a given EngineKind.
// If the runtime for the engine kind is not registered, returns false.
func FetchCryptoDriver(kind EngineKind) (CryptoDriver, bool) {
	object, exists := registry[kind]
	if !exists {
		return nil, false
	}

	return object.crypto, true
}
