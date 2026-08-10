package entity

import (
	"time"

	"gorm.io/datatypes"
)

type CreateSearchProfileRequest struct {
	InitialPrompt      string         `json:"initial_prompt"`
	RequiredConditions datatypes.JSON `json:"required_conditions"`
}

type UpdateSearchProfileRequest struct {
	InitialPrompt      string         `json:"initial_prompt"`
	RequiredConditions datatypes.JSON `json:"required_conditions"`
}

type SearchProfileResponse struct {
	ID                 int64          `json:"id"`
	InitialPrompt      string         `json:"initial_prompt"`
	RequiredConditions datatypes.JSON `json:"required_conditions"`
	CreatedAt          time.Time      `json:"created_at"`
}
