package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"content-service/internal/config"
	"content-service/internal/delivery/grpc/handler"
	"content-service/internal/delivery/grpc/interceptor"
	"content-service/internal/delivery/grpc/pb"
	"content-service/internal/infrastructure/database"
	"content-service/internal/infrastructure/repository"
	"content-service/internal/infrastructure/security"
	"content-service/internal/infrastructure/storage"
	"content-service/internal/usecase/content"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("[INFO] Starting Content Management Service...")

	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[FATAL] Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize Database Connection Pool
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)
	if err != nil {
		log.Fatalf("[FATAL] Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("[INFO] Database connection pool established successfully")

	// 3. Run Database Migrations
	if err := database.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		log.Fatalf("[FATAL] Failed to execute database migrations: %v", err)
	}

	// 4. Initialize Clean Architecture Layers
	storageService := storage.NewLocalStorageService(cfg.StorageUploadDir, cfg.StorageBaseURL)
	contentRepo := repository.NewContentRepository(pool)
	contentUsecase := content.NewContentUsecase(contentRepo, storageService)
	contentHandler := handler.NewContentHandler(contentUsecase)

	// 5. Initialize Auth Interceptor & gRPC Server
	jwtValidator := security.NewJWTValidator(cfg.JWTSecret)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(jwtValidator)),
	)
	pb.RegisterContentServiceServer(grpcServer, contentHandler)

	// Register gRPC reflection for development inspection
	reflection.Register(grpcServer)

	// 6. Start TCP Listener for gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("[FATAL] Failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	go func() {
		log.Printf("[INFO] Content Service gRPC server is listening on :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("[FATAL] Failed to serve gRPC: %v", err)
		}
	}()

	// 7. Start HTTP Static File Server
	_ = os.MkdirAll(cfg.StorageUploadDir, 0755)
	httpMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(cfg.StorageUploadDir))
	httpMux.Handle("/storage/uploads/", http.StripPrefix("/storage/uploads/", fileServer))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: httpMux,
	}

	go func() {
		log.Printf("[INFO] Content Service HTTP File Server is listening on :%s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[ERROR] HTTP File Server stopped: %v", err)
		}
	}()

	// 8. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down Content Service gracefully...")
	grpcServer.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = httpServer.Shutdown(shutdownCtx)

	log.Println("[INFO] Content Service stopped successfully")
}
