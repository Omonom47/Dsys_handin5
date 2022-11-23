package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	handin "handin5.dk/uni/grpc"
)

var port = flag.String("server", ":9100", "Tcp server")

func (s *server) SendBid(ctx context.Context, b *handin.Bid) (*handin.Ack, error) {

	if val, ok := s.auctionBids[b.Id]; ok {
		if val < b.BidAmount {
			ack := handin.Ack{Outcome: "SUCCES"}
			fmt.Printf("ACK: %v", ack)
			return &ack, nil
		} else {
			ack := handin.Ack{Outcome: "FAILURE"}
			fmt.Printf("ACK: %s", ack.Outcome)
			return &ack, nil
		}
	} else {
		s.auctionBids[b.Id] = b.BidAmount
		ack := handin.Ack{Outcome: "SUCCES"}
		fmt.Printf("ACK: %v", ack)
		return &ack, nil
	}
}

func (s *server) GetResults(ctx context.Context, p *emptypb.Empty) (*handin.Result, error) {
	//do something here to get result
	fmt.Printf("Get Result here")
	res := handin.Result{InProcess: true, HighestBid: 3}
	return &res, nil
}

func main() {
	flag.Parse()
	//log to different txt Log file
	LOG_FILE := "./txtLog"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	var network string

	network = fmt.Sprintf("%v", *port)
	lis, err := net.Listen("tcp", ("localhost" + network))
	fmt.Print("Connection to Port ", network)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	handin.RegisterAuctionServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

func newServer() *server {
	s := &server{
		auctionBids: make(map[int32]int32),
	}
	return s
}

type server struct {
	handin.UnimplementedAuctionServer
	auctionBids map[int32]int32
}
