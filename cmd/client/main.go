package main

import (
	"log"
	"testkafka/internal/client/client"
	"testkafka/internal/common/message"
	args "testkafka/internal/server/argumentparser"
	"time"
)

func main() {
	c := client.NewClient(args.Parse())

	gen := message.NewMessageGen(
		message.NewGenerator(
			message.NewKeyGenFromList([]string{"common", "personal"}, message.RandomIndex(2)).Gen,
			message.NewValueGen().Gen,
		),
	)

	if err := c.DoWithFakerDeadlineParallel(10, gen, time.Duration(10)*time.Second); err != nil {
		log.Fatalln(err)
	}
}
