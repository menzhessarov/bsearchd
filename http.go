package bsearchd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type ValueStore interface {
	GetIndex(value int) (Entry, error)
}

type HTTPServer struct {
	Server *http.Server
	store  ValueStore
}

func NewHTTPServer(port string, store ValueStore) *HTTPServer {
	return &HTTPServer{
		Server: &http.Server{
			Addr: fmt.Sprintf("localhost:%s", port),
		},
		store: store,
	}
}

func (s *HTTPServer) RegisterRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /values/{value}", s.GetIndex)

	s.Server.Handler = mux
}

type response struct {
	Index   *int   `json:"index,omitempty"`
	Value   *int   `json:"value,omitempty"`
	Message string `json:"msg,omitempty"`
}

func (s *HTTPServer) GetIndex(w http.ResponseWriter, r *http.Request) {
	value, err := strconv.Atoi(r.PathValue("value"))
	if err != nil {
		slog.ErrorContext(r.Context(), "invalid value", "err", err, "value", r.PathValue("value"))

		if err = encode(w, http.StatusBadRequest, &response{Value: &value, Message: "invalid value"}); err != nil {
			slog.ErrorContext(r.Context(), "encode error message", "err", err)
		}
		return
	}

	entry, err := s.store.GetIndex(value)
	if err != nil {
		slog.ErrorContext(r.Context(), "get index", "err", err, "value", value)

		status := http.StatusInternalServerError
		msg := "internal error"

		if errors.Is(err, ErrNotFound) {
			status = http.StatusNotFound
			msg = "not found"
		}

		if err = encode(w, status, &response{Value: &value, Message: msg}); err != nil {
			slog.ErrorContext(r.Context(), "encode error message", "err", err)
		}
		return
	}

	resp := &response{
		Index: &entry.Index,
		Value: &entry.Value,
	}

	if err = encode(w, http.StatusOK, resp); err != nil {
		slog.ErrorContext(r.Context(), "encode response", "err", err)
	}
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
