package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	"kage/backend/internal/api"
	"kage/backend/internal/auth"
	"kage/backend/internal/bidding"
	"kage/backend/internal/trip"
	"kage/backend/internal/ws"
)

// Application wires the HTTP server and supporting services.
type Application struct {
	Engine  *gin.Engine
	Hub     *ws.Hub
	cleanup func(context.Context) error
}

// Build assembles dependencies using the supplied config and logger.
func Build(cfg Config, logger *log.Logger) (*Application, error) {
	if logger == nil {
		logger = log.Default()
	}

	var db *sql.DB
	var err error
	if cfg.DBDSN != "" {
		// 1.- Open the MariaDB connection when a DSN is provided so repositories can persist data.
		db, err = sql.Open("mysql", cfg.DBDSN)
		if err != nil {
			return nil, fmt.Errorf("open database: %w", err)
		}
	}

	hub := ws.NewHub(logger)

	var bidRepo bidding.Repository
	if db != nil {
		bidRepo = bidding.NewSQLRepository(db)
	}
	arbiter := bidding.NewArbiter(bidRepo, bidding.WithTimeout(cfg.EvaluationTimeout), bidding.WithRadius(cfg.RadiusKm))

	var tripRepo trip.EventRepository
	if db != nil {
		tripRepo = trip.NewSQLEventRepository(db)
	}
	tripManager := trip.NewManager(tripRepo, nil)

	validator := auth.NewValidator(cfg.AuthSecret)

	router := gin.New()
	router.Use(gin.Recovery())

	server := api.NewServer(arbiter, tripManager, validator)
	server.RegisterRoutes(router)
	hub.RegisterRoutes(router)

	cleanup := func(ctx context.Context) error {
		// 2.- Stop background workers before closing shared connections.
		hub.Shutdown(ctx)
		if db != nil {
			return db.Close()
		}
		return nil
	}

	return &Application{Engine: router, Hub: hub, cleanup: cleanup}, nil
}

// Cleanup releases resources created by Build.
func (a *Application) Cleanup(ctx context.Context) error {
	if a.cleanup != nil {
		return a.cleanup(ctx)
	}
	return nil
}
