package protocol

import (
	"encoding/json"
	"fmt"
	"time"
)

type MessageType int

const (
	Chat MessageType = iota
	Join
	Leave
	Typing
)

func (m MessageType) String() string {
	switch m {
	case Chat:
		return "chat"
	case Join:
		return "join"
	case Leave:
		return "leave"
	case Typing:
		return "typing"
	default:
		return ""
	}
}

func (mt MessageType) MarshalJSON() ([]byte, error) {
	mtToStrMap := map[MessageType]string{
		Chat:   "chat",
		Join:   "join",
		Leave:  "leave",
		Typing: "typing",
	}
	value, ok := mtToStrMap[mt]
	if !ok {
		return nil, fmt.Errorf("unhandled message type %d", mt)
	}
	return []byte(`"` + value + `"`), nil // could have used json.Marshal(value) but that is more robust

}

func (mt *MessageType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	strToMtMap := map[string]MessageType{
		"chat":   Chat,
		"join":   Join,
		"leave":  Leave,
		"typing": Typing,
	}
	value, ok := strToMtMap[str]
	if !ok {
		return fmt.Errorf("unhandled message type %s", str)
	}
	*mt = value
	return nil
}

// Message represents a single message in the system
type Message struct {
	// The core message type, could be ("chat", "join", "leave", "typing", etc...)
	Type MessageType `json:"type"`
	// The unique identifier of the user who sent the message
	Sender string `json:"sender"`
	// Room is the identifier where the message belongs
	Room string `json:"room"`
	// The content of the message
	Content *string `json:"content"`
	// When the message was created
	Timestamp time.Time `json:"timestamp"`
	// Metadata field for future use when sending data not just messages
	Metadata map[string]any `json:"metadata,omitempty"`
}
