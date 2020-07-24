package model

import "time"

// Log is model for log table.
type Log struct {
	ID          int       `gorm:"primary_key;type:serial"`
	UserID      int       `gorm:"int"`
	Action      string    `gorm:"type:varchar"`
	RequestedAt time.Time `gorm:"type:timestamp"`
	CreatedAt   time.Time `gorm:"type:timestamp"`
}
