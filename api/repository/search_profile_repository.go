package repository

import (
	"context"

	"github.com/kikudesuyo/room-finder/api/entity"
	"gorm.io/gorm"
)

func CreateSearchProfile(ctx context.Context, db *gorm.DB, profile *entity.DBTableSearchProfile) error {
	return db.WithContext(ctx).Create(profile).Error
}

func GetSearchProfile(ctx context.Context, db *gorm.DB, id int64) (*entity.DBTableSearchProfile, error) {
	var profile entity.DBTableSearchProfile
	err := db.WithContext(ctx).First(&profile, id).Error
	return &profile, err
}

func ListSearchProfiles(ctx context.Context, db *gorm.DB) ([]entity.DBTableSearchProfile, error) {
	var profiles []entity.DBTableSearchProfile
	err := db.WithContext(ctx).Order("id ASC").Find(&profiles).Error
	return profiles, err
}
