package repository

import (
	"product-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductReportRepository interface {
	Create(productReport *model.ProductReport) error
	FindByID(id uuid.UUID) (*model.ProductReport, error)
	FindByProductID(productID uuid.UUID) ([]model.ProductReport, error)
	FindByUserID(userID uuid.UUID) ([]model.ProductReport, error)
	FindAll() ([]model.ProductReport, error)
	Update(productReport *model.ProductReport) error
	Delete(id uuid.UUID) error
}

type productReportRepository struct {
	db *gorm.DB
}

func NewProductReportRepository(db *gorm.DB) ProductReportRepository {
	return &productReportRepository{db: db}
}

func (r *productReportRepository) Create(productReport *model.ProductReport) error {
	return r.db.Create(productReport).Error
}

func (r *productReportRepository) FindByID(id uuid.UUID) (*model.ProductReport, error) {
	var productReport model.ProductReport
	err := r.db.First(&productReport, "id = ?", id).Error
	return &productReport, err
}

func (r *productReportRepository) FindByProductID(productID uuid.UUID) ([]model.ProductReport, error) {
	var productReports []model.ProductReport
	err := r.db.Where("product_id = ?", productID).Find(&productReports).Error
	return productReports, err
}

func (r *productReportRepository) FindByUserID(userID uuid.UUID) ([]model.ProductReport, error) {
	var productReports []model.ProductReport
	err := r.db.Where("user_id = ?", userID).Find(&productReports).Error
	return productReports, err
}

func (r *productReportRepository) FindAll() ([]model.ProductReport, error) {
	var productReports []model.ProductReport
	err := r.db.Find(&productReports).Error
	return productReports, err
}

func (r *productReportRepository) Update(productReport *model.ProductReport) error {
	return r.db.Save(productReport).Error
}

func (r *productReportRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ProductReport{}, "id = ?", id).Error
}
