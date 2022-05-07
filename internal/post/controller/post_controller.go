package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/auth/service"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/internal/post/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
)

type PostController struct {
	router         *mux.Router
	postService    service.PostServiceApi
	sessionService authServicePkg.SessionServiceApi
}

func NewController(router *mux.Router, postService service.PostServiceApi, sessionService authServicePkg.SessionServiceApi) *PostController {
	return &PostController{
		router:         router,
		postService:    postService,
		sessionService: sessionService,
	}
}

func (p *PostController) InitializeController() {
	p.router.HandleFunc("/post/create", p.createPostPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/post/create", p.createPostHandler).Methods(http.MethodPost)
	p.router.HandleFunc("/", p.postDashboardHandler).Methods(http.MethodGet)
}

func (p *PostController) createPostPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/post/createPost.html"))

	categoryList, err := p.postService.GetCategoryList(r.Context())
	data := map[string]interface{}{
		"Categories": categoryList,
	}

	if err == nil {
		err = tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
	} else {
		err = tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusInternalServerError, true, nil))
	}

	if err != nil {
		tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusInternalServerError, true, nil))
		return
	}
}

func (p *PostController) createPostHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := true
	c := &http.Cookie{}
	if storedCookie, _ := r.Cookie("session_token"); storedCookie != nil {
		c = storedCookie
	}
	if c.Value == "" {
		loggedIn = false
	}

	if err := r.ParseForm(); err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	category, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
	if err != nil {
		loggedIn = false
	}

	post := dto.BlogPostRequest{
		Category: category,
		Banner:   r.FormValue("cover"),
		Title:    r.FormValue("title"),
		Body:     r.FormValue("editordata"),
	}

	session, _ := p.sessionService.GetSession(r.Context(), c.Value)

	if !loggedIn {
		globalDTO.NewBaseResponse(http.StatusForbidden, true, "Authentication required").SendResponse(&w)
		return
	}

	postID, err := p.postService.AddPost(r.Context(), post, session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusCreated, false, fmt.Sprintf("/post/%d", postID)).SendResponse(&w)
}

func (p *PostController) postDashboardHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/index/index.html"))

	var err = tmpl.Execute(w, nil)

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}
