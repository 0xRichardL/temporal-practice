package temporal

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"

	"github.com/0xRichardL/temporal-practice/notification/internal/dtos"
	"github.com/0xRichardL/temporal-practice/notification/internal/services"
	"github.com/0xRichardL/temporal-practice/shared/activities"
)

type NotificationActivities struct {
	notificationService *services.NotificationService
}

func NewNotificationActivities(notificationService *services.NotificationService) activities.NotificationActivities {
	return &NotificationActivities{
		notificationService: notificationService,
	}
}

func (n *NotificationActivities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(n.NotifyPayment, activity.RegisterOptions{Name: activities.NotifyPaymentActivityName})
}

func (n *NotificationActivities) NotifyPayment(ctx context.Context, param activities.NotifyPaymentActivityParam) error {
	return n.notificationService.SendSMS(ctx, dtos.SendSMSRequest{
		AccountID: param.AccountID,
		Message:   fmt.Sprintf("Payment for Account: %s of %d", param.AccountID, param.Amount),
	})
}
