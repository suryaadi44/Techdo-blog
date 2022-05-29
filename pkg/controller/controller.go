package controller

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	middlewarePkg "github.com/suryaadi44/Techdo-blog/internal/middleware"
	postControllerPkg "github.com/suryaadi44/Techdo-blog/internal/post/controller"
	postServicePkg "github.com/suryaadi44/Techdo-blog/internal/post/service"
	userControllerPkg "github.com/suryaadi44/Techdo-blog/internal/user/controller"
	userServicePkg "github.com/suryaadi44/Techdo-blog/internal/user/service"
	globalMiddlewarePkg "github.com/suryaadi44/Techdo-blog/pkg/middleware"
)

func InitializeController(router *mux.Router, db *sql.DB) {
	router.Use(globalMiddlewarePkg.ErrorHandler)

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static/"))))

	SessionService := userServicePkg.NewSessionAuthService(db)
	AuthMiddleware := middlewarePkg.NewAuthMiddleware(SessionService)

	AuthService := userServicePkg.NewUserAuthService(db, SessionService)
	PostService := postServicePkg.NewPostService(db)
	UserService := userServicePkg.NewUserService(db)

	AuthController := userControllerPkg.NewController(router, AuthService, SessionService)
	PostController := postControllerPkg.NewController(router, PostService, SessionService, AuthMiddleware, UserService)
	UserController := userControllerPkg.NewUserController(router, UserService, SessionService, PostService, AuthMiddleware)

	AuthController.InitializeController()
	PostController.InitializeController()
	UserController.InitializeController()
}
