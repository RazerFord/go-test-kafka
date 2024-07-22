package message

import (
	"bytes"
	"fmt"
	"reflect"
	"testkafka/generated/messageproto"

	"github.com/go-faker/faker/v4"
	"google.golang.org/protobuf/proto"
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

		s := buff.String()
		if len(s) > 0 {
			s = s[:len(s)-1]
		}

		return s, nil
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
	return proto.Marshal(&messageproto.Message{From: m.From, To: m.To, Body: m.Body})
}

func FromBytes(buff []byte) (*Message, error) {
	m := messageproto.Message{}

	if err := proto.Unmarshal(buff, &m); err != nil {
		return nil, err
	}

	return &Message{From: m.From, To: m.To, Body: m.Body}, nil
}

func (m *Message) String() string {
	return fmt.Sprintf("{from:%s to:%s body:%s}", m.From, m.To, m.Body)
}
