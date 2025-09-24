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
	router.POST("/payment", c.CreatePaymentHandler)
	router.GET("/payment/:workflowID/status")
}

// CreatePaymentHandler API
// @Summary Create a new payment
// @Description Creates a new payment and starts a Temporal workflow.
// @Tags payment
// @Accept json
// @Produce json
// @Param request body dtos.CreatePaymentRequest true "Payment request"
// @Success 200 {object} dtos.CreatePaymentResponse "Payment created successfully"
// @Failure 400 {object} map[string]string "{"error": "Bad request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Router /payment [post]
func (c *PaymentController) CreatePaymentHandler(ctx *gin.Context) {
	var dto dtos.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.paymentService.CreatePayment(ctx, dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// GetPaymentStatusHandler API
// @Summary Get payment status
// @Description Gets the status of a payment.
// @Tags payment
// @Accept json
// @Produce json
// @Param workflowID path string true "Workflow ID"
// @Success 200 {object} dtos.GetPaymentStatusResponse "Payment status retrieved successfully"
// @Failure 400 {object} map[string]string "{"error": "Bad request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Router /payment/{workflowID}/status [get]
func (c *PaymentController) GetPaymentStatusHandler(ctx *gin.Context) {
	var dto dtos.GetPaymentStatusRequest
	if err := ctx.ShouldBindUri(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res, err := c.paymentService.GetPaymentStatus(ctx, dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
