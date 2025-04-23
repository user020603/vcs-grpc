package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/user020603/grpc-greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr = "localhost:50051"

func main() {
	tls := true
	opts := []grpc.DialOption{}

	if tls {
		certFile := "../ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("Failed to load TLS credentials")
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	fmt.Println("Starting gRPC client with TLS...")

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	doUnary(c)
	doServerStreaming(c)
	doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c pb.GreetServiceClient) {
	fmt.Println("\nStarting Unary RPC...")
	req := &pb.GreetRequest{
		Greeting: &pb.Greeting{
			FirstName: "Thanh",
			LastName:  "NT",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c pb.GreetServiceClient) {
	fmt.Println("\nStarting Server Streaming RPC...")
	req := &pb.GreetRequest{
		Greeting: &pb.Greeting{
			FirstName: "ThanhNT",
			LastName:  "208",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c pb.GreetServiceClient) {
	fmt.Println("\nStarting Client Streaming RPC...")

	requests := []*pb.GreetRequest{
		{
			Greeting: &pb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Jane",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Alex",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c pb.GreetServiceClient) {
	fmt.Println("\nStarting BiDi Streaming RPC...")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*pb.GreetRequest{
		{
			Greeting: &pb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Jane",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Alex",
			},
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				close(waitc)
				return
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
	}()

	<-waitc
}
