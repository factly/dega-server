package config

import "time"

// Base with id, created_at, updated_at & deleted_at
type Base struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
	CreatedByID uint       `json:"created_by_id"`
	UpdatedByID uint       `json:"updated_by_id"`
}
