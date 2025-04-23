package interceptor

import (
	"context"
	"log"
	"google.golang.org/grpc"
)

func UnaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("Unary RPC: %s, Request: %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error in Unary RPC %s: %v", info.FullMethod, err)
	}
	log.Printf("Unary RPC Response: %v", resp)
	return resp, err
}