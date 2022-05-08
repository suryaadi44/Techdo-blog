package controller

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	authControllerPkg "github.com/suryaadi44/Techdo-blog/internal/auth/controller"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/auth/service"
	postControllerPkg "github.com/suryaadi44/Techdo-blog/internal/post/controller"
	postServicePkg "github.com/suryaadi44/Techdo-blog/internal/post/service"
	userServicePkg "github.com/suryaadi44/Techdo-blog/internal/user/service"
)

func InitializeController(router *mux.Router, db *sql.DB) {
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static/"))))

	SessionService := authServicePkg.NewSessionAuthService(db)
	AuthService := authServicePkg.NewUserAuthService(db, SessionService)
	AuthController := authControllerPkg.NewController(router, AuthService, SessionService)
	AuthController.InitializeController()

	UserService := userServicePkg.NewUserService(db)

	PostService := postServicePkg.NewPostService(db)
	PostController := postControllerPkg.NewController(router, PostService, SessionService, UserService)
	PostController.InitializeController()
}
