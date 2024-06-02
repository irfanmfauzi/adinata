package tagHandler

import (
	"adinata/internal/model/request"
	"adinata/internal/model/response"
	"adinata/internal/service/tag"
	"adinata/middleware"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type tagHandler struct {
	tagService tagService.TagServiceProvider
}

func RegisterTagRoute(mux *http.ServeMux, service tagService.TagServiceProvider) {
	t := tagHandler{
		tagService: service,
	}

	mux.Handle("POST /api/tags", middleware.VerifyToken(http.HandlerFunc(t.CreateTagHandler)))
	mux.Handle("GET /api/tags", middleware.VerifyToken(http.HandlerFunc(t.GetTagHandler)))
	mux.Handle("PUT /api/tags/{id}", middleware.VerifyToken(http.HandlerFunc(t.GetTagHandler)))
	mux.Handle("DELETE /api/tags/{id}", middleware.VerifyToken(http.HandlerFunc(t.GetTagHandler)))
}

func (t *tagHandler) GetTagHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tags, err := t.tagService.GetTag(r.Context())
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.Write(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(tags)
	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

func (t *tagHandler) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}
	defer r.Body.Close()

	req := request.TagRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	err = t.tagService.CreateTag(r.Context(), req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Creating Tag"})

	w.Write(resp)
	w.WriteHeader(http.StatusCreated)

}

func (t *tagHandler) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tagId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}
	defer r.Body.Close()

	req := request.TagRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	err = t.tagService.UpdateTag(r.Context(), req, tagId)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Creating Tag"})

	w.Write(resp)
	w.WriteHeader(http.StatusCreated)

}

func (t *tagHandler) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tagId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	err = t.tagService.DeleteTag(r.Context(), tagId)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Delete Tag"})

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}
