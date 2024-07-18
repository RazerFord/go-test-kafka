package main

import (
	"log"
	args "testkafka/internal/server/argumentparser"
	"testkafka/internal/server/server"
)

func main() {
	s := server.NewServer(args.Parse())
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
