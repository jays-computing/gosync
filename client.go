package main

import (
	"log"
	"time"
	//"golang.org/x/net/context"
	pb "github.com/jays-computing/gosync/gosync"
	"google.golang.org/grpc"
	//"os"
	"context"
	//"io"
)

const (
	address = "127.0.0.1:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGoSyncClient(conn)
	session, err := client.JoinSession(context.Background())
	if err != nil {
		log.Fatal("Error joining session")
	}

	m := pb.GSMessage{
		Message:       "HelOoo",
		Time:          int64(time.Now().Nanosecond()),
		ExecuteOffset: 0,
	}
	if err := session.Send(&m); err != nil {
		log.Fatal("Error sending data")
	}
	in, err := session.Recv()
	print(in.Message)
	if err != nil {
		log.Fatal("Error receiving data ")
	}
	if err := session.CloseSend(); err != nil {
		log.Fatal("Failed to send EOF")
	}
}
