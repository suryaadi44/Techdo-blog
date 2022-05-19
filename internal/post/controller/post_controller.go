package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	middlewarePkg "github.com/suryaadi44/Techdo-blog/internal/middleware"
	"github.com/suryaadi44/Techdo-blog/internal/post/service"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/user/service"
	userServicePkg "github.com/suryaadi44/Techdo-blog/internal/user/service"
)

type PostController struct {
	router         *mux.Router
	postService    service.PostServiceApi
	sessionService authServicePkg.SessionServiceApi
	authMiddleware middlewarePkg.AuthMiddleware
	userService    userServicePkg.UserServiceApi
}

func NewController(router *mux.Router, postService service.PostServiceApi, sessionService authServicePkg.SessionServiceApi, authMiddleware middlewarePkg.AuthMiddleware, userService userServicePkg.UserServiceApi) *PostController {
	return &PostController{
		router:         router,
		postService:    postService,
		sessionService: sessionService,
		authMiddleware: authMiddleware,
		userService:    userService,
	}
}

func (p *PostController) InitializeController() {
	authRouter := p.router.PathPrefix("/").Subrouter()
	authRouter.Use(p.authMiddleware.AuthMiddleware())

	// With middlerware
	// Page
	authRouter.HandleFunc("/post/create", p.createPostPageHandler).Methods(http.MethodGet)

	// API
	authRouter.HandleFunc("/post/create", p.createPostHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/post/delete/{id:[0-9]+}", p.deletePostHandlder).Methods(http.MethodDelete)
	authRouter.HandleFunc("/post/{id:[0-9]+}/comment/add", p.addCommentHandler).Methods(http.MethodPost)

	// Without middleware
	// Page
	p.router.HandleFunc("/", p.postDashboardPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/search", p.searchPostPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/post/latest", p.latestPostPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/post/{id:[0-9]+}", p.viewPostPageHandlder).Methods(http.MethodGet)
	p.router.HandleFunc("/post/category/{category}", p.postInCategoryPageHandler).Methods(http.MethodGet)

	// API
	p.router.HandleFunc("/post/{id:[0-9]+}/comment", p.viewCommentHandler).Methods(http.MethodGet)
}
