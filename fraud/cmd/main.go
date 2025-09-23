package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	_ "github.com/0xRichardL/temporal-practice/fraud/docs"
	"github.com/0xRichardL/temporal-practice/fraud/internal/rest"
	"github.com/0xRichardL/temporal-practice/fraud/internal/services"
	"github.com/0xRichardL/temporal-practice/shared/workflows"
)

func main() {
	/// Temporal:
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()
	w := worker.New(c, workflows.FraudCheckWorkflowTaskQueue, worker.Options{})
	workflows.RegisterFraudCheckWorkflow(w)
	/// Services:
	fraudService := services.NewFraudService(c)
	/// Rest controllers:
	r := gin.Default()
	fraudController := rest.NewFraudController(fraudService)
	fraudController.RegisterRoutes(r)
	/// Swagger:
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger/", func(c *gin.Context) { c.Redirect(301, "/swagger/index.html") })
	// Start worker at another routine.
	go func() {
		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start Worker: %s, Error: %v", workflows.FraudCheckWorkflowTaskQueue, err)
		}
	}()
	// Start the rest server.
	log.Println("Starting Fraud server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalln("Unable to start Fraud server", err)
	}
}
