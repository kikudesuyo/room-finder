package entity

import (
	"time"

	"gorm.io/datatypes"
)

type DBTableProperty struct {
	ID                int64          `gorm:"column:id;primaryKey"`
	SearchProfileID   int64          `gorm:"column:search_profile_id"`
	Source            string         `gorm:"column:source"`
	SourcePropertyID  string         `gorm:"column:source_property_id"`
	SourceURL         string         `gorm:"column:source_url"`
	Name              *string        `gorm:"column:name"`
	Address           *string        `gorm:"column:address"`
	RentYen           *int64         `gorm:"column:rent_yen"`
	ManagementFeeYen  *int64         `gorm:"column:management_fee_yen"`
	DepositYen        *int64         `gorm:"column:deposit_yen"`
	KeyMoneyYen       *int64         `gorm:"column:key_money_yen"`
	RoomLayout        *string        `gorm:"column:room_layout"`
	AreaSquareMeters  *float64       `gorm:"column:area_square_meters"`
	BuiltYear         *int           `gorm:"column:built_year"`
	Floor             *string        `gorm:"column:floor"`
	Direction         *string        `gorm:"column:direction"`
	Access            datatypes.JSON `gorm:"column:access"`
	Amenities         datatypes.JSON `gorm:"column:amenities"`
	RealEstateCompany *string        `gorm:"column:real_estate_company"`
	Details           datatypes.JSON `gorm:"column:details"`
	CapturedAt        time.Time      `gorm:"column:captured_at"`
}

func (DBTableProperty) TableName() string {
	return "properties"
}
