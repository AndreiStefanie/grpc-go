package main

import (
	"context"
	"io"
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
	// sum(c)
	decompose(c)
}

func sum(c pb.CalcServiceClient) {
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

func decompose(c pb.CalcServiceClient) {
	req := &pb.DecompositionRequest{
		Decomposition: &pb.Decomposition{
			Number: 120,
		},
	}

	resStream, err := c.Decompose(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to decompose: %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive message from stream: %v", err)
		}

		log.Println(res.GetFactor())
	}
}
