package grpc

import (
	"context"
	"github/kijunpos/app"
	"github/kijunpos/app/handler/user"
	"github/kijunpos/app/lib/apm"
	pb "github/kijunpos/gen/proto/user"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(setup *app.SetupData) {
	apmTearDown := apm.InitTracer(context.Background(), apm.Option{
		ServiceName:  setup.ConfigData.Otel.ServiceName,
		CollectorURL: strings.TrimPrefix(setup.ConfigData.Otel.URL, "http://"),
		ApiKey:       setup.ConfigData.Otel.ApiKey,
		Environment:  setup.ConfigData.Otel.Env,
		Insecure:     setup.ConfigData.Otel.Insecure,
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Init user grpc
	pb.RegisterUserServiceServer(s, user.New(setup.ConfigData))

	reflection.Register(s)

	// Channel untuk menerima signal OS
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server dalam goroutine terpisah
	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	// Menunggu signal untuk shutdown
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	// Buat context dengan timeout untuk graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown
	done := make(chan bool)
	go func() {
		s.GracefulStop() // Menunggu semua request selesai
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout: forcing shutdown")
		s.Stop() // Force shutdown jika timeout
	case <-done:
		log.Println("Server shutdown gracefully")
	}

	// Handle APM teardown
	if err := apmTearDown(shutdownCtx); err != nil {
		log.Printf("Failed to teardown APM: %v", err)
	} else {
		log.Println("APM teardown completed successfully")
	}

	// Close database connections if needed
	if setup.DBManager != nil {
		setup.DBManager.CloseConnections()
		log.Println("Database connections closed")
	}

	log.Println("Server shutdown completed")
}
