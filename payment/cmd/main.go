package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"

	"github.com/0xRichardL/temporal-practice/payment/internal/rest"
	"github.com/0xRichardL/temporal-practice/payment/internal/services"
	"github.com/0xRichardL/temporal-practice/payment/internal/temporal/workflows"
)

func main() {
	/* Temporal: */
	// Connect to Temporal server
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
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

	log.Println("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalln("Unable to start server", err)
	}
}
