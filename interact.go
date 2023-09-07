package engineio

import "math/big"

// IxnType is an interface that describes the different kinds of
// interactions available on go-moi. Each IxnType has a unique ID
// and a string representation. It is implemented by common.IxType.
type IxnType interface {
	IxnID() int
	String() string
}

// IxnDriver represents a driver for interaction information.
// It describes the callsite and input calldata for execution calls along
// with other information such as the Interaction's fuel parameters or transfer funds.
type IxnDriver interface {
	IxnType() IxnType

	FuelPrice() *big.Int
	FuelLimit() uint64

	Callsite() string
	Calldata() []byte
}
