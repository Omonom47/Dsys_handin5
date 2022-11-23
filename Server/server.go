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

var portOld = flag.String("server", ":9100", "Tcp server")
var port = flag.Int("port", 9100, "Tcp server")
var timeStamp time.Time
var timeZero time.Time
var timeLimit time.Time

func (s *server) SendBid(ctx context.Context, b *handin.Bid) (*handin.Ack, error) {
	//set time stamp and time-limt
	if timeStamp == timeZero { //makes sure it's only set the first time the sendBid is called
		timeStamp = time.Now()
		timeLimit = timeStamp.Add(time.Duration(1) * time.Minute)
	}

	if _, ok := s.auctionBids[b.Id]; ok {
		//makes sure that the bidAmount can't be set to less than the current max
		var maxValue int32
		for _, currentBid := range s.auctionBids {
			if currentBid >= maxValue {
				maxValue = currentBid
			}
		}
		if maxValue < b.BidAmount {
			ack := handin.Ack{Outcome: "SUCCES"}
			fmt.Printf("ACK: %v", ack.Outcome)
			log.Printf("adjusted bid for user %v with value %v", b.Id, b.BidAmount)
			return &ack, nil
		} else {
			ack := handin.Ack{Outcome: "FAILURE"}
			fmt.Printf("ACK: %s", ack.Outcome)
			log.Printf("failure to adjust bid for user %v with value %v, because of smaller value", b.Id, b.BidAmount)
			return &ack, nil
		}
	} else {
		s.auctionBids[b.Id] = b.BidAmount
		ack := handin.Ack{Outcome: "SUCCES"}
		fmt.Printf("ACK: %v", ack.Outcome)
		log.Printf("added user %v with value %v to server map", b.Id, b.BidAmount)
		return &ack, nil
	}
}

func (s *server) GetResults(ctx context.Context, p *emptypb.Empty) (*handin.Result, error) {
	//do something here to get result

	var currentHighest int32
	if time.Now().Before(timeLimit) {
		for _, highestBid := range s.auctionBids {
			if highestBid > currentHighest {
				currentHighest = highestBid
			}
		}
		log.Printf("Returned Result to client with current max value of %v", currentHighest)
		res := handin.Result{InProcess: true, HighestBid: currentHighest}
		return &res, nil
	} else {
		res := handin.Result{InProcess: false, HighestBid: currentHighest}
		log.Printf("Time limited exceeded")
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

	var network string

	network = fmt.Sprintf("%v", *portOld)
	fmt.Printf("Network: %v", network)
	//hardcoded server test
	for i := 0; i < 3; i++ {
		// port := int32(5000) + int32(i)
		port := int32(*port) + int32(i)
		lis, err := net.Listen("tcp", ("localhost" + fmt.Sprintf(":%v", port)))
		fmt.Print("Connection to Port ", port)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		handin.RegisterAuctionServer(grpcServer, newServer())
		grpcServer.Serve(lis)
	}

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
