package client

import (
	"context"
	"sync"
	"sync/atomic"
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

func (c *Client) DoWithDeadline(m *kafka.Message, waiting time.Duration) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(c.ParsedArgs.Address()),
		Topic:    c.ParsedArgs.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	if err := w.WriteMessages(context.Background(), *m); err != nil {
		return err
	}

	return w.Close()
}

func (c *Client) Do(m *kafka.Message) error {
	return c.DoWithDeadline(m, 0)
}

func (c *Client) DoWithFakerDeadline(count uint, creator msg.CreateMessage, waiting time.Duration) error {
	for i := uint(0); i < count; i++ {
		if err := c.DoWithDeadline(creator.Create(), waiting); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) DoWithFakerDeadlineParallel(count uint, creator msg.CreateMessage, waiting time.Duration) error {
	ch := make(chan error, 1)
	defer close(ch)

	failed := atomic.Value{}
	failed.Store(false)

	wg := sync.WaitGroup{}
	wg.Add(int(count))
	for i := uint(0); i < count; i++ {
		select {
		case err := <-ch:
			{
				return err
			}
		default:
			{
				go func() {
					if err := c.DoWithDeadline(creator.Create(), waiting); err != nil {
						if failed.CompareAndSwap(false, true) {
							ch <- err
						}
					}
					wg.Done()
				}()
			}
		}
	}
	wg.Wait()
	select {
	case err := <-ch:
		{
			return err
		}
	default:
		{
			return nil
		}
	}
}
