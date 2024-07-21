package message

import (
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/segmentio/kafka-go"
)

const (
	StandardKey   = "key"
	StandardValue = "value"
)

type FuncGen interface {
	Gen() []byte
}

type AdapterFuncGen func() []byte

///////////////////////////////////////////////////////////////////

type Generator struct {
	KeyGen   AdapterFuncGen
	ValueGen AdapterFuncGen
}

func NewGenerator(keyGen, valueGen AdapterFuncGen) *Generator {
	return &Generator{KeyGen: keyGen, ValueGen: valueGen}
}

func NewStandard() *Generator {
	return NewGenerator(
		func() []byte { return []byte(StandardKey) },
		func() []byte { return []byte(StandardValue) },
	)
}

///////////////////////////////////////////////////////////////////

type storageIndex struct {
	curIndex int
	maxIndex int
}

func (s *storageIndex) get() int {
	r := s.curIndex
	s.curIndex = (s.curIndex + 1) % s.maxIndex
	return r
}

type KeyGenFromList struct {
	Keys     [][]byte
	genIndex GetIndex
}

type GetIndex func() int

func NewKeyGenFromList(keys []string, gi GetIndex) *KeyGenFromList {
	ks := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		ks[i] = []byte(keys[i])
	}

	if gi == nil {
		si := storageIndex{0, len(keys)}
		gi = si.get
	}

	return &KeyGenFromList{ks, gi}
}

func (k *KeyGenFromList) Gen() []byte {
	return k.Keys[k.genIndex()]
}

func RandomIndex(l int) GetIndex {
	return func() int {
		return rand.Int() % l
	}
}

///////////////////////////////////////////////////////////////////

type ValueGen struct {
}

func NewValueGen() *ValueGen {
	return &ValueGen{}
}

func (v *ValueGen) Gen() []byte {
	msg := Message{}
	faker.FakeData(&msg)
	bs, _ := msg.ToBytes()
	return bs
}

///////////////////////////////////////////////////////////////////

type CreateMessage interface {
	Create() *kafka.Message
}

type MessageGen struct {
	*Generator
}

func NewMessageGen(g *Generator) *MessageGen {
	return &MessageGen{g}
}

func (m *MessageGen) Create() *kafka.Message {
	key := m.KeyGen()
	value := m.ValueGen()
	return &kafka.Message{Key: key, Value: value}
}
