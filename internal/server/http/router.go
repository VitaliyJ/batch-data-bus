package http

import (
	"net/http"

	"data-bus/internal/server/http/middleware"
	"github.com/gorilla/mux"
)

// OpenAPI documentation info:
// @title       Items processing API
// @version     1.0
// @description This API represents interface for sending items to processing
// @BasePath    /
func (s *Server) newRouter() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(NotFound)
	r.Use(middleware.JSONResponse)

	// routes settings here
	r.HandleFunc("/items/batch", s.processBatchItems).Methods(http.MethodPost)

	return r
}
