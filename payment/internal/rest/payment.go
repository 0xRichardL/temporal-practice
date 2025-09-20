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
}

// PaymentHandler API
// @Summary Create a new payment
// @Description Creates a new payment and starts a Temporal workflow.
// @Tags payment
// @Accept json
// @Produce json
// @Param request body dtos.PaymentRequest true "Payment request"
// @Success 200 {object} dtos.PaymentResponse "Payment created successfully"
// @Failure 400 {object} map[string]string "{"error": "Bad request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Router /payment [post]
func (c *PaymentController) PaymentHandler(ctx *gin.Context) {
	var dto dtos.PaymentRequest
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
