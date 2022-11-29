package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	handin "handin5.dk/uni/grpc"
)

var (
	port           = flag.Int("port", 9100, "Tcp server")
	timeStamp      time.Time
	timeLimit      time.Time
	currentHighest int32
	maxValue       int32
)

func (s *server) SendBid(ctx context.Context, b *handin.Bid) (*handin.Ack, error) {
	//set time stamp and time-limt
	if timeStamp.IsZero() { //makes sure it's only set the first time the sendBid is called
		timeStamp = time.Now()
		timeLimit = timeStamp.Add(time.Duration(1) * time.Minute)
	}

	if b.BidAmount > maxValue && time.Now().Before(timeLimit) {
		maxValue = b.BidAmount
		if _, ok := s.auctionBids[b.Id]; ok {
			//makes sure that the bidAmount can't be set to less than the current max
			ack := handin.Ack{Outcome: "SUCCES"}
			s.auctionBids[b.Id] = b.BidAmount
			fmt.Printf("server %v: %v for changing bid\n", *port, ack.Outcome)
			log.Printf("server %v: adjusted bid for user %v with amount %v", *port, b.Id, b.BidAmount)
			return &ack, nil
		} else {
			s.auctionBids[b.Id] = b.BidAmount
			ack := handin.Ack{Outcome: "SUCCES"}
			log.Printf("server %v : added user %v with bid %v to server map", *port, b.Id, b.BidAmount)
			return &ack, nil
		}
	} else {
		ack := handin.Ack{Outcome: "FAILURE"}
		if time.Now().After(timeLimit) {
			fmt.Printf("server %v : %v for changing bid\n", *port, ack.Outcome)
			log.Printf("server %v : failure to adjust bid for user %v with amount %v, because the Auction has ended", *port, b.Id, b.BidAmount)
			return &ack, nil
		} else {
			fmt.Printf("server %v : %v for changing bid\n", *port, ack.Outcome)
			log.Printf("server %v : failure to adjust bid for user %v with amount %v, because the bid was smaller than then currently highest bid", *port, b.Id, b.BidAmount)
			return &ack, nil
		}
	}
}

func (s *server) GetResults(ctx context.Context, p *emptypb.Empty) (*handin.Result, error) {
	if time.Now().Before(timeLimit) {
		for _, highestBid := range s.auctionBids {
			if highestBid >= currentHighest {
				currentHighest = highestBid
			}
		}
		log.Printf("server %v : Returned Result to client with current highest bid of %v", *port, currentHighest)
		res := handin.Result{InProcess: true, HighestBid: currentHighest}
		return &res, nil
	} else {
		res := handin.Result{InProcess: false, HighestBid: currentHighest}
		log.Printf("server %v : Time limited exceeded with highest bid being %v", *port, currentHighest)
		return &res, nil
	}
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

	//hardcoded server test
	port := int32(*port)
	lis, err := net.Listen("tcp", ("localhost" + fmt.Sprintf(":%v", port)))
	fmt.Printf("Connection to Port %v\n ", port)
	log.Printf("Server connection to port %v", port)
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
