package controller

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

func (p *PostController) postDashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/homepage/index.html"))
	var err error

	latestPosts, err := p.postService.GetBriefsBlogPost(r.Context(), 1, 6)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	latestCategories, err := p.postService.GetTopCategoryPost(r.Context())
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	editorsPick, err := p.postService.GetEditorsPick(r.Context())
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}
	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn":         isLoggedIn,
		"EditorsPick":      editorsPick,
		"LatestPosts":      latestPosts,
		"LatestCategories": latestCategories,
	}

	if isLoggedIn {
		session, err := p.sessionService.GetSession(r.Context(), token)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}
		user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}

		if err == nil {
			data["User"] = user
		}
	}

	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (p *PostController) searchPostPageHandler(w http.ResponseWriter, r *http.Request) {
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
	category := queryVar.Get("category")

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

	limitConv, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}
	pageConv, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	contentCount, err := p.postService.GetCountOfSearchResult(r.Context(), q, dateStartPtr, dateEndPtr, category)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}
	maxPage := int64(math.Ceil(float64(contentCount) / float64(limitConv)))
	pageNavigation := utils.Paginate(pageConv, maxPage)
	startIndex := (pageConv-1)*limitConv + 1

	postData, err := p.postService.SearchBlogPost(r.Context(), q, pageConv, limitConv, dateStartPtr, dateEndPtr, category)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	categoryList, err := p.postService.GetCategoryList(r.Context())
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn":   isLoggedIn,
		"Query":      q,
		"PostsCount": contentCount,
		"StartIndex": startIndex,
		"EndIndex":   startIndex + int64(len(postData)) - 1,
		"Posts":      postData,
		"PageNav":    pageNavigation,
		"Categories": categoryList,
	}

	if isLoggedIn {
		session, err := p.sessionService.GetSession(r.Context(), token)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
		}

		user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
		}

		if err == nil {
			data["User"] = user
		}
	}

	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (p *PostController) viewPostPageHandlder(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/blog-view/blog-view.html"))

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	postData, err := p.postService.GetFullPost(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "No post") {
			panic(globalDTO.NewBaseResponse(http.StatusNotFound, true, err.Error()))
		}
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn": isLoggedIn,
		"Posts":    postData,
	}

	if isLoggedIn {
		session, err := p.sessionService.GetSession(r.Context(), token)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}
		user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}

		if err == nil {
			data["User"] = user
		}
	}

	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	err = p.postService.IncreaseView(r.Context(), id)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (p *PostController) createPostPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/post/createPost.html"))

	token, _ := utils.GetSessionToken(r)
	session, err := p.sessionService.GetSession(r.Context(), token)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}
	user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	categoryList, err := p.postService.GetCategoryList(r.Context())
	data := map[string]interface{}{
		"Categories": categoryList,
		"User":       user,
	}

	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, nil))
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (p *PostController) latestPostPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/see-more/see-more.html"))
	var err error

	queryVar := r.URL.Query()
	limit := queryVar.Get("limit")
	if limit == "" {
		limit = "12"
	}
	page := queryVar.Get("page")
	if page == "" {
		page = "1"
	}

	limitConv, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}
	pageConv, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	contentCount, err := p.postService.GetCountListOfPost(r.Context())
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}
	maxPage := int64(math.Ceil(float64(contentCount) / float64(limitConv)))
	pageNavigation := utils.Paginate(pageConv, maxPage)
	startIndex := (pageConv-1)*limitConv + 1

	postData, err := p.postService.GetBriefsBlogPost(r.Context(), pageConv, limitConv)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn":   isLoggedIn,
		"PostsCount": contentCount,
		"StartIndex": startIndex,
		"EndIndex":   startIndex + int64(len(postData)) - 1,
		"Posts":      postData,
		"PageNav":    pageNavigation,
		"In":         "Latest",
	}

	if isLoggedIn {
		session, err := p.sessionService.GetSession(r.Context(), token)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}
		user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}

		if err == nil {
			data["User"] = user
		}
	}

	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}

func (p *PostController) postInCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/see-more/see-more.html"))
	var err error

	queryVar := r.URL.Query()
	vars := mux.Vars(r)
	category := vars["category"]

	limit := queryVar.Get("limit")
	if limit == "" {
		limit = "12"
	}
	page := queryVar.Get("page")
	if page == "" {
		page = "1"
	}

	limitConv, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}
	pageConv, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
	}

	contentCount, err := p.postService.GetCountListOfPostInCategories(r.Context(), category)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}
	maxPage := int64(math.Ceil(float64(contentCount) / float64(limitConv)))
	pageNavigation := utils.Paginate(pageConv, maxPage)
	startIndex := (pageConv-1)*limitConv + 1

	postData, err := p.postService.GetBriefsBlogPostOfCategories(r.Context(), category, pageConv, limitConv)
	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	token, isLoggedIn := utils.GetSessionToken(r)
	data := map[string]interface{}{
		"LoggedIn":   isLoggedIn,
		"PostsCount": contentCount,
		"StartIndex": startIndex,
		"EndIndex":   startIndex + int64(len(postData)) - 1,
		"Posts":      postData,
		"PageNav":    pageNavigation,
		"In":         fmt.Sprint("Latest in ", strings.Title(category)),
	}

	if isLoggedIn {
		session, err := p.sessionService.GetSession(r.Context(), token)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}
		user, err := p.userService.GetUserMiniDetail(r.Context(), session.UID)
		if err != nil {
			panic(globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()))
		}

		if err == nil {
			data["User"] = user
		}
	}

	if err != nil {
		panic(globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()))
	}

	tmpl.Execute(w, globalDTO.NewBaseResponse(http.StatusOK, false, data))
}
