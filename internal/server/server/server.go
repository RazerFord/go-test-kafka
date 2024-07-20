package server

import (
	"context"
	"fmt"
	"log"
	args "testkafka/internal/server/argumentparser"
	"time"

	"github.com/segmentio/kafka-go"
)

type Server struct {
	ParsedArgs *args.Arguments
}

func NewServer(parsedArgs *args.Arguments) *Server {
	return &Server{ParsedArgs: parsedArgs}
}

func (s *Server) Run() error {
	log.Println(s)
	conn, err := kafka.DialLeader(context.Background(), "tcp", s.ParsedArgs.Address(), s.ParsedArgs.Topic, s.ParsedArgs.Partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	for {
		n, err := batch.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(n.Value)
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}

	return nil
}
