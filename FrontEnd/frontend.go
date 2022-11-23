package main

/*
import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	handin "handin5.dk/uni/grpc"
)

func main() {

	LOG_FILE := "./txtLog"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	ownPort := 5003

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fe := &frontend{
		clients: make(map[int32]handin.AuctionClient),
		ctx:     ctx,
	}

	// Create listener tcp on port ownPort
	go func() {
		for {
			list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
			if err != nil {
				log.Fatalf("Failed to listen on port: %v", err)
			}
			grpcServer := grpc.NewServer()
			handin.RegisterAuctionServer(grpcServer, fe)

			go func() {
				if err := grpcServer.Serve(list); err != nil {
					log.Fatalf("failed to server %v", err)
				}
			}()
		}
	}()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())
	var client handin.AuctionClient
	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		client = handin.NewAuctionClient(conn)
		//' client.Connect(conn)
		defer conn.Close()
	}
}

func (s *frontend) SendBid(ctx context.Context, b *handin.Bid) (*handin.Ack, error) {

	if _, ok := s.clients[b.Id]; ok {
		ack := handin.Ack{Outcome: "SUCCES"}
		fmt.Printf("ACK: %v", ack)
		return &ack, nil
	} else {
		s.clients[b.Id] = b.Id
		ack := handin.Ack{Outcome: "SUCCES"}
		fmt.Printf("ACK: %v", ack)
		return &ack, nil
	}

	sendBid(ctx)

}

func (s *frontend) GetResults(ctx context.Context, p *emptypb.Empty) (*handin.Result, error) {
	//do something here to get result
	fmt.Printf("Get Result here")
	res := handin.Result{InProcess: true, HighestBid: 3}
	return &res, nil
}

type frontend struct {
	handin.UnimplementedAuctionServer
	clients map[int32]int32
	ctx     context.Context
}

func sendBid(ctx context.Context, client handin.AuctionClient, bidAmount int32) {
	msg := handin.Bid{
		BidAmount: bidAmount,
		Id:        id,
	}
	_, err := client.SendBid(ctx, &msg)
	if err != nil {
		log.Printf("Cannot send bid: error: %v", err)
	}
}

func getResult(ctx context.Context, client handin.AuctionClient) {

	stream, err := client.GetResults(ctx, nil)
	if err != nil {
		log.Printf("Cannot send bid: error: %v", err)
	}

	log.Printf("result %v", stream.HighestBid)

}

*/
