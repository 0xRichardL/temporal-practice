package activities

var NotificationPaymentActivityName = "notification::payment"

type NotificationPaymentActivityParam struct {
	AccountID string
	Amount    int64
}
