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
	"github.com/0xRichardL/temporal-practice/shared/activities"
)

func main() {
	/// DBs:
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URI")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Account{})
	/// Temporal:
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	w := worker.New(c, activities.AccountActivityTaskQueue, worker.Options{})
	/// Services:
	accountService := services.NewAccountService(db)
	/// Temporal activities
	accountActivities := temporal.NewAccountActivities(accountService)
	accountActivities.Register(w)
	/// Seeds:
	// Initialize the test DB with 4 accounts.
	err = accountService.Seeds(context.Background())
	if err != nil {
		panic("failed to seed database.")
	}
	// Start the worker at main routine.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start Worker: %s, Error: %v", activities.AccountActivityTaskQueue, err)
	}
}
