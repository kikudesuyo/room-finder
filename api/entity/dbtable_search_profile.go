package entity

import (
	"time"

	"gorm.io/datatypes"
)

type DBTableSearchProfile struct {
	ID                 int64          `gorm:"column:id;primaryKey"`
	InitialPrompt      string         `gorm:"column:initial_prompt"`
	RequiredConditions datatypes.JSON `gorm:"column:required_conditions"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
}

func (DBTableSearchProfile) TableName() string {
	return "search_profiles"
}
