package models

import "time"

type Campaign struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Remarks     string     `json:"remarks"`
	Status      string     `json:"status"`
	StartTime   *time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime     *time.Time `json:"end_time" gorm:"column:end_time"`
	CreatedBy   int        `json:"created_by" gorm:"column:created_by"`
	UpdatedBy   int        `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt   time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt   time.Time  `json:"-" gorm:"column:updated_at"`
}
