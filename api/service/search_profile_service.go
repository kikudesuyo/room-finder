package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/kikudesuyo/room-finder/api/entity"
	"github.com/kikudesuyo/room-finder/api/library"
	"github.com/kikudesuyo/room-finder/api/repository"
	"github.com/kikudesuyo/room-finder/api/xerror"
)

func CreateSearchProfile(ctx context.Context, req entity.CreateSearchProfileRequest) (*entity.SearchProfileResponse, error) {
	if strings.TrimSpace(req.InitialPrompt) == "" {
		return nil, xerror.ClientValidationErr(nil, map[string]string{"field": "initial_prompt"})
	}
	if err := validateJSONMap(req.RequiredConditions); err != nil {
		return nil, xerror.ClientValidationErr(err, map[string]string{"field": "required_conditions"})
	}

	profile := &entity.DBTableSearchProfile{
		InitialPrompt:      req.InitialPrompt,
		RequiredConditions: req.RequiredConditions,
	}
	if err := repository.CreateSearchProfile(ctx, library.GetDB(ctx), profile); err != nil {
		return nil, xerror.UnknownDBErr(err)
	}

	return &entity.SearchProfileResponse{
		ID:                 profile.ID,
		InitialPrompt:      profile.InitialPrompt,
		RequiredConditions: profile.RequiredConditions,
		CreatedAt:          profile.CreatedAt,
	}, nil
}

func validateJSONMap(value []byte) error {
	if len(value) == 0 || !json.Valid(value) {
		return errors.New("required_conditions must be valid JSON")
	}
	var object map[string]any
	if err := json.Unmarshal(value, &object); err != nil || object == nil {
		if err == nil {
			return errors.New("required_conditions must be a JSON object")
		}
		return err
	}
	return nil
}
