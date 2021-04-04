package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/AndreiStefanie/grpc-go/calculator/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	// maximum(c)
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

	numbers := []int32{}

	for _, num := range numbers {
		stream.Send(&pb.AverageRequest{Number: num})
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			log.Println(status.Message())
			return
		} else {
			log.Fatalf("Unexpected error ocurred: %v", err)
		}
	}

	log.Println(res.GetResult())
}

func maximum(c pb.CalcServiceClient) {
	stream, err := c.Maximum(context.Background())
	if err != nil {
		log.Fatalf("Failed to open the stream: %v", err)
	}

	numbers := []int32{1, 5, 3, 6, 2, 20}

	waitc := make(chan struct{})

	go func() {
		for _, num := range numbers {
			err := stream.Send(&pb.MaxRequest{Number: num})
			if err != nil {
				log.Fatalf(" Failed to send the message: %v", err)
			}
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		defer close(waitc)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive the result: %v", err)
			}

			log.Printf("Current maximum: %v\n", res.GetMax())
		}
	}()

	<-waitc
}
