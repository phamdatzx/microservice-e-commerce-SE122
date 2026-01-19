package model

import (
	"time"

	"github.com/google/uuid"
)

type SearchHistory struct {
	ID         string    `bson:"_id" json:"id"`
	UserID     string    `bson:"user_id" json:"user_id"`
	Query      string    `bson:"query" json:"query"`
	SearchedAt time.Time `bson:"searched_at" json:"searched_at"`
}

func (sh *SearchHistory) BeforeCreate() {
	if sh.ID == "" {
		sh.ID = uuid.New().String()
	}
	if sh.SearchedAt.IsZero() {
		sh.SearchedAt = time.Now()
	}
}
