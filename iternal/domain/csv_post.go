package domain

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type PostCSV struct {
	AdminID   uuid.UUID      `gorm:"column:adminID"`
	CreatedAt time.Time      `gorm:"column:created-at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Size      int64          `gorm:"column:size"`
	FileName  string         `gorm:"column:file-name"`
}
