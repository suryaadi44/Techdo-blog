package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	internalMiddlewarePkg "github.com/suryaadi44/Techdo-blog/internal/middleware"
	"github.com/suryaadi44/Techdo-blog/internal/post/service"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/user/service"
	userServicePkg "github.com/suryaadi44/Techdo-blog/internal/user/service"
	middlewarePkg "github.com/suryaadi44/Techdo-blog/pkg/middleware"
)

type PostController struct {
	router         *mux.Router
	postService    service.PostServiceApi
	sessionService authServicePkg.SessionServiceApi
	authMiddleware internalMiddlewarePkg.AuthMiddleware
	userService    userServicePkg.UserServiceApi
}

func NewController(router *mux.Router, postService service.PostServiceApi, sessionService authServicePkg.SessionServiceApi, authMiddleware internalMiddlewarePkg.AuthMiddleware, userService userServicePkg.UserServiceApi) *PostController {
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

	// With auth middlerware
	// Page
	authRouter.HandleFunc("/post/create", p.createPostPageHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/post/{id:[0-9]+}/edit", p.editPostPageHandlder).Methods(http.MethodGet)

	// API
	secureApiRouter := authRouter.PathPrefix("/").Subrouter()
	secureApiRouter.Use(middlewarePkg.ApiMiddleware())

	secureApiRouter.HandleFunc("/user/post", p.userPostHandler).Methods(http.MethodGet)
	secureApiRouter.HandleFunc("/user/comment", p.userCommentHandler).Methods(http.MethodGet)
	secureApiRouter.HandleFunc("/post/create", p.createPostHandler).Methods(http.MethodPost)
	secureApiRouter.HandleFunc("/post/pick", p.editorPickHandler).Methods(http.MethodPost)
	secureApiRouter.HandleFunc("/post/{id:[0-9]+}/delete", p.deletePostHandler).Methods(http.MethodDelete)
	secureApiRouter.HandleFunc("/post/{id:[0-9]+}/edit", p.editPostHandler).Methods(http.MethodPost)
	secureApiRouter.HandleFunc("/post/{id:[0-9]+}/comment/add", p.addCommentHandler).Methods(http.MethodPost)
	secureApiRouter.HandleFunc("/post/comment/delete", p.deleteCommentHandler).Methods(http.MethodDelete)

	// Without auth middleware
	// Page
	p.router.HandleFunc("/", p.postDashboardPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/search", p.searchPostPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/post/latest", p.latestPostPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/post/{id:[0-9]+}", p.viewPostPageHandlder).Methods(http.MethodGet)
	p.router.HandleFunc("/post/category/{category}", p.postInCategoryPageHandler).Methods(http.MethodGet)

	// API
	apiRouter := authRouter.PathPrefix("/").Subrouter()
	apiRouter.Use(middlewarePkg.ApiMiddleware())

	apiRouter.HandleFunc("/post/{id:[0-9]+}/raw", p.viewRawPostHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc("/post/{id:[0-9]+}/comment", p.viewCommentHandler).Methods(http.MethodGet)
}
