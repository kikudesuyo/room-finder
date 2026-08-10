package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kikudesuyo/room-finder/api/xerror"
)

type BaseResponse struct {
	Meta     ResponseMeta   `json:"meta"`
	Err      *ResponseError `json:"error,omitempty"`
	DataType string         `json:"data_type"`
	Data     any            `json:"data,omitempty"`
}

type ResponseMeta struct {
	StatusCode int  `json:"status_code"`
	Success    bool `json:"success"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewObjectResponse(data any) BaseResponse {
	return BaseResponse{Meta: ResponseMeta{StatusCode: http.StatusOK, Success: true}, DataType: "object", Data: data}
}

func NewCreatedResponse(data any) BaseResponse {
	return BaseResponse{Meta: ResponseMeta{StatusCode: http.StatusCreated, Success: true}, DataType: "object", Data: data}
}

func NewErrorResponse(err error) BaseResponse {
	return BaseResponse{
		Meta:     ResponseMeta{StatusCode: xerror.GetHTTPErrorCode(err), Success: false},
		Err:      &ResponseError{Code: xerror.GetStringCode(err), Message: xerror.GetErrorMessage(err)},
		DataType: "error",
	}
}

func (r BaseResponse) Render(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Meta.StatusCode)
	_, _ = fmt.Fprint(w, r.GetBody())
}

func (r BaseResponse) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	r.Render(w)
}

func (r BaseResponse) GetStatusCode() int {
	return r.Meta.StatusCode
}

func (r BaseResponse) GetBody() string {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(r)
	return buf.String()
}
