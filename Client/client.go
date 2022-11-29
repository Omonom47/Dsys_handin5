package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	handin "handin5.dk/uni/grpc"
)

// flag for id
var flagId = flag.Int("id", 0, "Id for User")
var results = make([]handin.Result, 0, 0)
var clientConns = make([]handin.AuctionClient, 0, 0)

func main() {
	flag.Parse()
	//makes log file
	LOG_FILE := "./txtLog"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())
	var client handin.AuctionClient

	//dials all servers
	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Fatalf("Client %v : Failed to connect: %v", flagId, err)
		}

		client = handin.NewAuctionClient(conn)
		clientConns = append(clientConns, client)

		fmt.Printf("Client %v : Connected to server on port %v\n", *flagId, port)
		defer conn.Close()
	}
	log.Printf("Client %v connected to port %v, %v, %v", *flagId, 5000, 5001, 5002)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := strings.ToLower(scanner.Text())
		if strings.Contains(input, "bid") {
			fmt.Println("Please write the specified amount you wish to bid!")
			scanner.Scan()

			bid, _ := strconv.Atoi(scanner.Text())
			sendBid(ctx, client, int32(bid))
			time.Sleep(20)
		}
		if strings.Contains(input, "result") {
			getResult(ctx, client)
		}
	}
}

func sendBid(ctx context.Context, client handin.AuctionClient, bidAmount int32) {
	bid := handin.Bid{
		BidAmount: bidAmount,
		Id:        int32(*flagId),
	}
	for _, replica := range clientConns {
		go SendBidConcurrently(ctx, replica, &bid)
	}
	fmt.Printf("Client %v Sent Bid\n", *flagId)
}

// makes it run concurrently so bid can be send to all servers at same time
func SendBidConcurrently(ctx context.Context, client handin.AuctionClient, bid *handin.Bid) {
	ack, err := client.SendBid(ctx, bid)
	if err != nil {
		log.Printf("Client %v : Cannot send bid: error: %v", *flagId, err)
	}
	fmt.Println(ack)
}

func getResult(ctx context.Context, client handin.AuctionClient) {
	for _, replica := range clientConns {
		//Code for GetResult in here to call all clients
		result, err := replica.GetResults(ctx, new(emptypb.Empty))
		if err != nil {
			log.Printf("Client %v : Unable to get results: error: %v", *flagId, err)
		} else {
			fmt.Println(result)
			if result.InProcess {
				log.Printf("Client %v : Auction still ongoing with currently highest bid: %v", *flagId, result.HighestBid)
			} else {
				log.Printf("Client %v : Time limit exceeded with final highest bid: %v", *flagId, result.HighestBid)
			}
		}
	}
}
