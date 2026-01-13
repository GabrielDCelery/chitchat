package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageTypeMarshalJSON(t *testing.T) {
	t.Run("correctly marshals message type", func(t *testing.T) {
		tests := []struct {
			name     string
			input    MessageType
			expected string
		}{
			{"chat", Chat, `"chat"`},
			{"join", Join, `"join"`},
			{"leave", Leave, `"leave"`},
			{"typing", Typing, `"typing"`},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// when
				result, err := tt.input.MarshalJSON()

				// then
				assert.NoError(t, err)
				assert.Equal(t, result, []byte(tt.expected))
			})
		}
	})

	t.Run("throws if it sees incorrect message type", func(t *testing.T) {
		t.Parallel()
		// given
		mt := MessageType(999)

		// when
		_, err := mt.MarshalJSON()

		// then
		assert.Error(t, err)
	})
}

func TestMessageTypeUnmarshalJSON(t *testing.T) {
	t.Run("correctly unmarshals message type", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected MessageType
		}{
			{"chat", `"chat"`, Chat},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// given
				var mt MessageType

				// when
				err := mt.UnmarshalJSON([]byte(tt.input))

				// then
				assert.NoError(t, err)
				assert.Equal(t, mt, tt.expected)
			})
		}
	})

	t.Run("throws if it sees incorrect message type", func(t *testing.T) {
		t.Parallel()
		// given
		var mt MessageType

		// when
		err := mt.UnmarshalJSON([]byte("invalidvalue"))

		// then
		assert.Error(t, err)
	})
}
