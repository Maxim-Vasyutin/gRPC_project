// Package app is reserved for dependency wiring and application lifecycle.
package app

import (
	"context"
	userv1 "contracts/generated/user/v1"
	"net"
	"user-service/internal/adapters/inbound/grpcadapter"
	"user-service/internal/adapters/outbound/hasher"
	"user-service/internal/adapters/outbound/idgen"
	"user-service/internal/adapters/outbound/postgres"
	"user-service/internal/application"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	//адрес БД и порт
	dsn      = "postgres://user:password@localhost:5434/users?sslmode=disable"
	grpcPort = ":50051"
)

func connectionDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func buildHandler(pool *pgxpool.Pool) *grpcadapter.Server {
	repo := postgres.NewPostgresUserRepository(pool)

	bcryptHasher := hasher.BcryptHasher{}
	idGen := idgen.UUIDGenerator{}

	uc := application.NewRegisterUseCase(repo, bcryptHasher, idGen)

	return grpcadapter.NewServer(uc)
}

func startServer(handler *grpcadapter.Server) error {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	userv1.RegisterUserServiceServer(grpcServer, handler)

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func Run() error {
	ctx := context.Background()

	pool, err := connectionDB(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	handler := buildHandler(pool)

	err = startServer(handler)
	if err != nil {
		return err
	}

	return nil

}
