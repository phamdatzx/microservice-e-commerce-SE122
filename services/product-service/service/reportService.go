package service

import (
	"net/http"
	"product-service/client"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportService interface {
	CreateReport(report *model.Report) error
	GetReportByID(id string) (*model.Report, error)
	GetAllReports(page, limit int, status *model.ReportStatus, sortOrder int) ([]model.Report, int64, error)
	GetReportsByProductID(productID string, page, limit int) ([]model.Report, int64, error)
	GetReportsByUserID(userID string, page, limit int) ([]model.Report, int64, error)
	UpdateReportStatus(id string, status model.ReportStatus) error
	DeleteReport(id string) error
}

type reportService struct {
	repo        repository.ReportRepository
	orderClient *client.OrderServiceClient
	userClient  *client.UserServiceClient
}

func NewReportService(repo repository.ReportRepository, orderClient *client.OrderServiceClient, userClient *client.UserServiceClient) ReportService {
	return &reportService{
		repo:        repo,
		orderClient: orderClient,
		userClient:  userClient,
	}
}

func (s *reportService) CreateReport(report *model.Report) error {
	// Fetch user info from user-service
	userInfo, err := s.userClient.GetUserByID(report.User.ID)
	if err != nil {
		return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to fetch user information", err)
	}

	// Populate user data
	report.User.Name = userInfo.Name
	report.User.Email = userInfo.Email
	report.User.Image = userInfo.Image
	report.User.Phone = userInfo.Phone

	// Validate if the user has purchased the variant
	hasPurchased, err := s.orderClient.VerifyVariantPurchase(report.User.ID, report.ProductID, report.VariantID)
	if err != nil {
		return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to verify purchase", err)
	}

	if !hasPurchased {
		return appError.NewAppError(http.StatusForbidden, "You can only report products you have purchased")
	}

	// Check if user has already reported this product
	_, err = s.repo.FindByProductIDAndUserID(report.ProductID, report.User.ID)
	if err == nil {
		// Found document → user already reported
		return appError.NewAppError(http.StatusForbidden, "You have already reported this product")
	}

	if err != mongo.ErrNoDocuments {
		return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to check existing report", err)
	}

	// Validate report reason
	if !isValidReportReason(report.Reason) {
		return appError.NewAppError(http.StatusBadRequest, "Invalid report reason")
	}

	// Generate ID and timestamps
	report.ID = uuid.New().String()
	report.BeforeCreate()

	return s.repo.Create(report)
}

func (s *reportService) GetReportByID(id string) (*model.Report, error) {
	return s.repo.FindByID(id)
}

func (s *reportService) GetAllReports(page, limit int, status *model.ReportStatus, sortOrder int) ([]model.Report, int64, error) {
	skip := (page - 1) * limit
	return s.repo.FindAll(skip, limit, status, sortOrder)
}

func (s *reportService) GetReportsByProductID(productID string, page, limit int) ([]model.Report, int64, error) {
	skip := (page - 1) * limit
	return s.repo.FindByProductID(productID, skip, limit)
}

func (s *reportService) GetReportsByUserID(userID string, page, limit int) ([]model.Report, int64, error) {
	skip := (page - 1) * limit
	return s.repo.FindByUserID(userID, skip, limit)
}

func (s *reportService) UpdateReportStatus(id string, status model.ReportStatus) error {
	// Validate status
	if !isValidReportStatus(status) {
		return appError.NewAppError(http.StatusBadRequest, "Invalid report status")
	}

	report, err := s.repo.FindByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return appError.NewAppError(http.StatusNotFound, "Report not found")
		}
		return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to find report", err)
	}

	report.Status = status
	return s.repo.Update(report)
}

func (s *reportService) DeleteReport(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return appError.NewAppError(http.StatusNotFound, "Report not found")
		}
		return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to find report", err)
	}

	return s.repo.Delete(id)
}

func isValidReportReason(reason model.ReportReason) bool {
	switch reason {
	case model.ReportReasonFakeProduct,
		model.ReportReasonQualityIssue,
		model.ReportReasonMisleadingDescription,
		model.ReportReasonCounterfeit,
		model.ReportReasonDamagedProduct,
		model.ReportReasonOther:
		return true
	}
	return false
}

func isValidReportStatus(status model.ReportStatus) bool {
	switch status {
	case model.ReportStatusPending,
		model.ReportStatusReviewed,
		model.ReportStatusResolved,
		model.ReportStatusRejected:
		return true
	}
	return false
}
