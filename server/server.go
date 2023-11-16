package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example/server/handlers"
	"example/server/services"
	"example/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Serve(config types.Config) {
	ctrls := handlers.Setup(services.Setup(config))

	engine := echo.New()
	engine.HideBanner = true
	engine.Use(middleware.Recover(), middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "HTTP Request time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	reservations := engine.Group("/reservations")
	{
		ctrl := ctrls.Reservations
		reservations.GET("", ctrl.GetReservations)
		reservations.GET("/:id", ctrl.GetReservationByID)
		reservations.POST("", ctrl.CreateReservation)
	}

	go func() {
		if err := engine.Start(fmt.Sprintf(":%d", config.Port)); err != nil {
			slog.Info(fmt.Sprintf("Shutting down the server: %v\n", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down the server gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := engine.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("Unable to shutdown server gracefully: %v\n", err))
		return
	}
	slog.Info("Server has shut down gracefully")
}
