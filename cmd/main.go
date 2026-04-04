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
	"github.com/booking-api-service/pkg/logger"
)

func main() {
	log := logger.GetLogger()

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "root:12345@tcp(localhost:3306)/booking_api?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

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
	router.Use(appMiddleware.RequestIDMiddleware(log))

	// Routes
	// Availability endpoints
	router.Post("/coaches/{coach_id}/availability", availabilityHandler.SetAvailability)
	router.Post("/coaches/{coach_id}/exceptions", availabilityHandler.AddException)
	router.Get("/coaches/{coach_id}/slots", availabilityHandler.GetSlots)

	// Booking endpoints
	router.Post("/bookings", bookingHandler.CreateBooking)
	router.Get("/users/{user_id}/bookings", bookingHandler.GetUserBookings)
	router.Put("/bookings/{id}", bookingHandler.ModifyBooking)
	router.Delete("/bookings/{id}", bookingHandler.CancelBooking)

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Info("starting server", slog.String("address", addr))

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")

	// Create context with timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", slog.String("error", err.Error()))
	}

	// Wait for stuck queries to complete (with timeout)
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	// Close database connection
	if err := closeDatabase(dbCtx, db, log); err != nil {
		log.Error("failed to close database", slog.String("error", err.Error()))
	}

	log.Info("server stopped")
}

// closeDatabase closes the database connection with context
func closeDatabase(ctx context.Context, db *sql.DB, log *slog.Logger) error {
	// Wait for connections to finish within timeout
	done := make(chan error, 1)
	go func() {
		done <- db.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Error("error closing database", slog.String("error", err.Error()))
			return err
		}
		log.Info("database closed successfully")
		return nil
	case <-ctx.Done():
		log.Error("database close timeout exceeded")
		return ctx.Err()
	}
}
