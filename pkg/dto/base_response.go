package dto

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Code  int         `json:"code"`
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
}

func (baseResponse *BaseResponse) SendResponse(w *http.ResponseWriter) error {
	(*w).WriteHeader(baseResponse.Code)
	return json.NewEncoder(*w).Encode(baseResponse)
}

func NewBaseResponse(Code int, Error bool, Data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:  Code,
		Error: Error,
		Data:  Data,
	}
}
