package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/auth/service"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	"github.com/suryaadi44/Techdo-blog/internal/post/service"
	userServicePkg "github.com/suryaadi44/Techdo-blog/internal/user/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

type PostController struct {
	router         *mux.Router
	postService    service.PostServiceApi
	sessionService authServicePkg.SessionServiceApi
	userService    userServicePkg.UserServiceApi
}

func NewController(router *mux.Router, postService service.PostServiceApi, sessionService authServicePkg.SessionServiceApi, userService userServicePkg.UserServiceApi) *PostController {
	return &PostController{
		router:         router,
		postService:    postService,
		sessionService: sessionService,
		userService:    userService,
	}
}

func (p *PostController) InitializeController() {
	p.router.HandleFunc("/post/create", p.createPostPageHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/post/create", p.createPostHandler).Methods(http.MethodPost)
	p.router.HandleFunc("/", p.postDashboardHandler).Methods(http.MethodGet)
}

func (p *PostController) postDashboardHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/index/index.html"))
	var err error

	queryVar := r.URL.Query()
	limit := queryVar.Get("limit")
	if limit == "" {
		limit = "8"
	}
	page := queryVar.Get("page")
	if page == "" {
		page = "1"
	}
	limitConv, _ := strconv.ParseInt(limit, 10, 64)
	pageConv, _ := strconv.ParseInt(page, 10, 64)

	postData, err := p.postService.GetBriefsBlogPost(r.Context(), pageConv, limitConv)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn": isLoggedIn,
		"Posts":    postData,
	}

	if isLoggedIn {
		session, err := p.sessionService.GetSession(r.Context(), token)
		user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)

		if err == nil {
			data["User"] = user
		}
	}

	if err == nil {
		tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
	} else {
		tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusInternalServerError, true, nil))
	}
}

func (p *PostController) createPostPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/post/createPost.html"))

	token, _ := utils.GetSessionToken(r)
	session, err := p.sessionService.GetSession(r.Context(), token)
	user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)

	categoryList, err := p.postService.GetCategoryList(r.Context())
	data := map[string]interface{}{
		"Categories": categoryList,
		"User":       user,
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

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	category, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	uploadedFile, handler, err := r.FormFile("banner")
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

	post := dto.BlogPostRequest{
		Category:   category,
		Banner:     buf.Bytes(),
		BannerName: handler.Filename,
		Title:      r.FormValue("title"),
		Body:       r.FormValue("editordata"),
	}

	postID, err := p.postService.AddPost(r.Context(), post, session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusCreated, false, fmt.Sprintf("/post/%d", postID)).SendResponse(&w)
}
