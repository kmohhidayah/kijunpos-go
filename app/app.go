package app

import (
	"github/kijunpos/config"
	"github/kijunpos/config/db"
	"log"
)

type SetupData struct {
	ConfigData *config.Config
	DBManager  *db.Manager
	// Handler
}

func Init() *SetupData {
	// config init
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("error when load config data, err: %v", err)
	}

	configData := config.GetConfig()

	// DB init
	dbManager := db.NewManager()
	if err := dbManager.InitConnections(configData.Databases...); err != nil {
		log.Fatalf("error when init database, err: %v", err)
	}

	kijunConn, err := dbManager.GetConnection(db.KIJUNDB)
	if err != nil {
		log.Fatalf("error when get kijundb connections, err: %v", err)
	}

	_ = kijunConn
	// grpcServer := grpc.NewServer()

	return &SetupData{
		ConfigData: configData,
		DBManager:  dbManager,
	}
}
