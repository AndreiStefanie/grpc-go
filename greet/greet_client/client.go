package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/AndreiStefanie/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c, 1*time.Second)
	doUnary(c, 2*time.Second)
	// doServerStream(c)
	// doClientStream(c)
	// doBiDi(c)
}

func doUnary(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andrei",
			LastName:  "Stefanie"}}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.Greet(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("Deadline exceeded")
			} else {
				log.Println(statusErr.Message())
			}
			return
		} else {
			log.Fatalf("Could not receive greeting: %v", err)
		}
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

func doBiDi(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating the stream: %v", err)
	}

	requests := []*greetpb.GreetRequest{
		{
			Greeting: &greetpb.Greeting{FirstName: "Andrei"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Petru"},
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			stream.Send(req)
			time.Sleep(time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Failed to close the stream: %v", err)
		}
	}()

	go func() {
		defer close(waitc)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
			}
			log.Println(res.GetResult())
		}
	}()

	<-waitc
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("ssl/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
