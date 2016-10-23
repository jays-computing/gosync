package main

import (
	"log"
	pb "github.com/jays-computing/gosync/gosync"
	"google.golang.org/grpc"
	"context"
	"time"
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
	session, err := client.JoinSession(context.Background(), &pb.JoinRequest{Name:"Jamzz"} )
	if err != nil {
		log.Fatal("Error joining session")
	}

	_, err = client.PublishEvent(context.Background(), &pb.PublishRequest{
		Session: session,
		Message: &pb.GSMessage{
			Message: "Helllowowowwo",
			Time: int32(time.Now().Nanosecond()),
		},
	})
	if err != nil {
		log.Fatal("Error publishing event")
	}
	log.Println("Published")

}
