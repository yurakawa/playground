package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pancake/maker/gen/api"
	"pancake/maker/handler"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Logger を追加
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	grpc_zap.ReplaceGrpcLogger(zapLogger)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(zapLogger),
				grpc_auth.UnaryServerInterceptor(auth),
			),
		),
	)
	api.RegisterPancakeBakerServiceServer(server, handler.NewBakerHandler())
	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()

}

func auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	if token != "hi/mi/tsu" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid bearer token")
	}

	return context.WithValue(ctx, "UserName", "God"), nil

}
