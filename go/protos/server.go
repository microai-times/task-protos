package protos

import (
	"net"

	pb "github.com/thebigbrain/microai-protos/go/gen"
	"google.golang.org/grpc"
)

func StartServer(lis net.Listener) (err error) {
	s := grpc.NewServer()
	pb.RegisterRegistrationServer(s, &TaskCenter{})

	return s.Serve(lis)
}
