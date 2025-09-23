package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/0xRichardL/temporal-practice/notification/internal/services"
	"github.com/0xRichardL/temporal-practice/notification/temporal"
	"github.com/0xRichardL/temporal-practice/shared/activities"
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
	w := worker.New(c, activities.NotificationActivityTaskQueue, worker.Options{})
	/// Services:
	notificationService := services.NewNotificationService()
	/// Temporal Activities:
	notificationActivities := temporal.NewNotificationActivities(notificationService)
	notificationActivities.Register(w)
	// Run worker at main routine.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start Worker: %s, Error: %v", activities.NotificationActivityTaskQueue, err)
	}
}
