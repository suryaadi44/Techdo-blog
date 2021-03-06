package controller

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	internalMiddlewarePkg "github.com/suryaadi44/Techdo-blog/internal/middleware"

	postPkg "github.com/suryaadi44/Techdo-blog/internal/post/service"
	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/internal/user/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	middlewarePkg "github.com/suryaadi44/Techdo-blog/pkg/middleware"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type UserController struct {
	router         *mux.Router
	userService    service.UserServiceApi
	sessionService service.SessionServiceApi
	postService    postPkg.PostServiceApi
	authMiddleware internalMiddlewarePkg.AuthMiddleware
}

func NewUserController(router *mux.Router, userService service.UserServiceApi, sessionService service.SessionServiceApi, postService postPkg.PostServiceApi, authMiddleware internalMiddlewarePkg.AuthMiddleware) *UserController {
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

	//with auth middleware
	//API
	secureApiRouter := authRouter.PathPrefix("/").Subrouter()
	secureApiRouter.Use(middlewarePkg.ApiMiddleware())

	secureApiRouter.HandleFunc("/user/detail", u.updateUserDetailHandler).Methods(http.MethodPost)
	secureApiRouter.HandleFunc("/user/detail/picture", u.updateUserPictureHandler).Methods(http.MethodPost)
	secureApiRouter.HandleFunc("/user/detail", u.getUserDetailHandler).Methods(http.MethodGet)
	secureApiRouter.HandleFunc("/user/mini-detail", u.getUserMiniDetailHandler).Methods(http.MethodGet)
	secureApiRouter.HandleFunc("/user/delete", u.deleteUserHandler).Methods(http.MethodDelete)

	// Page
	authRouter.HandleFunc("/user/settings", u.settingPageHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/user", u.userDashboardPageHandler).Methods(http.MethodGet)

	//without auth middleware
	//API

	// Page
	u.router.HandleFunc("/user/{id:[0-9]+}", u.userProfilePageHandler).Methods(http.MethodGet)

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

	postStats, err := u.postService.GetUserPostStatisticOfEachCategory(r.Context(), user.UserID)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	recentComments, err := u.postService.GetCommentsByUser(r.Context(), user.UserID, 5)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	recentPost, err := u.postService.GetMiniBlogPostsByUser(r.Context(), user.UserID, 5)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	data := map[string]interface{}{
		"User":           user,
		"CommentCount":   totalComment,
		"PostCount":      totalPost,
		"PostStats":      postStats,
		"RecentComments": recentComments,
		"RecentPost":     recentPost,
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (u *UserController) userProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/user/user_profile.html"))

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	user, err := u.userService.GetUserDetail(r.Context(), id)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	totalComment, err := u.postService.GetUserTotalCommentCount(r.Context(), id)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	totalPost, err := u.postService.GetUserTotalPostCount(r.Context(), id)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	postStats, err := u.postService.GetUserPostStatisticOfEachCategory(r.Context(), id)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	recentComments, err := u.postService.GetCommentsByUser(r.Context(), id, 5)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	recentPost, err := u.postService.GetMiniBlogPostsByUser(r.Context(), id, 5)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn":       isLoggedIn,
		"User":           user,
		"CommentCount":   totalComment,
		"PostCount":      totalPost,
		"PostStats":      postStats,
		"RecentComments": recentComments,
		"RecentPost":     recentPost,
	}

	if isLoggedIn {
		session, err := u.sessionService.GetSession(r.Context(), token)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}
		currentUser, err := u.userService.GetUserMiniDetail(r.Context(), session.UID)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}

		if err == nil {
			data["CurrentUser"] = currentUser
		}
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
