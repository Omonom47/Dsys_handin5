package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	handin "handin5.dk/uni/grpc"
)

var id int32

func main() {
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
	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		client = handin.NewAuctionClient(conn)
		defer conn.Close()
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		input := scanner.Text()
		if strings.Contains(input, "Bid") {
			bid, _ := strconv.Atoi(scanner.Text())
			go sendBid(ctx, client, int32(bid))
		}
		if strings.Contains(input, "Result") {
			go getResult(ctx, client)
		}

	}
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
		log.Printf("Unable to get results: error: %v", err)
	}
	log.Printf("Auction still ongoing: %v, currently highest bid: %v", stream.InProcess, stream.HighestBid)

}
