package main

import (
	pb "github.com/jays-computing/gosync/gosync"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct {
	messages map[string][]*pb.GSMessage
}

func (s *server) JoinSession(stream pb.GoSync_JoinSessionServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		println(in.Message)
		println("time ", in.Time)
		if err := stream.Send(in); err != nil {
			return err
		}

	}
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGoSyncServer(s, &server{})
	s.Serve(lis)
}
