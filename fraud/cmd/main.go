package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"

	"github.com/0xRichardL/temporal-practice/fraud/internal/rest"
	"github.com/0xRichardL/temporal-practice/fraud/internal/services"
)

func main() {
	// Create the Temporal client
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	fraudService := services.NewFraudService(c)

	fraudController := rest.NewFraudController(fraudService)

	router := gin.Default()
	fraudController.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalln("Unable to start Fraud server", err)
	}
}
