package main

import (
	gs "github.com/jays-computing/gosync/gosync"
	"google.golang.org/grpc"
	"log"
	"net"
	"golang.org/x/net/context"
	"github.com/beevik/ntp"
	"time"
	"math"
	"github.com/jays-computing/gosync/utils"
)

const (
	port = "127.0.0.1:50051"
	ntpHost = "0.pool.ntp.org"
	//ntpHost = "10.0.75.2"
)

type server struct {
	sessions map[string] *gs.Session
	sessionGuests map[string] map[int32] *Guest
}
type Guest struct {
	stream *gs.GoSync_GetEventsServer
	name string
	avg_rtt int32
}
var online = true

func (s *server) JoinSession(context context.Context, jr *gs.JoinRequest) (*gs.JoinReply, error) {
	println (jr.Name, " joined!!")
	session := ServerSessions.sessions["1"]
	guestId := len(ServerSessions.sessionGuests["1"]) + 1
	ServerSessions.sessionGuests[session.SessionId][int32(guestId)] = &Guest{
		name:jr.Name,
	}
	return &gs.JoinReply{
		Session: session,
		GuestId: int32(guestId),
	}, nil
}

func (s *server) GetEvents(ger *gs.GetEventsRequest, stream gs.GoSync_GetEventsServer) error {
	guest := ServerSessions.sessionGuests[ger.Session.SessionId][ger.GuestId]
	println("Getting events for session Id: " , ger.Session.SessionId, " for ", guest.name, " RTT: ", guest.avg_rtt)
	guest.stream = &stream
	guest.avg_rtt = ger.NtpTimeRtt
	for online {

	}
	return nil
}

func (s *server) PublishEvent(context context.Context, pr *gs.PublishRequest) (*gs.PublishResult , error) {
	// Find max RTT
	var maxRtt int64
	guests := ServerSessions.sessionGuests[pr.Session.SessionId]
	for _, g := range guests {
		maxRtt = int64(math.Max(float64(maxRtt), float64(g.avg_rtt)))
	}
	maxOwt := maxRtt / 2
	ntpNow := ( utils.ToMilli(time.Now().UnixNano()) - StartTimeLocal ) + StartTimeNtp

	var i = 0
	for _, g := range guests {

		if g.stream == nil {
			continue
		}
		// Calculate offset from NtpTime

		message := gs.GSMessage{
			Message: pr.Message.Message,
			Time: int32(ntpNow + maxOwt),
		}
		println("I: " , i)
		// publish
		if err := (*g.stream).Send(&message); err != nil {
			log.Println("Failed to send to guest ")
		}
	}
	return &gs.PublishResult{
	}, nil
}

var ServerSessions server;
var StartTimeLocal int64;
var StartTimeNtp int64;


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	ServerSessions = server {
		sessions:make(map[string]*gs.Session),
		sessionGuests:make(map[string]map[int32]*Guest),
	}

	ntpTime, err := ntp.Query(ntpHost, 4)
	if err != nil {
		log.Fatal("NTP Query failed", err)
		return
	}
	StartTimeNtp = utils.ToMilli(ntpTime.Time.UnixNano())
	StartTimeLocal = utils.ToMilli(time.Now().UnixNano())

	ServerSessions.sessions["1"] = &gs.Session{
		SessionId: "1",
		SessionName: "Sessh",
		NtpHost: ntpHost,
	}
	ServerSessions.sessionGuests["1"] = make(map[int32]*Guest, 0)
	log.Println("Server Started")
	gs.RegisterGoSyncServer(s, &server{})
	s.Serve(lis)
}
