package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.temporal.io/sdk/client"

	_ "github.com/0xRichardL/temporal-practice/fraud/docs"
	"github.com/0xRichardL/temporal-practice/fraud/internal/rest"
	"github.com/0xRichardL/temporal-practice/fraud/internal/services"
)

func main() {
	// Create the Temporal client
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	fraudService := services.NewFraudService(c)

	fraudController := rest.NewFraudController(fraudService)

	router := gin.Default()
	fraudController.RegisterRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		log.Fatalln("Unable to start Fraud server", err)
	}
}
