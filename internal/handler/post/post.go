package postHandler

import (
	"adinata/internal/model/request"
	"adinata/internal/model/response"
	postService "adinata/internal/service/post"
	"adinata/middleware"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type postHandler struct {
	postService postService.PostServiceProvider
}

func RegisterPostHandler(mux *http.ServeMux, postService postService.PostServiceProvider) {
	handler := &postHandler{
		postService: postService,
	}

	mux.Handle("GET /api/posts/", middleware.VerifyToken(http.HandlerFunc(handler.GetPostHandler)))
	mux.Handle("POST /api/posts", middleware.VerifyToken(http.HandlerFunc(handler.CreatePostHandler)))
	mux.Handle("PUT /api/posts/{post_id}", middleware.VerifyToken(http.HandlerFunc(handler.UpdatePostHandler)))
	mux.Handle("DELETE /api/posts/{post_id}", middleware.VerifyToken(http.HandlerFunc(handler.DeletePostHandler)))
	mux.Handle("GET /api/posts/{post_id}", middleware.VerifyToken(http.HandlerFunc(handler.GetPostDetailHandler)))
}

func (p *postHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Fail when reading body", "Error", err)
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return

	}
	defer r.Body.Close()

	req := request.PostRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		return
	}

	err = p.postService.CreatePost(r.Context(), req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Create Post Success"})

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (p *postHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postId, err := strconv.ParseInt(r.PathValue("post_id"), 10, 64)
	if err != nil {
		slog.Error("Fail when parsing path", "Error", err)
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Fail when reading body", "Error", err)
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	defer r.Body.Close()

	req := request.PostRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		return
	}

	err = p.postService.UpdatePost(r.Context(), req, postId)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Update Post Success"})

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (p *postHandler) GetPostDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	postId, err := strconv.ParseInt(r.PathValue("post_id"), 10, 64)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	postData, err := p.postService.GetPostDetail(r.Context(), postId)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GetPostByIdResponse{GenericResponse: response.GenericResponse{Success: true, Message: "Create Post Success"}, Data: postData})

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}

func (p *postHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postId, err := strconv.ParseInt(r.PathValue("post_id"), 10, 64)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	err = p.postService.DeletePost(r.Context(), postId)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.Write(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Deleting Data"})
	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

func (p *postHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	searchParams := r.URL.Query().Get("search")

	postData, err := p.postService.GetPosts(r.Context(), searchParams)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GetPostResponse{GenericResponse: response.GenericResponse{Success: true, Message: "Get Post Success"}, Data: postData})

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}
