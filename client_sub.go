package main

import (
	"log"
	pb "github.com/jays-computing/gosync/gosync"
	"google.golang.org/grpc"
	"context"
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
	log.Println("Successfully joined session: ", session.SessionName)

	events, err := client.GetEvents(context.Background(), &pb.GetEventsRequest{Session: session})
	if err != nil {
		log.Fatal("Failed to subscribe to events")
	}
	for {
		in, err := events.Recv()
		if err != nil {
			log.Fatal("Failed to receive from event stream:" , err)
			return
		}
		log.Println("Received message: ", in.Message)

	}
}
