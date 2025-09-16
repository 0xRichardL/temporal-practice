package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/0xRichardL/temporal-practice/payment/internal/dtos"
	"github.com/0xRichardL/temporal-practice/payment/internal/services"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

func (c *PaymentController) RegisterRoutes(router *gin.Engine) {
	router.POST("/payment", c.PaymentHandler)
	router.POST("/fraud-check", c.FraudCheckHandler)
}

func (c *PaymentController) PaymentHandler(ctx *gin.Context) {
	var dto dtos.PaymentRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	res, err := c.paymentService.CreatePayment(ctx, dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *PaymentController) FraudCheckHandler(ctx *gin.Context) {
	var dto dtos.FraudCheckRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
