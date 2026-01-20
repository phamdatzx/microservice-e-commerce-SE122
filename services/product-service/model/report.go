package model

import (
	"time"

	"github.com/google/uuid"
)

// ReportReason represents the reason for reporting a product
type ReportReason string

const (
	ReportReasonFakeProduct           ReportReason = "fake_product"
	ReportReasonQualityIssue          ReportReason = "quality_issue"
	ReportReasonMisleadingDescription ReportReason = "misleading_description"
	ReportReasonCounterfeit           ReportReason = "counterfeit"
	ReportReasonDamagedProduct        ReportReason = "damaged_product"
	ReportReasonOther                 ReportReason = "other"
)

// ReportStatus represents the status of a report
type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"
	ReportStatusReviewed ReportStatus = "reviewed"
	ReportStatusResolved ReportStatus = "resolved"
	ReportStatusRejected ReportStatus = "rejected"
)

type Report struct {
	ID        string `bson:"_id" json:"id"`
	ProductID string `bson:"product_id" json:"product_id"`
	VariantID string `bson:"variant_id" json:"variant_id"`

	User        User         `bson:"user" json:"user"`
	Reason      ReportReason `bson:"reason" json:"reason"`
	Description string       `bson:"description" json:"description"`
	Status      ReportStatus `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (r *Report) BeforeCreate() {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	if r.Status == "" {
		r.Status = ReportStatusPending
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
}
