package controller

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suryaadi44/Techdo-blog/internal/post/service"
	globalDTO "github.com/suryaadi44/Techdo-blog/pkg/dto"
)

type PostController struct {
	router      *mux.Router
	postService service.PostServiceApi
}

func NewController(router *mux.Router, postService service.PostServiceApi) *PostController {
	return &PostController{
		router:      router,
		postService: postService,
	}
}

func (p *PostController) InitializeController() {
	p.router.HandleFunc("/post/create", p.createPostHandler).Methods(http.MethodGet)
	p.router.HandleFunc("/", p.postDashboardHandler).Methods(http.MethodGet)

	p.router.HandleFunc("/api/upload/image", p.uploadImageHandler).Methods(http.MethodPost)
	//p.router.HandleFunc("/api/upload/post", p.uploadPostHandler).Methods(http.MethodPost)
}

func (p *PostController) createPostHandler(w http.ResponseWriter, r *http.Request) {
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

func (p *PostController) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	part, err := reader.NextPart()
	if err != nil {
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	file := StreamToByte(part)

	response, err := p.postService.UploadImage(r.Context(), part.FileName(), file)
	if err != nil {
		log.Println("[Post] Error Uploading image")
		globalDTO.NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}

	globalDTO.NewBaseResponse(http.StatusOK, false, response).SendResponse(&w)
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
