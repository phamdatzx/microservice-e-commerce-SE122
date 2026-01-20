package controller

import (
	"net/http"
	"product-service/error"
	"product-service/model"
	"product-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	service service.ReportService
}

func NewReportController(service service.ReportService) *ReportController {
	return &ReportController{service: service}
}

type CreateReportRequest struct {
	ProductID   string             `json:"product_id" binding:"required"`
	VariantID   string             `json:"variant_id" binding:"required"`
	Reason      model.ReportReason `json:"reason" binding:"required"`
	Description string             `json:"description"`
}

type UpdateReportStatusRequest struct {
	Status model.ReportStatus `json:"status" binding:"required"`
}

func (c *ReportController) CreateReport(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(error.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	var req CreateReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	report := &model.Report{
		ProductID:   req.ProductID,
		VariantID:   req.VariantID,
		Reason:      req.Reason,
		Description: req.Description,
		User: model.User{
			ID: userID,
		},
	}

	if err := c.service.CreateReport(report); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, report)
}

func (c *ReportController) GetReportByID(ctx *gin.Context) {
	id := ctx.Param("id")

	report, err := c.service.GetReportByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, report)
}

func (c *ReportController) GetAllReports(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse status filter (optional)
	var status *model.ReportStatus
	if statusStr := ctx.Query("status"); statusStr != "" {
		s := model.ReportStatus(statusStr)
		status = &s
	}

	// Parse sort order (default: descending -1, ascending: 1)
	sortOrder := -1 // default descending (newest first)
	if sortStr := ctx.Query("sort"); sortStr == "asc" {
		sortOrder = 1
	}

	reports, total, err := c.service.GetAllReports(page, limit, status, sortOrder)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (c *ReportController) GetReportsByProductID(ctx *gin.Context) {
	productID := ctx.Param("productId")

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	reports, total, err := c.service.GetReportsByProductID(productID, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (c *ReportController) GetMyReports(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(error.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	reports, total, err := c.service.GetReportsByUserID(userID, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (c *ReportController) UpdateReportStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var req UpdateReportStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	if err := c.service.UpdateReportStatus(id, req.Status); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report status updated successfully"})
}

func (c *ReportController) DeleteReport(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteReport(id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}
