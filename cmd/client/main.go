package main

import (
	"log"
	"testkafka/internal/client/client"
	"testkafka/internal/common/message"
	args "testkafka/internal/server/argumentparser"
	"time"
)

func main() {
	parsed := args.Parse()
	c := client.NewClient(parsed)

	gen := message.NewMessageGen(
		message.NewGenerator(
			message.NewKeyGenFromList(
				[]string{"common", "personal"},
				message.RandomIndex(2),
			).Gen,
			message.NewValueGen().Gen,
		),
	)

	for {
		if err := c.DoWithFakerDeadlineParallel(
			uint(parsed.CountMessage),
			gen,
			time.Duration(parsed.Waiting)*time.Second,
		); err != nil {
			log.Fatalln(err)
		}
		if !parsed.Repeat {
			break
		}
		time.Sleep(time.Duration(parsed.Sleep) * time.Second)
	}
}
