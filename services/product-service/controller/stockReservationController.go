package controller

import (
	"net/http"
	"product-service/dto"
	"product-service/service"

	"github.com/gin-gonic/gin"
)

type StockReservationController struct {
	stockReservationService service.StockReservationService
}

func NewStockReservationController(stockReservationService service.StockReservationService) *StockReservationController {
	return &StockReservationController{
		stockReservationService: stockReservationService,
	}
}

// ReserveStock handles POST /api/v1/stock/reserve
func (ctrl *StockReservationController) ReserveStock(c *gin.Context) {
	var req dto.ReserveStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.stockReservationService.ReserveStock(req.OrderID, req.Items); err != nil {
		// Determine status code based on error message
		statusCode := http.StatusInternalServerError
		if err.Error() == "variant not found" || err.Error() == "no reservations found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "insufficient stock" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock reserved successfully"})
}

// ReleaseStock handles POST /api/v1/stock/release
func (ctrl *StockReservationController) ReleaseStock(c *gin.Context) {
	var req dto.ReleaseStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.stockReservationService.ReleaseStock(req.OrderID); err != nil {
		// Determine status code based on error message
		statusCode := http.StatusInternalServerError
		if err.Error() == "no reservations found" || err.Error() == "no reserved stock found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock released successfully"})
}
