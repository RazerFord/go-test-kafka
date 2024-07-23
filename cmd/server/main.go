package main

import (
	"context"
	"flag"
	"fmt"
	pb "gotestgrpc/generated/sortingroto"
	"log"
	"net"
	"sort"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSorterServer
}

type Ints []int32

func (a Ints) Len() int           { return len(a) }
func (a Ints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Ints) Less(i, j int) bool { return a[i] < a[j] }

func (s *server) Sort(_ context.Context, req *pb.Request) (*pb.Response, error) {
	ints := req.GetData()
	sort.Sort(Ints(ints))
	return &pb.Response{Data: ints}, nil
}

func main() {
	port := 0
	flag.IntVar(&port, "port", 0, "server port")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalln(err)
	}
	defer lis.Close()
	log.Printf("hosted on %s\n", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterSorterServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
