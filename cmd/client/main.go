package main

import (
	"log"
	"testkafka/internal/client/client"
	args "testkafka/internal/server/argumentparser"
	"time"
)

func main() {
	c := client.NewClient(args.Parse())

	if err := c.DoWithFakerDeadline(10, time.Duration(10)*time.Second); err != nil {
		log.Fatalln(err)
	}
}
