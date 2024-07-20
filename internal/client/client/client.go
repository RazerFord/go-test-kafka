package client

import (
	"context"
	msg "testkafka/internal/common/message"
	args "testkafka/internal/server/argumentparser"
	"time"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	ParsedArgs *args.Arguments
}

func NewClient(parsedArgs *args.Arguments) *Client {
	return &Client{ParsedArgs: parsedArgs}
}

func (c *Client) DoWithDeadline(m *msg.Message, waiting time.Duration) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", c.ParsedArgs.Address(), c.ParsedArgs.Topic, c.ParsedArgs.Partition)
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(waiting))

	bs, err := m.ToBytes()
	if err != nil {
		return err
	}

	_, err = conn.WriteMessages(kafka.Message{Value: bs})
	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Do(m *msg.Message) error {
	return c.DoWithDeadline(m, 0)
}
