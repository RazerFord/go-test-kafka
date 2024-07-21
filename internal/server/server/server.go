package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	args "testkafka/internal/server/argumentparser"
	"time"

	"github.com/segmentio/kafka-go"
)

type Server struct {
	ParsedArgs *args.Arguments
	Timeout    time.Duration
	mtx        sync.RWMutex
	storage    []string
}

func NewServer(parsedArgs *args.Arguments, timeout time.Duration) *Server {
	return &Server{ParsedArgs: parsedArgs, Timeout: timeout}
}

func (s *Server) Run() error {
	go s.startServer()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{s.ParsedArgs.Address()},
		Topic:   s.ParsedArgs.Topic,
	})

	for {
		ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
		m, err := r.ReadMessage(ctx)
		cancel()
		if err != nil {
			break
		}

		s.mtx.Lock()
		s.storage = append(s.storage, string(m.Value))
		s.mtx.Unlock()
	}

	return r.Close()
}

func (s *Server) startServer() {
	localServer := http.NewServeMux()
	localServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		s.mtx.RLock()
		if data, err :=  json.Marshal(s.storage); err == nil {
			w.Write(data)
		}
		s.mtx.RUnlock()
	})
	http.ListenAndServe(fmt.Sprintf(":%s", s.ParsedArgs.ServerPort), localServer)
}
