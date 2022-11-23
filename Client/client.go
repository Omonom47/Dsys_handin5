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

var responses = make([]handin.Ack, 0, 0)
var results = make([]handin.Result, 0, 0)
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

	clientConns := make([]grpc.ClientConn, 0, 0)

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		client = handin.NewAuctionClient(conn)
		clientConns = append(clientConns, *conn)
		//' client.Connect(conn)
		defer conn.Close()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		input := scanner.Text()
		if strings.Contains(input, "Bid") {
			bid, _ := strconv.Atoi(scanner.Text())
			responses = make([]handin.Ack, 0, 0)
			for i := 0; i < 3; i++ {
				go sendBid(ctx, client, int32(bid), clientConns[i])
			}
			if responses[0].Outcome != responses[1].Outcome && responses[0].Outcome != responses[2].Outcome {
				if responses[1].Outcome != responses[2].Outcome {
					log.Printf(responses[0].Outcome)
				} else {
					log.Printf(responses[1].Outcome)
				}
			} else {
				log.Printf(responses[0].Outcome)
			}
		}
		if strings.Contains(input, "Result") {
			responses = make([]handin.Ack, 0, 0)
			for i := 0; i < 3; i++ {
				go getResult(ctx, client, clientConns[i])
			}
			if results[0].String() != results[1].String() && results[0].String() != results[2].String() {
				if results[1].String() != results[2].String() {
					log.Printf(results[0].String())
				} else {
					log.Printf(results[1].String())
				}
			} else {
				log.Printf(results[0].String())
			}
		}

	}
}

func sendBid(ctx context.Context, client handin.AuctionClient, bidAmount int32, con grpc.ClientConn) {
	msg := handin.Bid{
		BidAmount: bidAmount,
		Id:        id,
	}
	ack, err := client.SendBid(ctx, &msg)
	if err != nil {
		log.Printf("Cannot send bid: error: %v", err)
	} else {
		responses = append(responses, *ack)
	}

}

func getResult(ctx context.Context, client handin.AuctionClient, con grpc.ClientConn) {

	result, err := client.GetResults(ctx, nil)
	if err != nil {
		log.Printf("Unable to get results: error: %v", err)
	} else {
		results = append(results, *result)
	}
	log.Printf("Auction still ongoing: %v, currently highest bid: %v", result.InProcess, result.HighestBid)

}
