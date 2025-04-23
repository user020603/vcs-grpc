package interceptor

import (
	"log"
	"google.golang.org/grpc"
)

func StreamLoggingInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("Streaming RPC: %s", info.FullMethod)
	err := handler(srv, ss)
	if err != nil {
		log.Printf("Error in Streaming RPC %s: %v", info.FullMethod, err)
	}
	return err
}