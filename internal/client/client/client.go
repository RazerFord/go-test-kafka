package client

import (
	"context"
	msg "testkafka/internal/common/message"
	args "testkafka/internal/server/argumentparser"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/segmentio/kafka-go"
)

type Client struct {
	ParsedArgs *args.Arguments
}

func NewClient(parsedArgs *args.Arguments) *Client {
	return &Client{ParsedArgs: parsedArgs}
}

func (c *Client) DoWithDeadline(m *kafka.Message, waiting time.Duration) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", c.ParsedArgs.Address(), c.ParsedArgs.Topic, c.ParsedArgs.Partition)
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(waiting))

	_, err = conn.WriteMessages(*m)
	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Do(m *kafka.Message) error {
	return c.DoWithDeadline(m, 0)
}

func (c *Client) DoWithFakerDeadline(count uint, waiting time.Duration) error {
	m := &msg.Message{}
	for i := uint(0); i < count; i++ {
		faker.FakeData(m)
		bs, err := m.ToBytes()
		if err != nil {
			return err
		}
		if err := c.DoWithDeadline(&kafka.Message{Value: bs}, waiting); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) DoWithFakerDeadlineParallel(count uint, waiting time.Duration) error {
	ch := make(chan error, 1)
	for i := uint(0); i < count; i++ {
		select {
		case err := <-ch:
			{
				return err
			}
		default:
			{
				go func() {
					m := &msg.Message{}
					faker.FakeData(m)
					bs, err := m.ToBytes()
					if err != nil {
						ch <- err
					} else if err := c.DoWithDeadline(&kafka.Message{Value: bs}, waiting); err != nil {
						ch <- err
					}
				}()
			}
		}
	}
	return nil
}
