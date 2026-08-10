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

func SaveRentalOffer(ctx context.Context, profileID int64, req entity.SaveRentalOfferRequest) (*entity.RentalOfferResponse, error) {
	if profileID <= 0 {
		return nil, xerror.ClientValidationErr(errors.New("profile id must be positive"))
	}
	if strings.TrimSpace(req.Source) == "" || strings.TrimSpace(req.SourceOfferID) == "" {
		return nil, xerror.ClientValidationErr(errors.New("source and source_offer_id are required"))
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

	offer := &entity.DBTableRentalOffer{
		SearchProfileID: profileID, Source: req.Source, SourceOfferID: req.SourceOfferID,
		SourceURL: req.SourceURL, Name: req.Name, Address: req.Address, RentYen: req.RentYen,
		ManagementFeeYen: req.ManagementFeeYen, DepositYen: req.DepositYen, KeyMoneyYen: req.KeyMoneyYen,
		RoomLayout: req.RoomLayout, AreaSquareMeters: req.AreaSquareMeters, BuiltYear: req.BuiltYear,
		Floor: req.Floor, Direction: req.Direction, Access: req.Access, Amenities: req.Amenities,
		RealEstateCompany: req.RealEstateCompany, Details: req.Details, CapturedAt: req.CapturedAt,
	}
	if err := repository.UpsertRentalOffer(ctx, db, offer); err != nil {
		return nil, xerror.UnknownDBErr(err)
	}

	return &entity.RentalOfferResponse{
		ID: offer.ID, SearchProfileID: offer.SearchProfileID, Source: offer.Source,
		SourceOfferID: offer.SourceOfferID, SourceURL: offer.SourceURL, Name: offer.Name,
		Address: offer.Address, RentYen: offer.RentYen, ManagementFeeYen: offer.ManagementFeeYen,
		DepositYen: offer.DepositYen, KeyMoneyYen: offer.KeyMoneyYen, RoomLayout: offer.RoomLayout,
		AreaSquareMeters: offer.AreaSquareMeters, BuiltYear: offer.BuiltYear, Floor: offer.Floor,
		Direction: offer.Direction, Access: offer.Access, Amenities: offer.Amenities,
		RealEstateCompany: offer.RealEstateCompany, Details: offer.Details, CapturedAt: offer.CapturedAt,
	}, nil
}
