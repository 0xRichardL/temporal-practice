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

// FraudCheckHandler API
// @Summary Perform a fraud check
// @Description Sends a signal to perform a fraud check for a given account.
// @Tags fraud
// @Accept json
// @Produce json
// @Param request body dtos.FraudCheckRequest true "Fraud check request"
// @Success 200 {object} map[string]string "{"message": "Fraud check signal sent"}"
// @Failure 400 {object} map[string]string "{"error": "Bad request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Router /fraud-check [post]
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
