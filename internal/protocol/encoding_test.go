package protocol

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ptr(s string) *string {
	return &s
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

type failWriter struct{}

func (f *failWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write failed")
}

func TestJSONEncoder(t *testing.T) {
	t.Run("correctly encodes message", func(t *testing.T) {
		// given
		encoder := NewJSONEncoder()
		buf := bytes.Buffer{}
		msg := Message{
			Type:      Chat,
			Sender:    "Gabe",
			Room:      "Room1",
			Content:   ptr("Hello World"),
			Timestamp: mustParseTime("2026-01-14T19:22:10Z"),
			Metadata:  make(map[string]any),
		}

		// when
		err := encoder.Encode(&buf, &msg)

		// then
		assert.NoError(t, err)
		assert.Equal(t, `{"type":"chat","sender":"Gabe","room":"Room1","content":"Hello World","timestamp":"2026-01-14T19:22:10Z"}`, strings.TrimSpace(buf.String()))
	})

	t.Run("throws when writer fails to write", func(t *testing.T) {
		// given
		encoder := NewJSONEncoder()
		msg := Message{
			Type:      Chat,
			Sender:    "Gabe",
			Room:      "Room1",
			Content:   ptr("Hello World"),
			Timestamp: mustParseTime("2026-01-14T19:22:10Z"),
			Metadata:  make(map[string]any),
		}

		// when
		err := encoder.Encode(&failWriter{}, &msg)

		// then
		assert.Error(t, err)
	})
}

func TestJSONDecoder(t *testing.T) {
	t.Run("correctly decodes message", func(t *testing.T) {
		// given
		encoder := NewJSONEncoder()
		reader := strings.NewReader(`{"type":"chat","sender":"Gabe","room":"Room1","content":"Hello World","timestamp":"2026-01-14T19:22:10Z"}`)
		msg := Message{}

		// when
		err := encoder.Decode(reader, &msg)

		// then
		assert.NoError(t, err)
		assert.Equal(t, Message{
			Type:      Chat,
			Sender:    "Gabe",
			Room:      "Room1",
			Content:   ptr("Hello World"),
			Timestamp: mustParseTime("2026-01-14T19:22:10Z"),
		}, msg)
	})

	t.Run("throws when message is incorrect", func(t *testing.T) {
		// given
		encoder := NewJSONEncoder()
		reader := strings.NewReader(`{"type":"invalid","sender":"Gabe","room":"Room1","content":"Hello World","timestamp":"2026-01-14T19:22:10Z"}`)
		msg := Message{}

		// when
		err := encoder.Decode(reader, &msg)

		// then
		assert.Error(t, err)
	})
}
