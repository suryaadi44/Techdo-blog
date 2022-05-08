package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	sessionServicePkg "github.com/suryaadi44/Techdo-blog/internal/auth/service"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type AuthMiddleware struct {
	session sessionServicePkg.SessionServiceApi
}

func (m *AuthMiddleware) AuthMiddleware() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, isLoggedIn := utils.GetSessionToken(r)
			sessionDetail, err := m.session.GetSession(r.Context(), token)

			if err != nil || sessionDetail.ExpireAt.Before(time.Now()) || !isLoggedIn {
				//http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func NewAuthMiddleware(session sessionServicePkg.SessionServiceApi) AuthMiddleware {
	return AuthMiddleware{
		session: session,
	}
}
