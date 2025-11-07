package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/zhunismp/intent-products-api/internal/adapters/primary/grpc/product"
	core "github.com/zhunismp/intent-products-api/internal/core/infrastructure/config"
	"go.uber.org/zap"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	productv1 "github.com/zhunismp/intent-proto/product/gen/go/proto/v1"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	cfg      core.AppConfigProvider
	log      *zap.Logger
	grpcApp  *gogrpc.Server
	listener net.Listener
}

func NewGrpcServer(cfg core.AppConfigProvider, zapLogger *zap.Logger) *GrpcServer {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GetGrpcServerPort()))
	if err != nil {
		zapLogger.Fatal("Failed to start grpc listener")
	}

	grpcPanicRecoveryHandler := func(p any) error {
		zapLogger.Error("panic recovered", zap.Any("panic", p))
		return status.Errorf(codes.Internal, "internal server error")
	}

	opts := []gogrpc.ServerOption{
		gogrpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 15 * time.Minute, // disconnect idle clients
			Timeout:           20 * time.Second, // timeout for ping ack
			Time:              2 * time.Minute,  // ping interval
		}),
		gogrpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Minute, // minimum ping interval from client
			PermitWithoutStream: true,
		}),
		gogrpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	}
	app := gogrpc.NewServer(opts...)

	zapLogger.Info("Grpc server core initialized")
	return &GrpcServer{
		cfg:      cfg,
		log:      zapLogger,
		grpcApp:  app,
		listener: lis,
	}
}

func (s *GrpcServer) Start() {
	serverAddr := fmt.Sprintf("%s:%s", s.cfg.GetServerHost(), s.cfg.GetGrpcServerPort())
	s.log.Info("Attempting to start GRPC server...", zap.String("address", serverAddr))

	go func() {
		if err := s.grpcApp.Serve(s.listener); err != nil {
			s.log.Error("Failed to start GRPC server listener", zap.Error(err))
		}
	}()

	s.log.Info("GRPC server listener started.", zap.String("address", serverAddr))
}

func (s *GrpcServer) RegisterServices(productGrpcHandler *product.ProductGrpcHandler) {
	if productGrpcHandler == nil {
		s.log.Fatal("Failed to setup grpc handler")
	}

	// setup
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Registration
	healthpb.RegisterHealthServer(s.grpcApp, healthServer)
	productv1.RegisterProductServiceServer(s.grpcApp, productGrpcHandler)
	reflection.Register(s.grpcApp)
}

func (s *GrpcServer) GracefulShutdown() {
	s.log.Info("Gracefully shutting down GRPC server...")
	s.grpcApp.GracefulStop()

	s.log.Info("GRPC server shutdown gracefully")
}
