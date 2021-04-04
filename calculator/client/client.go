package main

import (
	"context"
	"io"
	"log"
	"time"

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
	// decompose(c)
	average(c)
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

func average(c pb.CalcServiceClient) {
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Failed to start streaming data: %v", err)
	}

	numbers := []int32{1, 2, 3, 4}

	for _, num := range numbers {
		stream.Send(&pb.AverageRequest{Number: num})
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive result: %v", err)
	}

	log.Println(res.GetResult())
}
