package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/suryaadi44/Techdo-blog/internal/post/dto"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
	"github.com/suryaadi44/Techdo-blog/pkg/utils"
)

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

	globalDTO.NewBaseResponse(http.StatusOK, false, "Post deleted").SendResponse(&w)
}

func (p *PostController) createPostHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, err := p.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()).SendResponse(&w)
	}

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

func (p *PostController) userPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()).SendResponse(&w)
		return
	}

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

	postData, err := p.postService.GetMiniBlogPostsByUser(r.Context(), id, pageConv, limitConv)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, postData).SendResponse(&w)
}

func (p *PostController) userCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()).SendResponse(&w)
		return
	}

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

	postData, err := p.postService.GetCommentsByUser(r.Context(), id, pageConv, limitConv)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, postData).SendResponse(&w)
}

func (p *PostController) viewCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()).SendResponse(&w)
		return
	}

	commentsData, err := p.postService.GetComments(r.Context(), id)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, commentsData).SendResponse(&w)
}

func (p *PostController) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetSessionToken(r)
	session, err := p.sessionService.GetSession(r.Context(), token)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()).SendResponse(&w)
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusBadRequest, true, err.Error()).SendResponse(&w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	payload := dto.CommentRequest{}

	if err := decoder.Decode(&payload); err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	payload.UserID = session.UID
	payload.PostID = id

	err = p.postService.AddComment(r.Context(), payload)
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusCreated, false, nil).SendResponse(&w)
}
