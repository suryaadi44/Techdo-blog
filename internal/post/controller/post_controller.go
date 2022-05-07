package controller

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/auth/service"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/internal/post/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
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
	token, _ := utils.GetSessionToken(r)
	session, _ := p.sessionService.GetSession(r.Context(), token)
	// if !loggedIn {
	// 	globalDTO.NewBaseResponse(http.StatusForbidden, true, "Authentication required").SendResponse(&w)
	// 	return
	// }

	if err := r.ParseForm(); err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	category, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	post := dto.BlogPostRequest{
		Category: category,
		Banner:   r.FormValue("cover"),
		Title:    r.FormValue("title"),
		Body:     r.FormValue("editordata"),
	}

	err = p.postService.AddPost(r.Context(), post, session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusCreated, false, nil).SendResponse(&w)
}

func (p *PostController) postDashboardHandler(w http.ResponseWriter, r *http.Request) {
	token, isLoggedIn := utils.GetSessionToken(r)
	session, _ := p.sessionService.GetSession(r.Context(), token)

	//TODO : Get user detail and pass its value to front end

	var tmpl = template.Must(template.ParseFiles("web/template/index/index.html"))

	data := map[string]interface{}{
		"LoggedIn": isLoggedIn,
		"User":     session, //only placehodler
	}

	var err = tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))

	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}
