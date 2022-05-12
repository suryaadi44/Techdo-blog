package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	authServicePkg "github.com/suryaadi44/Techdo-blog/internal/auth/service"
	middlewarePkg "github.com/suryaadi44/Techdo-blog/internal/middleware"
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
	authMiddleware middlewarePkg.AuthMiddleware
	userService    userServicePkg.UserServiceApi
}

func NewController(router *mux.Router, postService service.PostServiceApi, sessionService authServicePkg.SessionServiceApi, authMiddleware middlewarePkg.AuthMiddleware, userService userServicePkg.UserServiceApi) *PostController {
	return &PostController{
		router:         router,
		postService:    postService,
		sessionService: sessionService,
		authMiddleware: authMiddleware,
		userService:    userService,
	}
}

func (p *PostController) InitializeController() {
	createRouter := p.router.PathPrefix("/").Subrouter()
	createRouter.Use(p.authMiddleware.AuthMiddleware())
	createRouter.HandleFunc("/post/create", p.createPostPageHandler).Methods(http.MethodGet)
	createRouter.HandleFunc("/post/create", p.createPostHandler).Methods(http.MethodPost)
	createRouter.HandleFunc("/post/delete/{id:[0-9]+}", p.deletePostHandlder).Methods(http.MethodDelete)

	p.router.HandleFunc("/", p.postDashboardHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/search", p.searchBlogPostHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/post/{id:[0-9]+}", p.viewPostHandlder).Methods(http.MethodGet)
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

func (p *PostController) searchBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO : Change search page template
	var tmpl = template.Must(template.ParseFiles("web/template/search-blog/search-blog.html"))
	var err error
	var dateStart, dateEnd time.Time

	queryVar := r.URL.Query()
	limit := queryVar.Get("limit")
	if limit == "" {
		limit = "8"
	}
	page := queryVar.Get("page")
	if page == "" {
		page = "1"
	}

	start := queryVar.Get("start")
	end := queryVar.Get("end")

	dateStart, err = time.Parse("2006-01-02", start)
	dateStartPtr := &dateStart
	if err != nil || start == "" {
		dateStartPtr = nil
	}

	dateEnd, err = time.Parse("2006-01-02", end)
	dateEndPtr := &dateEnd
	if err != nil || end == "" {
		dateEndPtr = nil
	}

	q := queryVar.Get("q")

	limitConv, _ := strconv.ParseInt(limit, 10, 64)
	pageConv, _ := strconv.ParseInt(page, 10, 64)

	postData, err := p.postService.SearchBlogPost(r.Context(), q, pageConv, limitConv, dateStartPtr, dateEndPtr)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn": isLoggedIn,
		"Query":    q,
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

func (p *PostController) viewPostHandlder(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/blog-view/blog-view.html"))

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	postData, err := p.postService.GetFullPost(r.Context(), id)
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

func (p *PostController) deletePostHandlder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	token, _ := utils.GetSessionToken(r)
	session, err := p.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	postAuthor, err := p.postService.GetPostAuthorIdFromId(r.Context(), id)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	if postAuthor != session.UID {
		globalDTO.NewBaseResponse(http.StatusUnauthorized, true, "Cannot delete other user post").SendResponse(&w)
		return
	}

	err = p.postService.DeletePost(r.Context(), id)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, true, "Post deleted").SendResponse(&w)
}

func (p *PostController) createPostHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, _ := p.sessionService.GetSession(r.Context(), token)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	title := r.FormValue("title")
	body := r.FormValue("editordata")

	category, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, "Category error").SendResponse(&w)
		return
	}

	if strings.TrimSpace(title) == "" {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, "Title and Body are required").SendResponse(&w)
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
		Title:      title,
		Body:       body,
	}

	postID, err := p.postService.AddPost(r.Context(), post, session.UID)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusCreated, false, fmt.Sprintf("/post/%d", postID)).SendResponse(&w)
}
