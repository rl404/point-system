package model

import "time"

// UserPoint is model for user_point table.
type UserPoint struct {
	UserID    int       `gorm:"primary_key;auto_increment:false"`
	Point     int       `gorm:"type:int"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
}
