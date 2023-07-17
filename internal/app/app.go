package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/DesmondTo/minichallenge/internal/database"
	"github.com/DesmondTo/minichallenge/internal/transport/handler/flight"
)

func Run() {
	client := database.Connect()
	router := gin.Default()

	// Inject the database connection client into the handlers
	flightHandler := flight.NewHandler(client)
	router.GET("/flight", flightHandler.GetCheapest)
	err := router.Run(":8080")

	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
}
