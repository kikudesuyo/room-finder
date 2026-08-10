package repository

import (
	"context"

	"github.com/kikudesuyo/room-finder/api/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpsertRentalOffer(ctx context.Context, db *gorm.DB, offer *entity.DBTableRentalOffer) error {
	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "search_profile_id"},
			{Name: "source"},
			{Name: "source_offer_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"source_url", "name", "address", "rent_yen", "management_fee_yen",
			"deposit_yen", "key_money_yen", "room_layout", "area_square_meters",
			"built_year", "floor", "direction", "access", "amenities",
			"real_estate_company", "details", "captured_at",
		}),
	}).Create(offer).Error
}
