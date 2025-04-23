package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/user020603/grpc-greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedGreetServiceServer
}

// Unary RPC
func (*server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hola " + firstName
	res := &pb.GreetResponse{
		Result: result,
	}
	return res, nil
}

// Server Streaming RPC
func (*server) GreetManyTimes(req *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := fmt.Sprintf("Hola %s, number %d", firstName, i)
		res := &pb.GreetResponse{
			Result: result,
		}

		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

// Client Streaming RPC
func (*server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.GreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

// BiDi Streaming RPC
func (*server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"

		sendErr := stream.Send(&pb.GreetResponse{
			Result: result,
		})

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func main() {
	fmt.Println("Starting gRPC server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &server{})

	reflection.Register(s)

	fmt.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
