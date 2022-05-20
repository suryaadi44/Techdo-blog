package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	middlewarePkg "github.com/suryaadi44/Techdo-blog/internal/middleware"
	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/internal/user/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type UserhController struct {
	router         *mux.Router
	userService    service.UserServiceApi
	sessionService service.SessionServiceApi
	authMiddleware middlewarePkg.AuthMiddleware
}

func NewUserController(router *mux.Router, userService service.UserServiceApi, sessionService service.SessionServiceApi, authMiddleware middlewarePkg.AuthMiddleware) *UserhController {
	return &UserhController{
		router:         router,
		userService:    userService,
		sessionService: sessionService,
		authMiddleware: authMiddleware,
	}
}

func (u *UserhController) InitializeController() {
	authRouter := u.router.PathPrefix("/").Subrouter()
	authRouter.Use(u.authMiddleware.AuthMiddleware())

	//with middleware
	//API
	authRouter.HandleFunc("/user/detail", u.updateUserDetailHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/user/detail", u.getUserDetailHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/user/mini-detail", u.getUserMiniDetailHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/user/delete", u.deleteUserHandler).Methods(http.MethodDelete)

	// Page
	authRouter.HandleFunc("/user/settings", u.settingPageHandler).Methods(http.MethodGet)

	//without middleware
	//API

	// Page

}

func (u *UserhController) settingPageHandler(w http.ResponseWriter, r *http.Request) {
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

func (u *UserhController) updateUserDetailHandler(w http.ResponseWriter, r *http.Request) {
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

func (u *UserhController) getUserDetailHandler(w http.ResponseWriter, r *http.Request) {
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

func (u *UserhController) getUserMiniDetailHandler(w http.ResponseWriter, r *http.Request) {
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

func (u *UserhController) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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
