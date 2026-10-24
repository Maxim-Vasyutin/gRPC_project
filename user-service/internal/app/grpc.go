package app

import (
	"fmt"
	"log"
	"net"

	userv1 "contracts/generated/user/v1"
	"user-service/internal/adapters/inbound/grpcadapter"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// newGRPCServer собирает транспортный слой: интерсепторы и хендлеры.
func newGRPCServer(uc useCases) (*grpc.Server, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("create validator: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcadapter.ValidationInterceptor(v),
		),
	)

	// reflection нужен Postman/grpcurl, чтобы узнать список методов
	reflection.Register(grpcServer)

	userv1.RegisterUserServiceServer(grpcServer, grpcadapter.NewServer(uc.register, uc.login))

	return grpcServer, nil
}

// serve занимает порт и блокируется, пока сервер работает.
func serve(grpcServer *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen %s: %w", port, err)
	}

	log.Printf("gRPC server listening on %s", port)
	return grpcServer.Serve(lis)
}
