package netio

import (
	"encoding/json"
)

type MessageLine struct {
	//type of message i.e. the communications protocol
	MessageType string
	//Specific message command
	Command string
	//any data, can be empty. gets interpreted downstream to other structs
	Data json.RawMessage `json:"data,omitempty"`
}

// func NewLineNMessage(m Message) (MessageLine, error) {
// 	return
// }
