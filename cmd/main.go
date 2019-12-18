package main

import (
	"context"
	"net"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	pb "github.com/willzhao90/hellobackend/out"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"github.com/willzhao90/hellogrpc/pkg/rpc"
	
)

const (
	port        = "0.0.0.0:8030"
	serviceName = "service.hello"
)

type HelloService struct {
	rpc    *rpc.Server
	health *health.Server
	//db     *mongo.Database
}

type Service interface {
	Run(ctx context.Context)
}

func envOrDefaultString(envVar string, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue
	}

	return value
}

func (s *HelloService) Run(ctx context.Context) {
	lis, err := net.Listen("tcp", envOrDefaultString("hello_rpc:server:port", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	pb.RegisterHelloServiceServer(gs, s.rpc)
	grpc_health_v1.RegisterHealthServer(gs, s.health)
	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	// Register reflection service on gRPC server.
	reflection.Register(gs)

	go func() {
		select {
		case <-ctx.Done():
			gs.GracefulStop()
		}
	}()

	log.Infof("Listening at %v...\n", port)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	log.Info("Start:")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	helloServer := &HelloService{
		rpc:    rpc.NewServer(),
		health: health.NewServer(),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		helloServer.Run(ctx)
	}()
	wg.Wait()
}
