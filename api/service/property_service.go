package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/kikudesuyo/room-finder/api/entity"
	"github.com/kikudesuyo/room-finder/api/library"
	"github.com/kikudesuyo/room-finder/api/repository"
	"github.com/kikudesuyo/room-finder/api/xerror"
	"gorm.io/gorm"
)

func SaveProperty(ctx context.Context, profileID int64, req entity.SavePropertyRequest) (*entity.PropertyResponse, error) {
	if profileID <= 0 {
		return nil, xerror.ClientValidationErr(errors.New("profile id must be positive"))
	}
	if strings.TrimSpace(req.Source) == "" || strings.TrimSpace(req.SourcePropertyID) == "" {
		return nil, xerror.ClientValidationErr(errors.New("source and source_property_id are required"))
	}
	parsedURL, err := url.Parse(req.SourceURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return nil, xerror.ClientValidationErr(errors.New("source_url must be an absolute http or https URL"))
	}
	if req.CapturedAt.IsZero() {
		return nil, xerror.ClientValidationErr(errors.New("captured_at is required"))
	}
	if req.Access == nil || req.Amenities == nil || req.Details == nil {
		return nil, xerror.ClientValidationErr(errors.New("access, amenities, and details are required"))
	}

	db := library.GetDB(ctx)
	if _, err := repository.GetSearchProfile(ctx, db, profileID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.ClientResourceNotFoundErr()
		}
		return nil, xerror.UnknownDBErr(err)
	}

	property := &entity.DBTableProperty{
		SearchProfileID: profileID, Source: req.Source, SourcePropertyID: req.SourcePropertyID,
		SourceURL: req.SourceURL, Name: req.Name, Address: req.Address, RentYen: req.RentYen,
		ManagementFeeYen: req.ManagementFeeYen, DepositYen: req.DepositYen, KeyMoneyYen: req.KeyMoneyYen,
		RoomLayout: req.RoomLayout, AreaSquareMeters: req.AreaSquareMeters, BuiltYear: req.BuiltYear,
		Floor: req.Floor, Direction: req.Direction, Access: req.Access, Amenities: req.Amenities,
		RealEstateCompany: req.RealEstateCompany, Details: req.Details, CapturedAt: req.CapturedAt,
	}
	if err := repository.UpsertProperty(ctx, db, property); err != nil {
		return nil, xerror.UnknownDBErr(err)
	}

	return &entity.PropertyResponse{
		ID: property.ID, SearchProfileID: property.SearchProfileID, Source: property.Source,
		SourcePropertyID: property.SourcePropertyID, SourceURL: property.SourceURL, Name: property.Name,
		Address: property.Address, RentYen: property.RentYen, ManagementFeeYen: property.ManagementFeeYen,
		DepositYen: property.DepositYen, KeyMoneyYen: property.KeyMoneyYen, RoomLayout: property.RoomLayout,
		AreaSquareMeters: property.AreaSquareMeters, BuiltYear: property.BuiltYear, Floor: property.Floor,
		Direction: property.Direction, Access: property.Access, Amenities: property.Amenities,
		RealEstateCompany: property.RealEstateCompany, Details: property.Details, CapturedAt: property.CapturedAt,
	}, nil
}
