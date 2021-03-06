package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/suryaadi44/Techdo-blog/internal/user/dto"
	"github.com/suryaadi44/Techdo-blog/internal/user/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"

	middlewarePkg "github.com/suryaadi44/Techdo-blog/pkg/middleware"
)

type UserAuthController struct {
	router         *mux.Router
	userService    service.UserServiceApi
	sessionService service.SessionServiceApi
}

func NewController(router *mux.Router, userAuthService service.UserServiceApi, sessionService service.SessionServiceApi) *UserAuthController {
	return &UserAuthController{
		router:         router,
		userService:    userAuthService,
		sessionService: sessionService,
	}
}

func (u *UserAuthController) InitializeController() {
	//API
	apiRouter := u.router.PathPrefix("/").Subrouter()
	apiRouter.Use(middlewarePkg.ApiMiddleware())

	apiRouter.HandleFunc("/login", u.loginHandler).Methods(http.MethodPost)
	apiRouter.HandleFunc("/signup", u.signUpHandler).Methods(http.MethodPost)

	// Page
	u.router.HandleFunc("/login", u.loginPageHandler).Methods(http.MethodGet)
	u.router.HandleFunc("/signup", u.signUpPageHandler).Methods(http.MethodGet)
	u.router.HandleFunc("/logout", u.logOutHandler).Methods(http.MethodGet)
}

func (u *UserAuthController) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := dto.AuthRequest{}

	if err := decoder.Decode(&payload); err != nil {
		log.Println("[Decode] Error decoding JSON")
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	session, err := u.userService.AuthenticateUser(r.Context(), payload)
	if err != nil {
		log.Println("[Auth] Login failed :", payload.Username)
		globalDTO.NewBaseResponse(http.StatusUnauthorized, true, "Inccorect Username or Password").SendResponse(&w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session.Token,
		Expires: session.ExpireAt,
	})

	globalDTO.NewBaseResponse(http.StatusSeeOther, false, "/").SendResponse(&w)

}

func (u *UserAuthController) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/login/login.html"))

	var err = tmpl.Execute(w, nil)

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}

func (u *UserAuthController) signUpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := dto.SignUpRequest{}

	if err := decoder.Decode(&payload); err != nil {
		log.Println("[Decode] Error decoding JSON")
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	if len(payload.Username) > 20 {
		globalDTO.NewBaseResponse(http.StatusBadRequest, true, "Username must be shorter than 20 character").SendResponse(&w)
		return
	}

	err := u.userService.RegisterUser(r.Context(), payload)
	if err == nil {
		log.Println("[Auth] Success :", payload.Username, "created")
		globalDTO.NewBaseResponse(http.StatusSeeOther, false, "/login").SendResponse(&w)
		return
	}

	if !strings.Contains(err.Error(), "Duplicate") {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	log.Println("[Auth] Failed creating account :", payload.Username, "Already exist")
	globalDTO.NewBaseResponse(http.StatusOK, true, fmt.Sprintf("Accout with username %s already exist", payload.Username)).SendResponse(&w)

}

func (u *UserAuthController) signUpPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/signup/signup.html"))

	var err = tmpl.Execute(w, nil)

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}

func (u *UserAuthController) logOutHandler(w http.ResponseWriter, r *http.Request) {
	if storedCookie, _ := r.Cookie("session_token"); storedCookie != nil {
		u.sessionService.DeleteSession(r.Context(), storedCookie.Value)
	}

	utils.DeleteSessionCookie(&w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
