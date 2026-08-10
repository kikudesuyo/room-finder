package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kikudesuyo/room-finder/api/entity"
	"github.com/kikudesuyo/room-finder/api/service"
	"github.com/kikudesuyo/room-finder/api/xerror"
)

func HandleCreateSearchProfile(r *http.Request, _ map[string]interface{}) (http.Handler, error) {
	var req entity.CreateSearchProfileRequest
	if err := decodeJSON(r, &req); err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	profile, err := service.CreateSearchProfile(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return entity.NewCreatedResponse(profile), nil
}

func HandleUpdateSearchProfile(r *http.Request, _ map[string]interface{}) (http.Handler, error) {
	profileID, err := parseID(r)
	if err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	var req entity.UpdateSearchProfileRequest
	if err := decodeJSON(r, &req); err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	if err := service.UpdateSearchProfile(r.Context(), profileID, req); err != nil {
		return nil, err
	}
	return nil, errors.New("search profile update unexpectedly succeeded")
}

func HandleSaveRentalOffer(r *http.Request, _ map[string]interface{}) (http.Handler, error) {
	profileID, err := parseID(r)
	if err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	var req entity.SaveRentalOfferRequest
	if err := decodeJSON(r, &req); err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	offer, err := service.SaveRentalOffer(r.Context(), profileID, req)
	if err != nil {
		return nil, err
	}
	return entity.NewCreatedResponse(offer), nil
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func parseID(r *http.Request) (int64, error) {
	return parseInt64(chi.URLParam(r, "id"))
}
