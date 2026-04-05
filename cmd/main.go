package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"

	"github.com/booking-api-service/internal/handler"
	appMiddleware "github.com/booking-api-service/internal/middleware"
	"github.com/booking-api-service/internal/repository"
	"github.com/booking-api-service/internal/service"
	"github.com/booking-api-service/pkg/config"
	"github.com/booking-api-service/pkg/logger"
)

func main() {
	log := logger.GetLogger()

	// Load configuration from .env file
	cfg := config.Load()

	// Build database connection string from config
	dsn := cfg.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(1 * time.Minute)

	// Create context with timeout for ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		log.Error("failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("database connected successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	coachRepo := repository.NewCoachRepository(db)
	availabilityRepo := repository.NewAvailabilityRepository(db)
	availabilityExceptionRepo := repository.NewAvailabilityExceptionRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// Initialize services
	availabilityService := service.NewAvailabilityService(availabilityRepo, availabilityExceptionRepo, bookingRepo, coachRepo)
	bookingService := service.NewBookingService(bookingRepo, userRepo, coachRepo, db)

	// Initialize handlers
	availabilityHandler := handler.NewAvailabilityHandler(availabilityService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Router setup
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(appMiddleware.RequestTimeoutMiddleware)

	// Routes
	// Availability endpoints
	router.Post("/api/v1/coaches/{coach_id}/availability/weekly", availabilityHandler.SetWeeklyAvailability)
	router.Post("/api/v1/coaches/{coach_id}/exceptions", availabilityHandler.AddException)
	router.Get("/api/v1/coaches/{coach_id}/slots", availabilityHandler.GetSlots)

	// Booking endpoints
	router.Post("/api/v1/bookings", bookingHandler.CreateBooking)
	router.Get("/api/v1/users/{user_id}/bookings", bookingHandler.GetUserBookings)
	router.Put("/api/v1/bookings/{id}", bookingHandler.ModifyBooking)
	router.Delete("/api/v1/bookings/{id}", bookingHandler.CancelBooking)

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Info("starting server", slog.String("address", server.Addr))

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for Signal for temination
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully...")

	// Create a single context for the entire shutdown process
	// This gives the whole process 30s to complete
	shutdown_ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown Server first (stop accepting new traffic)
	if err := server.Shutdown(shutdown_ctx); err != nil {
		log.Error("server shutdown failed", slog.Any("error", err))
	}

	// Shutdown DB (Close after server is done with queries)
	if err := db.Close(); err != nil {
		log.Error("database close failed", slog.Any("error", err))
	} else {
		log.Info("database closed successfully")
	}

	log.Info("server stopped")
}
