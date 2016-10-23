package main

import (
	pb "github.com/jays-computing/gosync/gosync"
	"google.golang.org/grpc"
	"log"
	"net"
	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type server struct {
	sessions map[string] *pb.Session
	sessionGuests map[string] [] *pb.GoSync_GetEventsServer
}
var online = true

func (s *server) JoinSession(context context.Context, jr *pb.JoinRequest) (*pb.Session, error) {
	println (jr.Name, " joined!!")
	return ServerSessions.sessions["1"], nil
}

func (s *server) GetEvents(ger *pb.GetEventsRequest, stream pb.GoSync_GetEventsServer) error {
	println("Getting events for session Id: " , ger.Session.SessionId)
	ServerSessions.sessionGuests[ger.Session.SessionId] = append(ServerSessions.sessionGuests[ger.Session.SessionId], &stream)
	for online {

	}
	return nil
}

func (s *server) PublishEvent(context context.Context, pr *pb.PublishRequest) (*pb.PublishResult , error) {
	guests := ServerSessions.sessionGuests[pr.Session.SessionId]
	for _, g := range guests {
		if err := (*g).Send(pr.Message); err != nil {
			log.Println("Failed to send to guest ")
		}
	}
	return &pb.PublishResult{
	}, nil
}

func pd() {
	for _, s := range ServerSessions.sessions {
		println("Session: " , s.SessionId, " --name: " + s.SessionName)
		for range ServerSessions.sessionGuests[s.SessionId] {
			println("Guest ")
		}
	}
}
var ServerSessions server;
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	ServerSessions = server {
		sessions:make(map[string]*pb.Session),
		sessionGuests:make(map[string][]*pb.GoSync_GetEventsServer),
	}
	ServerSessions.sessions["1"] = &pb.Session{
		SessionId: "1",
		SessionName: "Sessh",
	}
	ServerSessions.sessionGuests["1"] = make([]*pb.GoSync_GetEventsServer, 0)
	pd()
	pb.RegisterGoSyncServer(s, &server{})
	s.Serve(lis)
}
