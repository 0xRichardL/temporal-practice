package main

import (
	"context"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/0xRichardL/temporal-practice/account/internal/models"
	"github.com/0xRichardL/temporal-practice/account/internal/services"
	"github.com/0xRichardL/temporal-practice/account/internal/temporal"
)

func main() {
	/* DBs: */
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URI")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Account{})
	/* Services */
	accountService := services.NewAccountService(db)

	/* Seeds */
	// Initialize the test DB with 4 accounts.
	err = accountService.Seeds(context.Background())
	if err != nil {
		panic("failed to seed database.")
	}
	/* Temporal: */
	// Connect to Temporal server
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	w := worker.New(c, "account-tasks", worker.Options{})
	accountActivities := temporal.NewAccountActivities(accountService)
	accountActivities.Register(w)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
