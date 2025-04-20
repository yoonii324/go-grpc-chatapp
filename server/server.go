package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "grpcchat/chatproto"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

func main() {
	flag.Parse()

	// wait for the connection
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// register Chat service
	pb.RegisterChatServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

func newServer() *chatServer {
	return &chatServer{}
}

// object which implemented Chat service interface
type chatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	streams []pb.ChatService_ChatServer
}

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
	// add to streams list
	s.mu.Lock()
	s.streams = append(s.streams, stream)
	s.mu.Unlock()

	var err error

	for {
		// read the value sent to client
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		// broadcast to all clients
		s.mu.Lock()
		for _, strm := range s.streams {
			strm.Send(&pb.ChatMsg{
				Sender:  in.Sender,
				Message: in.Message,
			})
		}

		s.mu.Unlock()
	}

	// delete from the list after disconnection
	s.mu.Lock()
	for i, strm := range s.streams {
		if strm == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			break
		}
	}

	s.mu.Unlock()

	return err
}
