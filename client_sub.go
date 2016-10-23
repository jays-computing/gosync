package main

import (
	"log"
	gs "github.com/jays-computing/gosync/gosync"
	"google.golang.org/grpc"
	"context"
	"github.com/jays-computing/gosync/utils"
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
	log.Println("Successfully joined session: ",joinReply.Session.SessionName)
	guestId := joinReply.GuestId

	rttAverage, err := utils.GetAverageNtpRtt(joinReply.Session.NtpHost, 3)
	if err != nil {
		log.Fatal("Error getting average ntp", err)
	}

	events, err := client.GetEvents(context.Background(), &gs.GetEventsRequest{Session:joinReply.Session, GuestId: guestId, NtpTimeRtt: rttAverage})
	if err != nil {
		log.Fatal("Failed to subscribe to events")
	}
	for {
		in, err := events.Recv()
		if err != nil {
			log.Fatal("Failed to receive from event stream:" , err)
			return
		}
		log.Println("Received message: ", in.Message, " time: " , in.Time)

	}
}
