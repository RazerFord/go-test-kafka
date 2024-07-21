package argumentparser

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Arguments struct {
	Host         string
	Port         string
	Topic        string
	ServerPort   string
	Partition    int
	CountMessage int
	Waiting      int
	Sleep        int
	Repeat       bool
}

const (
	host         = "localhost"
	port         = "29092"
	topic        = "test-topic"
	serverPort   = "8081"
	partition    = 0
	countMessage = 10
	waiting      = 10
	sleep        = 10
	repeat       = false
)

func Parse() *Arguments {
	a := Arguments{}

	flag.StringVar(&a.Host, "host", host, "kafka host")
	if s, r := os.LookupEnv("host"); r {
		a.Host = s
	}

	flag.StringVar(&a.Port, "port", port, "kafka port")
	if s, r := os.LookupEnv("port"); r {
		a.Port = s
	}

	flag.StringVar(&a.Topic, "topic", topic, "kafka topic")
	if s, r := os.LookupEnv("topic"); r {
		a.Topic = s
	}

	flag.StringVar(&a.ServerPort, "server", serverPort, "server port")
	if s, r := os.LookupEnv("server"); r {
		a.ServerPort = s
	}

	flag.IntVar(&a.Partition, "partition", partition, "kafka partition")
	if s, r := os.LookupEnv("partition"); r {
		a.Partition, _ = strconv.Atoi(s)
	}

	flag.IntVar(&a.CountMessage, "count", countMessage, "count of messages sent")
	if s, r := os.LookupEnv("count"); r {
		a.CountMessage, _ = strconv.Atoi(s)
	}

	flag.IntVar(&a.Waiting, "duration", waiting, "message sending time")
	if s, r := os.LookupEnv("duration"); r {
		a.Waiting, _ = strconv.Atoi(s)
	}

	flag.IntVar(&a.Sleep, "sleep", sleep, "sleep after sending messages")
	if s, r := os.LookupEnv("sleep"); r {
		a.Sleep, _ = strconv.Atoi(s)
	}

	flag.BoolVar(&a.Repeat, "repeat", repeat, "resend")
	if s, r := os.LookupEnv("repeat"); r {
		a.Repeat = s == "true"
	}

	flag.Parse()

	return &a
}

func (a *Arguments) Address() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func (a *Arguments) String() string {
	return fmt.Sprintf(
		"{host:%s port:%s topic:%s server-port:%s partition:%d count-message:%d waiting:%d repeat:%t}",
		a.Host,
		a.Port,
		a.Topic,
		a.ServerPort,
		a.Partition,
		a.CountMessage,
		a.Waiting,
		a.Repeat,
	)
}
