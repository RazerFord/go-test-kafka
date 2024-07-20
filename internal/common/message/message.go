package message

import "encoding/json"

type Message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}

func NewMessage(from, to, body string) *Message {
	return &Message{from, to, body}
}

func (m *Message) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func FromBytes(buff []byte) (*Message, error) {
	m := Message{}

	if err := json.Unmarshal(buff, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
