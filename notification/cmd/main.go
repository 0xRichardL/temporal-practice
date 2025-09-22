package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/0xRichardL/temporal-practice/notification/internal/services"
	"github.com/0xRichardL/temporal-practice/notification/temporal"
)

func main() {
	/* Services: */
	notificationService := services.NewNotificationService()

	/* Temporal: */
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()
	w := worker.New(c, "notification-tasks", worker.Options{})
	notificationActivities := temporal.NewNotificationActivities(notificationService)
	notificationActivities.Register(w)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
