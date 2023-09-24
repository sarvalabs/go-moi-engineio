package engineio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRefProvider map[string]any

func (m mockRefProvider) GetReference(ref ReferenceValue) (any, bool) {
	val, ok := m[string(ref)]

	return val, ok
}

func TestReferenceVal_String(t *testing.T) {
	ref := ReferenceValue("foo")
	assert.Equal(t, "ref<foo>", ref.String())
}

func TestEncodeValues(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		refs   ReferenceProvider
		output []byte
		err    string
	}{
		{
			name:   "encode simple int",
			input:  100,
			output: []byte{0x03, 0x64},
			err:    "",
		},
		{
			name:   "encode simple string",
			input:  "hello world",
			output: []byte{0x6, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64},
			err:    "",
		},
		{
			name:   "encode reference",
			input:  ReferenceValue("myref"),
			refs:   &mockRefProvider{"myref": 100},
			output: []byte{0x03, 0x64},
			err:    "",
		},
		{
			name:   "encode reference without provider",
			input:  ReferenceValue("myref"),
			output: nil,
			err:    "encountered reference value without a ref provider",
		},
		{
			name: "encode object",
			input: map[string]any{
				"foo": 1,
				"bar": 2,
			},
			output: []byte{0xd, 0x5f, 0x6, 0x35, 0x56, 0x85, 0x1, 0x62, 0x61, 0x72, 0x3, 0x2, 0x66, 0x6f, 0x6f, 0x3, 0x1},
			err:    "",
		},
		{
			name:   "encode slice",
			input:  []any{1, 2},
			output: []byte{0xe, 0x2f, 0x3, 0x13, 0x1, 0x2},
			err:    "",
		},
		{
			name: "encode map",
			input: map[any]any{
				"foo": 1,
				"bar": 2,
			},
			output: []byte{0xe, 0x4f, 0x6, 0x33, 0x46, 0x73, 0x62, 0x61, 0x72, 0x2, 0x66, 0x6f, 0x6f, 0x1},
			err:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := EncodeValues(test.input, test.refs)

			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.output, output)
			}
		})
	}
}
