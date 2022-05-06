package controller

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
)

type PostController struct {
	router *mux.Router
}

func NewController(router *mux.Router) *PostController {
	return &PostController{
		router: router,
	}
}

func (p *PostController) InitializeController() {
	p.router.HandleFunc("/post/create", p.createPostHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/", p.postDashboardHandler).Methods(http.MethodGet)
}

func (p *PostController) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/post/createPost.html"))

	var err = tmpl.Execute(w, nil)

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}

func (p *PostController) postDashboardHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/index/index.html"))

	var err = tmpl.Execute(w, nil)

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}
