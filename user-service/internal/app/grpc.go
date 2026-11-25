package app

import (
	"context"
	"fmt"
	"net"

	userv1 "contracts/generated/user/v1"
	"user-service/internal/adapters/inbound/grpcadapter"
	"user-service/internal/application"

	"buf.build/go/protovalidate"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// newGRPCServer собирает транспортный слой: интерсепторы и хендлеры.
func newGRPCServer(
	register *application.RegisterUseCase,
	login *application.LoginUseCase,
	refresh *application.RefreshUseCase,
	logout *application.LogoutUseCase,
	logger *zap.Logger) (*grpc.Server, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("create validator: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcadapter.RecoveryInterceptor(logger),
			grpcadapter.ValidationInterceptor(v),
		),
	)

	// reflection нужен Postman/grpcurl, чтобы узнать список методов
	reflection.Register(grpcServer)

	userv1.RegisterUserServiceServer(
		grpcServer,
		grpcadapter.NewServer(
			register,
			login,
			refresh,
			logout,
			logger,
		),
	)

	return grpcServer, nil
}

func registerServer(lc fx.Lifecycle, srv *grpc.Server, cfg config, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", cfg.grpcPort)
			if err != nil {
				return fmt.Errorf("listen %s: %w", cfg.grpcPort, err)
			}

			logger.Info("gRPC server listening", zap.String("port", cfg.grpcPort))

			go func() {
				err := srv.Serve(lis)
				if err != nil {
					logger.Error("service stopped with error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})
}
