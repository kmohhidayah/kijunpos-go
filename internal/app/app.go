package app

import (
	"github/kijunpos/config"
	"github/kijunpos/config/db"
	"github/kijunpos/internal/delivery/grpc"
	"github/kijunpos/internal/repository"
	userUseCase "github/kijunpos/internal/usecase/user"
	"log"
)

// Application represents the application with all its dependencies
type Application struct {
	Config      *config.Config
	DBManager   *db.Manager
	GRPCHandler grpc.UserHandler
}

// NewApplication creates and initializes a new application
func NewApplication() *Application {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("error when loading config data: %v", err)
	}
	configData := config.GetConfig()

	// Initialize database connections
	dbManager := db.NewManager()
	if err := dbManager.InitConnections(configData.Databases...); err != nil {
		log.Fatalf("error when initializing database: %v", err)
	}

	// Get database connection
	kijunConn, err := dbManager.GetConnection(db.KIJUNDB)
	if err != nil {
		log.Fatalf("error when getting kijundb connection: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(kijunConn)

	// Initialize use cases
	userUC := userUseCase.NewUserUseCase(userRepo)

	// Initialize gRPC handlers
	userHandler := grpc.NewUserHandler(userUC)

	return &Application{
		Config:      configData,
		DBManager:   dbManager,
		GRPCHandler: userHandler,
	}
}

// Start starts the application
func (app *Application) Start() {
	// Start the gRPC server
	grpc.StartGRPCServer(app.Config, app.GRPCHandler)
}
