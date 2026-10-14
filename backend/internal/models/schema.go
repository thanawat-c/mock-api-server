package models

import (
	"time"

	"gorm.io/datatypes"
)

type Endpoint struct {
	ID      uint           `gorm:"primaryKey" json:"id"`
	Path    string         `gorm:"index:idx_path_method,unique;not null" json:"path"`
	Method  string         `gorm:"index:idx_path_method,unique;not null" json:"method"`
	Delay   int            `gorm:"default:0" json:"delay"` // Response delay in milliseconds
	Headers datatypes.JSON `gorm:"type:jsonb" json:"headers"`
	Schema  datatypes.JSON `gorm:"type:jsonb" json:"schema"`
}

type MockRecord struct {
	ID uint `gorm:"primaryKey" json:"id"`
	EndpointID uint `gorm:"index;not null" json:"endpoint_id"`
	Data datatypes.JSON `gorm:"type:jsonb" json:"data"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

