package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-faker/faker/v4"
)

const MaxLength = 10

func init() {
	_ = faker.AddProvider("bodyFaker", func(v reflect.Value) (interface{}, error) {
		buff := bytes.Buffer{}

		for i := 0; i < MaxLength; i++ {
			if _, err := buff.WriteString(faker.Word()); err != nil {
				return nil, err
			}
			if err := buff.WriteByte(' '); err != nil {
				return nil, err
			}
		}

		return buff.String(), nil
	})
}

type Message struct {
	From string `json:"from" faker:"name"`
	To   string `json:"to" faker:"name"`
	Body string `json:"body" faker:"bodyFaker"`
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

func (m *Message) String() string {
	return fmt.Sprintf("{from:%s to:%s body:%s}", m.From, m.To, m.Body)
}
