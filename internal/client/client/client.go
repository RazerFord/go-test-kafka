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

	mtx := sync.Mutex{}

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
						mtx.Lock()
						if failed.Load().(bool) {
							mtx.Unlock()
							return
						}
						failed.Store(true)
						mtx.Unlock()
						ch <- err
					}
					wg.Done()
				}()
			}
		}
	}
	wg.Wait()
	return nil
}
