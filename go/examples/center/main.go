package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/microai-times/task-protos/go/examples/center/task"
	pb "github.com/microai-times/task-protos/go/protos"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 7060, "The server port")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterRegistrationServer(s, &task.TaskCenter{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
