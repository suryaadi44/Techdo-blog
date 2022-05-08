package middleware

import (
	"net/http"

	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				switch response := r.(type) {
				case *globalDTO.BaseResponse:
					globalDTO.NewBaseResponse(response.Code, true, response.Data).SendResponse(&w)
				case error:
					globalDTO.NewBaseResponse(http.StatusInternalServerError, true, response.Error()).SendResponse(&w)
				default:
					globalDTO.NewBaseResponse(http.StatusInternalServerError, true, "runtime error").SendResponse(&w)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
