package main

import (
	"context"
	"log"
	"net"

	"github.com/AndreiStefanie/grpc-go/calculator/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalcServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.CalcRequest) (*pb.CalcResponse, error) {
	res := &pb.CalcResponse{
		Result: req.Operands.First + req.Operands.Second,
	}

	return res, nil
}

func main() {
	log.Println("Starting the calculator server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalcServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
