package engineio

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sarvalabs/go-polo"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestContextStateMatrix(t *testing.T) {
	tests := []struct {
		object     ContextStateMatrix
		persistent bool
		ephemeral  bool
	}{
		{ContextStateMatrix{PersistentState: 0, EphemeralState: 1}, true, true},
		{ContextStateMatrix{PersistentState: 0}, true, false},
		{ContextStateMatrix{EphemeralState: 1}, false, true},
		{ContextStateMatrix{}, false, false},
	}

	for _, test := range tests {
		require.Equal(t, test.persistent, test.object.Persistent())
		require.Equal(t, test.ephemeral, test.object.Ephemeral())
	}
}

func TestContextScopeKind_String(t *testing.T) {
	require.True(t, len(contextStateKindFromString) == len(contextStateKindToString))

	require.Equal(t, "persistent", PersistentState.String())
	require.Equal(t, "ephemeral", EphemeralState.String())

	require.PanicsWithValue(t, "unknown ContextStateKind variant", func() {
		_ = ContextStateKind(10).String()
	})
}

//nolint:dupl
func TestContextStateKind_Serialization(t *testing.T) {
	t.Run("POLO", func(t *testing.T) {
		for kind, str := range contextStateKindToString {
			encoded, err := polo.Polorize(kind)

			require.Nil(t, err)
			require.Equal(t, bytes.Join([][]byte{{6}, []byte(str)}, []byte{}), encoded)

			decoded := new(ContextStateKind)
			err = polo.Depolorize(decoded, encoded)

			require.Nil(t, err)
			require.Equal(t, kind, *decoded)
		}

		require.PanicsWithValue(t, "unknown ContextStateKind variant", func() {
			_, _ = polo.Polorize(ContextStateKind(10))
		})

		malformed := bytes.Join([][]byte{{6}, []byte("malformed")}, []byte{})
		err := polo.Depolorize(new(ContextStateKind), malformed)
		require.EqualError(t, err, "invalid ContextStateKind value")
	})

	t.Run("JSON", func(t *testing.T) {
		for kind, str := range contextStateKindToString {
			encoded, err := json.Marshal(kind)

			require.Nil(t, err)
			require.Equal(t, bytes.Join([][]byte{{34}, []byte(str), {34}}, []byte{}), encoded)

			decoded := new(ContextStateKind)
			err = json.Unmarshal(encoded, decoded)

			require.Nil(t, err)
			require.Equal(t, kind, *decoded)
		}

		require.PanicsWithValue(t, "unknown ContextStateKind variant", func() {
			_, _ = json.Marshal(ContextStateKind(10))
		})

		malformed := bytes.Join([][]byte{{34}, []byte("malformed"), {34}}, []byte{})
		err := json.Unmarshal(malformed, new(ContextStateKind))
		require.EqualError(t, err, "invalid ContextStateKind value")
	})

	t.Run("YAML", func(t *testing.T) {
		for kind, str := range contextStateKindToString {
			encoded, err := yaml.Marshal(kind)

			require.Nil(t, err)
			require.Equal(t, bytes.Join([][]byte{[]byte(str), {10}}, []byte{}), encoded)

			decoded := new(ContextStateKind)
			err = yaml.Unmarshal(encoded, decoded)

			require.Nil(t, err)
			require.Equal(t, kind, *decoded)
		}

		require.PanicsWithValue(t, "unknown ContextStateKind variant", func() {
			_, _ = yaml.Marshal(ContextStateKind(10))
		})

		malformed := bytes.Join([][]byte{{34}, []byte("malformed"), {34}}, []byte{})
		err := yaml.Unmarshal(malformed, new(ContextStateKind))
		require.EqualError(t, err, "invalid ContextStateKind value")
	})
}
