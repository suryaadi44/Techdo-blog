package controller

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	authControllerPkg "github.com/suryaadi44/Techdo-blog/internal/auth/controller"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/auth/service"
)

func InitializeController(router *mux.Router, db *sql.DB) {
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static/"))))

	AuthService := authServicePkg.NewUserAuthService(db)
	AuthController := authControllerPkg.NewController(router, AuthService)
	AuthController.InitializeController()
}
