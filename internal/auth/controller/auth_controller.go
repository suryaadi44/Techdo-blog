package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/suryaadi44/Techdo-blog/internal/auth/dto"
	"github.com/suryaadi44/Techdo-blog/internal/auth/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
)

type UserAuthController struct {
	router          *mux.Router
	userAuthService service.UserAuthServiceApi
}

func NewController(router *mux.Router, userAuthService service.UserAuthServiceApi) *UserAuthController {
	return &UserAuthController{
		router:          router,
		userAuthService: userAuthService,
	}
}

func (u *UserAuthController) InitializeController() {
	u.router.HandleFunc("/login", u.loginHandler).Methods(http.MethodPost)
	u.router.HandleFunc("/login", u.loginPageHandler).Methods(http.MethodGet)
	u.router.HandleFunc("/signup", u.signUpHandler).Methods(http.MethodPost)
	u.router.HandleFunc("/signup", u.signUpPageHandler).Methods(http.MethodGet)
}

func (u *UserAuthController) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := dto.AuthRequest{}

	if err := decoder.Decode(&payload); err != nil {
		log.Println("[Decode] Error decoding JSON")
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	session, err := u.userAuthService.AuthenticateUser(r.Context(), payload)
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
	return
}

func (u *UserAuthController) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/login/index.html"))

	var err = tmpl.Execute(w, nil)

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}

func (u *UserAuthController) signUpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := dto.AuthRequest{}

	if err := decoder.Decode(&payload); err != nil {
		log.Println("[Decode] Error decoding JSON")
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	err := u.userAuthService.RegisterUser(r.Context(), payload)
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
	return
}

func (u *UserAuthController) signUpPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/signup/index.html"))

	var err = tmpl.Execute(w, nil)

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}
