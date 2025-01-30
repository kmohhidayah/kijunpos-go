package db

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type PSQLName string

const (
	KIJUNDB PSQLName = "kijundb"
)

type Config struct {
	Name        PSQLName
	HostURL     string
	MaxIdleConn int
	MaxOpenConn int
}

type Connection struct {
	DB   *sqlx.DB
	Name PSQLName
}

type Manager struct {
	connections map[PSQLName]Connection
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[PSQLName]Connection),
	}
}

func (m *Manager) InitConnections(configs ...Config) error {
	for _, config := range configs {
		db, err := sqlx.Connect("pgx", config.HostURL)
		if err != nil {
			return fmt.Errorf("failed init connection to %s: %w", config.Name, err)
		}

		m.connections[config.Name] = Connection{Name: config.Name, DB: db}
	}
	return nil
}

func (m *Manager) GetConnection(name PSQLName) (*Connection, error) {
	connected, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("no database connection found with name: %s", name)
	}

	if err := connected.DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed ping connection to %s: %w", name, err)
	}

	return &connected, nil
}

func (m *Manager) CloseConnections() {
	for name, conn := range m.connections {
		if err := conn.DB.Close(); err != nil {
			log.Printf("Failed to close database connection (%s): %v", name, err)
		}
	}
}
