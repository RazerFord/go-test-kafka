package argumentparser

import (
	"flag"
	"fmt"
)

type Arguments struct {
	Host      string
	Port      string
	Topic     string
	Partition int
}

const (
	host      = "localhost"
	port      = "29092"
	topic     = "test-topic"
	partition = 0
)

func Parse() *Arguments {
	a := Arguments{}

	flag.StringVar(&a.Host, "address", host, "kafka host")
	flag.StringVar(&a.Port, "port", port, "kafka port")
	flag.StringVar(&a.Topic, "topic", topic, "kafka topic")
	flag.IntVar(&a.Partition, "partition", partition, "kafka partition")
	flag.Parse()

	return &a
}

func (a *Arguments) Address() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func (a *Arguments) String() string {
	return fmt.Sprintf("{host:%s port:%s topic:%s partition:%d}", a.Host, a.Port, a.Topic, a.Partition)
}
