package handler

import (
	"net/http"

	"github.com/kikudesuyo/room-finder/api/entity"
)

type ProcessFunc func(r *http.Request, requestData map[string]interface{}) (http.Handler, error)

func HandleRequestAndResponse(r *http.Request, w http.ResponseWriter, processFn ProcessFunc) {
	respData, err := processFn(r, map[string]interface{}{})
	if err != nil {
		entity.NewErrorResponse(err).ServeHTTP(w, r)
		return
	}
	respData.ServeHTTP(w, r)
}
