package server

import (
	"context"
	"fmt"
	"log"
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
	engine.Use(middleware.Recover(), middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
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
			log.Printf("Shutting down the server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down the server gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := engine.Shutdown(ctx); err != nil {
		log.Printf("Unable to shutdown server gracefully: %v\n", err)
		return
	}
	log.Println("Server has shut down gracefully")
}
