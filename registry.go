package engineio

// registry is an in-memory registry of EngineRuntime instances and their respective cryptography drivers.
// Support for different engine runtimes is only available if they are first registered with this package.
var registry = map[EngineKind]entry{}

type entry struct {
	runtime EngineRuntime
	crypto  CryptographyDriver
}

// RegisterRuntime registers an EngineRuntime with the package along with a CryptographyDriver for the runtime.
// If a runtime instance already exists for the EngineKind, it is overwritten.
func RegisterRuntime(runtime EngineRuntime, crypto CryptographyDriver) {
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

// FetchCryptographyDriver retrieves an CryptoDriver for a given EngineKind.
// If the runtime for the engine kind is not registered, returns false.
func FetchCryptographyDriver(kind EngineKind) (CryptographyDriver, bool) {
	object, exists := registry[kind]
	if !exists {
		return nil, false
	}

	return object.crypto, true
}
