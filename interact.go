package engineio

import "math/big"

// InteractionType is an interface that describes the different kinds
// of interactions available on go-moi. Each InteractionType has a unique
// ID and a string representation. It is implemented by common.IxType.
type InteractionType interface {
	IxnID() int
	String() string
}

// InteractionDriver represents a driver for interaction information.
// It describes the callsite and input calldata for execution calls along with
// other information such as the Interaction's fuel parameters or transfer funds.
type InteractionDriver interface {
	InteractionType() InteractionType

	FuelPrice() *big.Int
	FuelLimit() uint64

	Callsite() string
	Calldata() []byte
}
