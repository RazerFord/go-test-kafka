package main

import (
	"log"
	args "testkafka/internal/server/argumentparser"
	"testkafka/internal/server/server"
)

func main() {
	s := server.NewServer(args.Parse(), 0)
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
