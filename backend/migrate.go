package main

import (
	"context"
	"log"

	"go-pro-backend/internal/config"
	"go-pro-backend/internal/repository/postgres"
	migrations "go-pro-backend/internal/repository/postgres/migrations"
	"go-pro-backend/pkg/logger"
)

func main() {
	cfg := config.Load()
	db, err := postgres.NewConnection(&postgres.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.Database,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create logger for migration manager
	l := logger.New("info", "text")

	// Create migration manager and register all migrations
	migrationManager := postgres.NewMigrationManager(db, l)
	migrationManager.RegisterMultiple(migrations.GetAllMigrations())

	// Run migrations
	if err := migrationManager.Up(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations completed successfully")
}
