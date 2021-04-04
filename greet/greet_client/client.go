package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/AndreiStefanie/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	// doUnary(c)
	// doServerStream(c)
	doClientStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andrei",
			LastName:  "Stefanie"}}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not receive greeting: %v", err)
	}

	log.Println(res.Result)
}

func doServerStream(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andrei",
			LastName:  "Stefanie",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to greet many times: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			log.Println("Stream closed")
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}

		log.Println(msg.Result)
	}
}

func doClientStream(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Could not long greet: %v", err)
	}

	requests := []*greetpb.GreetRequest{
		{
			Greeting: &greetpb.Greeting{FirstName: "Andrei"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Petru"},
		},
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while closing the client stream: %v", err)
	}

	log.Println(res.Result)
}
