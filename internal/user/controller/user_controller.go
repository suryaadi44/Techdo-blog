package controller

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	middlewarePkg "github.com/suryaadi44/Techdo-blog/internal/middleware"
	postPkg "github.com/suryaadi44/Techdo-blog/internal/post/service"
	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/internal/user/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type UserController struct {
	router         *mux.Router
	userService    service.UserServiceApi
	sessionService service.SessionServiceApi
	postService    postPkg.PostServiceApi
	authMiddleware middlewarePkg.AuthMiddleware
}

func NewUserController(router *mux.Router, userService service.UserServiceApi, sessionService service.SessionServiceApi, postService postPkg.PostServiceApi, authMiddleware middlewarePkg.AuthMiddleware) *UserController {
	return &UserController{
		router:         router,
		userService:    userService,
		sessionService: sessionService,
		postService:    postService,
		authMiddleware: authMiddleware,
	}
}

func (u *UserController) InitializeController() {
	authRouter := u.router.PathPrefix("/").Subrouter()
	authRouter.Use(u.authMiddleware.AuthMiddleware())

	//with middleware
	//API
	authRouter.HandleFunc("/user/detail", u.updateUserDetailHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/user/detail/picture", u.updateUserPictureHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/user/detail", u.getUserDetailHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/user/mini-detail", u.getUserMiniDetailHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/user/delete", u.deleteUserHandler).Methods(http.MethodDelete)

	// Page
	authRouter.HandleFunc("/user/settings", u.settingPageHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/user", u.userDashboardPageHandler).Methods(http.MethodGet)

	//without middleware
	//API

	// Page

}

func (u *UserController) settingPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/user/edit-user-profiles.html"))
	var err error

	token, _ := utils.GetSessionToken(r)
	session, err := u.sessionService.GetSession(r.Context(), token)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	user, err := u.userService.GetUserDetail(r.Context(), session.UID)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	data := map[string]interface{}{
		"User":     user,
		"Username": session.Username,
	}

	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (u *UserController) userDashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/user/dashboard.html"))

	token, _ := utils.GetSessionToken(r)
	session, err := u.sessionService.GetSession(r.Context(), token)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}
	user, err := u.userService.GetUserMiniDetail(r.Context(), session.UID)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	totalComment, err := u.postService.GetUserTotalCommentCount(r.Context(), user.UserID)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	totalPost, err := u.postService.GetUserTotalPostCount(r.Context(), user.UserID)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	recentComments, err := u.postService.GetCommentsByUser(r.Context(), user.UserID, 1, 5)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	recentPost, err := u.postService.GetMiniBlogPostsByUser(r.Context(), user.UserID, 1, 5)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	data := map[string]interface{}{
		"User":           user,
		"CommentCount":   totalComment,
		"PostCount":      totalPost,
		"RecentComments": recentComments,
		"RecentPost":     recentPost,
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (u *UserController) updateUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, err := u.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	payload := dto.UserDetailRequest{}

	if err := decoder.Decode(&payload); err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	payload.UserID = session.UID

	err = u.userService.UpdateUserDetail(r.Context(), payload)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, nil).SendResponse(&w)
}

func (u *UserController) updateUserPictureHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, err := u.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	uploadedFile, handler, err := r.FormFile("profile-pic")
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
	defer uploadedFile.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, uploadedFile); err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	err = u.userService.UpdateUserPicture(r.Context(), buf.Bytes(), handler.Filename, session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, nil).SendResponse(&w)
}

func (u *UserController) getUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, err := u.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	user, err := u.userService.GetUserDetail(r.Context(), session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, user).SendResponse(&w)
}

func (u *UserController) getUserMiniDetailHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, err := u.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	user, err := u.userService.GetUserMiniDetail(r.Context(), session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, user).SendResponse(&w)
}

func (u *UserController) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, err := u.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	err = u.userService.DeleteUser(r.Context(), session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	utils.DeleteSessionCookie(&w)
	globalDTO.NewBaseResponse(http.StatusOK, false, nil).SendResponse(&w)
}
