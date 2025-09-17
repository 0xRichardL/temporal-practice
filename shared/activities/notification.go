package activities

import (
	"context"

	"go.temporal.io/sdk/worker"
)

var NotifyPaymentActivityName = "notification::payment:notify"

type NotifyPaymentActivityParam struct {
	AccountID string
	Amount    int64
}

type NotificationActivities interface {
	Register(w worker.Worker)
	NotifyPayment(ctx context.Context, param NotifyPaymentActivityParam) error
}
