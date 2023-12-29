package engineio

// CallResult is the output of the Call method of Engine.
// It expresses the amount of engine fuel consumed for the call along with
// result of the call which can either be some outputs or an error response
type CallResult interface {
	// Ok specifies whether the execution call was successful.
	// If true, Error() must return nil.
	Ok() bool
	// Fuel specifies the amount of EngineFuel that was consumed
	// for the execution call regardless of its successful run.
	Fuel() EngineFuel

	// Engine specifies the engine kind
	// that generated the CallResult
	Engine() EngineKind

	// Outputs returns the outputs for the execution call.
	// The output data must be polo document-encoded and
	// may be nil if the call has no return values
	Outputs() []byte

	// Error returns the encoded error for the execution call (if any).
	// Must return a non-nil value if Ok() is false and vice versa.
	// The output bytes must be decodable into an ErrorResult
	// using the DecodeErrorResult method of EngineRuntime
	Error() []byte
}

// ErrorResult is an interface for an engine specific error message.
// It is returned as raw bytes within CallResult if an execution call fails.
//
// It can be decoded from the raw data using the DecodeErrorResult method of EngineRuntime
type ErrorResult interface {
	// Engine specifies the engine kind that generated the ErrorResult
	Engine() EngineKind
	// String returns a string representation of ErrorResult
	String() string
	// Bytes returns the bytes representation of ErrorResult
	Bytes() []byte
	// Reverted returns whether the ErrorResult has a reversion flag
	Reverted() bool
}
