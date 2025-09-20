package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.temporal.io/sdk/client"

	_ "github.com/0xRichardL/temporal-practice/payment/docs"
	"github.com/0xRichardL/temporal-practice/payment/internal/rest"
	"github.com/0xRichardL/temporal-practice/payment/internal/services"
	"github.com/0xRichardL/temporal-practice/shared/workflows"
)

func main() {
	/* Temporal: */
	// Connect to Temporal server
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	// Register workflows
	workflows.RegisterPaymentWorkFlow(c)

	/* Services: */
	paymentService := services.NewPaymentService(c)

	/* REST server: */
	r := gin.Default()
	paymentController := rest.NewPaymentController(paymentService)
	paymentController.RegisterRoutes(r)
	// Swagger:
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalln("Unable to start Payment server", err)
	}
}
