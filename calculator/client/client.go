package main

import (
	"context"
	"log"

	"github.com/AndreiStefanie/grpc-go/calculator/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to the server: %v", err)
	}
	defer conn.Close()

	c := pb.NewCalcServiceClient(conn)

	req := &pb.CalcRequest{
		Operands: &pb.Operands{
			First:  10,
			Second: 3,
		},
	}

	res, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call the add function: %v", err)
	}

	log.Println(res.Result)
}
