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

	m := message.NewMessage("Egor", "Egor2", "Hi, Egor2!")

	if err := c.DoWithDeadline(m, time.Duration(10)*time.Second); err != nil {
		log.Fatalln(err)
	}
}
