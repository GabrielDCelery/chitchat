package protocol

import (
	"encoding/json"
	"io"
)

type Encoder interface {
	Encode(w io.Writer, msg *Message) error
	Decode(r io.Reader, msg *Message) error
}

type JSONEncoder struct{}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

func (j *JSONEncoder) Encode(w io.Writer, msg *Message) error {
	return json.NewEncoder(w).Encode(msg)
	// NOTE: did not know we had an idiomatic way of doing this so kept the code for learning purposes

	// data, err := json.Marshal(msg)
	// if err != nil {
	// 	return err
	// }
	// n, err := w.Write(data)
	// if err != nil {
	// 	return err
	// }
	// if n != len(data) {
	// 	return fmt.Errorf("only sent %d bytes out of %d", n, len(data))
	// }
	// return nil
}

func (j *JSONEncoder) Decode(r io.Reader, msg *Message) error {
	return json.NewDecoder(r).Decode(msg)
	// NOTE: did not know we had an idiomatic way of doing this so kept the code for learning purposes

	// data, err := io.ReadAll(r)
	// if err != nil {
	// 	return err
	// }
	// err = json.Unmarshal(data, msg)
	// if err != nil {
	// 	return err
	// }
	// return nil
}
