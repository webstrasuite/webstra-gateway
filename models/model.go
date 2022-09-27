package models

import (
	"database/sql"
	"time"
)

type Model struct {
	ID        uint `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}
