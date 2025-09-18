package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/0xRichardL/temporal-practice/fraud/internal/dtos"
	"github.com/0xRichardL/temporal-practice/fraud/internal/services"
)

type FraudController struct {
	fraudService *services.FraudService
}

func NewFraudController(fraudService *services.FraudService) *FraudController {
	return &FraudController{
		fraudService: fraudService,
	}
}

func (c *FraudController) RegisterRoutes(router *gin.Engine) {
	router.POST("/fraud-check", c.FraudCheckHandler)
}

func (c *FraudController) FraudCheckHandler(ctx *gin.Context) {
	var dto dtos.FraudCheckRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.fraudService.Check(ctx, dto); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Fraud check signal sent"})
}
