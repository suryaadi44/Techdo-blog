package middleware

import (
	"net/http"
	"strings"
	"text/template"

	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				switch response := r.(type) {
				case *globalDTO.BaseResponse:
					if strings.Contains(response.Data.(string), "no session") {
						utils.DeleteSessionCookie(&w)
					}
					ErrorPage(&w, globalDTO.NewBaseResponse(response.Code, true, response.Data))
				case error:
					// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					ErrorPage(&w, globalDTO.NewBaseResponse(http.StatusInternalServerError, true, response.Error()))
				default:
					ErrorPage(&w, globalDTO.NewBaseResponse(http.StatusInternalServerError, true, "runtime error"))
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ErrorPage(w *http.ResponseWriter, data *globalDTO.BaseResponse) {
	tmpl := template.Must(template.ParseFiles("web/template/error-route/error-route.html"))

	err := tmpl.Execute(*w, *data)

	if err != nil {
		panic(err)
	}
}
