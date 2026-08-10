package entity

import (
	"time"

	"gorm.io/datatypes"
)

type SavePropertyRequest struct {
	Source            string         `json:"source"`
	SourcePropertyID  string         `json:"source_property_id"`
	SourceURL         string         `json:"source_url"`
	Name              *string        `json:"name"`
	Address           *string        `json:"address"`
	RentYen           *int64         `json:"rent_yen"`
	ManagementFeeYen  *int64         `json:"management_fee_yen"`
	DepositYen        *int64         `json:"deposit_yen"`
	KeyMoneyYen       *int64         `json:"key_money_yen"`
	RoomLayout        *string        `json:"room_layout"`
	AreaSquareMeters  *float64       `json:"area_square_meters"`
	BuiltYear         *int           `json:"built_year"`
	Floor             *string        `json:"floor"`
	Direction         *string        `json:"direction"`
	Access            datatypes.JSON `json:"access"`
	Amenities         datatypes.JSON `json:"amenities"`
	RealEstateCompany *string        `json:"real_estate_company"`
	Details           datatypes.JSON `json:"details"`
	CapturedAt        time.Time      `json:"captured_at"`
}

type PropertyResponse struct {
	ID                int64          `json:"id"`
	SearchProfileID   int64          `json:"search_profile_id"`
	Source            string         `json:"source"`
	SourcePropertyID  string         `json:"source_property_id"`
	SourceURL         string         `json:"source_url"`
	Name              *string        `json:"name"`
	Address           *string        `json:"address"`
	RentYen           *int64         `json:"rent_yen"`
	ManagementFeeYen  *int64         `json:"management_fee_yen"`
	DepositYen        *int64         `json:"deposit_yen"`
	KeyMoneyYen       *int64         `json:"key_money_yen"`
	RoomLayout        *string        `json:"room_layout"`
	AreaSquareMeters  *float64       `json:"area_square_meters"`
	BuiltYear         *int           `json:"built_year"`
	Floor             *string        `json:"floor"`
	Direction         *string        `json:"direction"`
	Access            datatypes.JSON `json:"access"`
	Amenities         datatypes.JSON `json:"amenities"`
	RealEstateCompany *string        `json:"real_estate_company"`
	Details           datatypes.JSON `json:"details"`
	CapturedAt        time.Time      `json:"captured_at"`
}
