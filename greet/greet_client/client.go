package main

import (
	"context"
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
	doUnary(c)
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
