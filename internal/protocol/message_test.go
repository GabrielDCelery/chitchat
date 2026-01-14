package protocol

import (
	"encoding/json"
	"testing"
	"time"

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
				// when
				result, err := tt.input.MarshalJSON()

				// then
				assert.NoError(t, err)
				assert.Equal(t, result, []byte(tt.expected))
			})
		}
	})

	t.Run("throws if it sees incorrect message type", func(t *testing.T) {
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
		// given
		var mt MessageType

		// when
		err := mt.UnmarshalJSON([]byte("invalidvalue"))

		// then
		assert.Error(t, err)
	})
}

func TestMessageMarshalJSON(t *testing.T) {
	t.Run("correctly marshals message", func(t *testing.T) {
		// given
		content := "Hello World"
		message := Message{
			Type:      Chat,
			Sender:    "Gabe",
			Room:      "Room1",
			Content:   &content,
			Timestamp: time.Date(2026, time.January, 14, 19, 22, 10, 0, time.UTC),
			Metadata:  make(map[string]any),
		}

		// when
		parsed, err := json.Marshal(message)

		// then
		assert.NoError(t, err)
		assert.Equal(t, parsed, []byte(`{"type":"chat","sender":"Gabe","room":"Room1","content":"Hello World","timestamp":"2026-01-14T19:22:10Z"}`))
	})
}

func TestMessageUnmarshalJSON(t *testing.T) {
	t.Run("correctly unmarshals message", func(t *testing.T) {
		// given
		message := `{"type":"chat","sender":"Gabe","room":"Room1","content":"Hello World","timestamp":"2026-01-14T19:22:10Z"}`

		// when
		m := Message{}
		err := json.Unmarshal([]byte(message), &m)

		// then
		content := "Hello World"
		assert.NoError(t, err)
		assert.Equal(t, m, Message{
			Type:      Chat,
			Sender:    "Gabe",
			Room:      "Room1",
			Content:   &content,
			Timestamp: time.Date(2026, time.January, 14, 19, 22, 10, 0, time.UTC),
		})
	})
}
