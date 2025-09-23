package services

import (
	"context"
	"fmt"

	"github.com/0xRichardL/temporal-practice/notification/internal/dtos"
)

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendSMS(ctx context.Context, dto dtos.SendSMSRequest) error {
	fmt.Println(dto.Message)
	return nil
}
