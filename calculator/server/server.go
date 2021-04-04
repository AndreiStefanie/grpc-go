package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/AndreiStefanie/grpc-go/calculator/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *server) Decompose(req *pb.DecompositionRequest, stream pb.CalcService_DecomposeServer) error {
	number := req.GetDecomposition().GetNumber()

	var factor int32 = 2
	for number > 1 {
		if number%factor == 0 {
			stream.Send(&pb.DecompositionResponse{Factor: factor})
			number /= factor
		} else {
			factor++
		}
	}

	return nil
}

func (s *server) Average(stream pb.CalcService_AverageServer) error {
	var sum int32
	var count int

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			result := 0.0
			if count == 0 {
				return status.Error(codes.InvalidArgument, "At least one number expected")
			}
			result = (float64(sum)) / (float64(count))
			stream.SendAndClose(&pb.AverageResponse{Result: result})
			break
		}
		if err != nil {
			log.Fatalf("Failed while receiving stream request: %v", err)
		}

		sum += req.GetNumber()
		count++
	}

	return nil
}

func (s *server) Maximum(stream pb.CalcService_MaximumServer) error {
	var currentMax int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if req.GetNumber() > currentMax {
			currentMax = req.GetNumber()
			stream.Send(&pb.MaxResponse{Max: currentMax})
		}
	}

	return nil
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
