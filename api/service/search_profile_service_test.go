package service

import (
	"context"
	"testing"

	"github.com/kikudesuyo/room-finder/api/entity"
	"github.com/morikuni/failure"
)

func TestUpdateSearchProfileRejectsChanges(t *testing.T) {
	err := UpdateSearchProfile(context.Background(), 1, entity.UpdateSearchProfileRequest{})
	if err == nil {
		t.Fatal("UpdateSearchProfile() accepted a change")
	}
	message, ok := failure.MessageOf(err)
	if !ok || message != "request validation failed" {
		t.Fatalf("UpdateSearchProfile() error message = %q", message)
	}
}
