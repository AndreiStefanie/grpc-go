package main

import (
	"context"
	"io"
	"log"

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
	doServerStream(c)
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
