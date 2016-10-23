package main

import (
	"log"
	gs "github.com/jays-computing/gosync/gosync"
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
	client := gs.NewGoSyncClient(conn)
	joinReply, err := client.JoinSession(context.Background(), &gs.JoinRequest{Name:"Jamzz"} )
	if err != nil {
		log.Fatal("Error joining session")
	}
	if err != nil {
		log.Fatal("Error getting server time:", err)
	}

	publishTime := time.Now().Nanosecond()
	_, err = client.PublishEvent(context.Background(), &gs.PublishRequest{
		Session: joinReply.Session,
		Message: &gs.GSMessage{
			Message: "Helllowowowwo",
			Time: int32(publishTime),
		},
	})
	if err != nil {
		log.Fatal("Error publishing event")
	}
	log.Println("Published")

}
