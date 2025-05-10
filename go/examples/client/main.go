package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/microai-times/task-protos/go/protos"
)

var nodeId = "node-1"

func handleRpc(c pb.RegistrationClient) {
	for {
		r, err := c.Heartbeat(context.Background(), &pb.NodePulse{NodeId: nodeId})
		if err != nil {
			log.Printf("发送心跳失败: %v\n", err)
			return
		}

		t, err := r.Recv()
		if err != nil {
			log.Printf("接收任务错误: %v\n", err)
			return
		}
		log.Printf("Heartbeat: %v\n", t.TaskId)
	}
}

func handleRegistration(conn *grpc.ClientConn) {
	c := pb.NewRegistrationClient(conn)

	r, err := c.Register(context.Background(), &pb.NodeInfo{NodeId: nodeId})
	if err != nil {
		log.Printf("注册失败: %v\n", err)
		return
	}

	log.Printf("Assigned: %s, %v", r.GetAssignedId(), r.GetSuccess())

	handleRpc(c)
}

func main() {
	addr := "localhost:7060"
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()

	for {
		handleRegistration(conn)
		time.Sleep(time.Second * 15)
	}

}
