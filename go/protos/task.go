package protos

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "microai.com/protos/gen/protos"
)

// 核心数据结构
type TaskNode struct {
	ID           string
	LastActive   time.Time
	TaskChannel  chan *pb.Task
	Capabilities *pb.DeviceCapabilities
}

type TaskCenter struct {
	pb.UnimplementedRegistrationServer
	nodes sync.Map // Concurrent map: node_id -> *AndroidNode
}

// 注册实现
func (s *TaskCenter) Register(ctx context.Context, req *pb.NodeInfo) (*pb.RegistrationResponse, error) {
	node := &TaskNode{
		ID:           req.NodeId,
		TaskChannel:  make(chan *pb.Task, 10), // 带缓冲的通道
		Capabilities: req.Capabilities,
	}

	s.nodes.Store(req.NodeId, node)
	return &pb.RegistrationResponse{
		Success:           true,
		HeartbeatInterval: 30, // 建议心跳间隔(秒)
	}, nil
}

// 心跳长连接实现
func (s *TaskCenter) Heartbeat(req *pb.NodePulse, stream pb.Registration_HeartbeatServer) error {
	val, ok := s.nodes.Load(req.NodeId)
	if !ok {
		return status.Error(codes.NotFound, "node not registered")
	}

	node := val.(*TaskNode)
	for {
		select {
		case task := <-node.TaskChannel:
			if err := stream.Send(task); err != nil {
				s.nodes.Delete(req.NodeId)
				return err
			}
		case <-time.After(30 * time.Second):
			// 保持连接活跃
			if err := stream.Send(&pb.Task{TaskId: "PING"}); err != nil {
				return err
			}
		}
	}
}

// 基于设备能力的任务分配
func (s *TaskCenter) DispatchTask(task *pb.Task) error {
	var selectedNode *TaskNode

	// s.nodes.Range(func(key, value interface{}) bool {
	// 	node := value.(*TaskNode)
	// 	if node.Capabilities.GetMemoryMb() >= task.RequiredMemory {
	// 		selectedNode = node
	// 		return false // 终止遍历
	// 	}
	// 	return true
	// })

	if selectedNode != nil {
		select {
		case selectedNode.TaskChannel <- task:
			return nil
		default:
			return errors.New("node task queue full")
		}
	}
	return errors.New("no suitable node available")
}

// Go服务端健康检查
func (s *TaskCenter) checkNodeHealth() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.nodes.Range(func(key, value interface{}) bool {
			node := value.(*TaskNode)
			if time.Since(node.LastActive) > 5*time.Minute {
				s.nodes.Delete(key)
				log.Printf("Evicted inactive node: %s", key)
			}
			return true
		})
	}
}

// Go服务端压缩
func compressTask(task *pb.Task) {
	if len(task.Payload) > 1024 { // 超过1KB压缩
		var buf bytes.Buffer
		z := zlib.NewWriter(&buf)
		z.Write(task.Payload)
		z.Close()
		task.Payload = buf.Bytes()
		task.Compression = pb.CompressionType_ZLIB
	}
}
