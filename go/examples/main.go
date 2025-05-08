package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/thebigbrain/microai-protos/protos"
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
	if err := protos.StartServer(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)

	}
}
