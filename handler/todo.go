package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (t *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		todoReq := &model.CreateTODORequest{}
		todoRes := &model.CreateTODOResponse{}
		if err := json.NewDecoder(r.Body).Decode(todoReq); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if todoReq.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := t.svc.CreateTODO(r.Context(), todoReq.Subject, todoReq.Description)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		todoRes.TODO = *todo
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "PUT":
		todoReq := &model.UpdateTODORequest{}
		todoRes := &model.UpdateTODOResponse{}
		if err := json.NewDecoder(r.Body).Decode(todoReq); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if todoReq.Subject == "" || todoReq.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := t.svc.UpdateTODO(r.Context(), todoReq.ID, todoReq.Subject, todoReq.Description)
		if errors.Is(err, &model.ErrNotFound{}) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		todoRes.TODO = *todo
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "GET":
		todoReq := &model.ReadTODORequest{}
		todoRes := &model.ReadTODOResponse{}
		todoReq.PrevID, _ = strconv.ParseInt(r.URL.Query().Get("prev_id"), 10, 64)
		todoReq.Size, _ = strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)
		//log.Println(todoReq)
		if todoReq.Size == 0 {
			todoReq.Size = 5
		}
		todos, err := t.svc.ReadTODO(r.Context(), todoReq.PrevID, todoReq.Size)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		todoRes.TODOs = todos
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "DELETE":
		todoReq := &model.DeleteTODORequest{}
		todoRes := &model.DeleteTODOResponse{}
		if err := json.NewDecoder(r.Body).Decode(todoReq); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(todoReq.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := t.svc.DeleteTODO(r.Context(), todoReq.IDs)
		if errors.Is(err, &model.ErrNotFound{}) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
