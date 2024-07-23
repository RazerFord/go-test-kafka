package main

import (
	"context"
	"flag"
	"fmt"
	pb "gotestgrpc/generated/sortingroto"
	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := 0
	length := 0
	numberСlients := 1
	flag.IntVar(&port, "port", 0, "server port")
	flag.IntVar(&length, "length", 1024, "array length")
	flag.IntVar(&numberСlients, "clients", 1, "number of clients")
	flag.Parse()

	arr := make([]int32, length)
	for i := 0; i < length; i++ {
		arr[i] = rand.Int31()
	}

	wg := sync.WaitGroup{}
	wg.Add(numberСlients)

	for i := 0; i < numberСlients; i++ {
		go func(n int) {
			defer wg.Done()

			conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()

			c := pb.NewSorterClient(conn)

			req := pb.Request{Data: arr}
			resp, err := c.Sort(context.Background(), &req)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("%d: %v\n", n, resp.Data)
		}(i)
	}
	wg.Wait()
}
