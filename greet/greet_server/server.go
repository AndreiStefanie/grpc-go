package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/AndreiStefanie/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()

	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetResponse{
			Result: "Hello " + firstName + " for the " + strconv.Itoa(i) + " time",
		}
		stream.Send(res)
		time.Sleep(2 * time.Second)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&greetpb.GreetResponse{
				Result: result,
			})
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		result += "Hello " + req.GetGreeting().GetFirstName() + "! "
	}

	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		result := "Hello " + res.GetGreeting().GetFirstName() + "!"
		err = stream.Send(&greetpb.GreetResponse{Result: result})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("Starting server :D")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
