package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// processBatchItems godoc
// @summary Process batch of items
// @description Send items for processing
// @tags Item
// @accept json
// @produce json
// @param req body ProcessBatchItemsRequest true "Request object"
// @success 200 {array} OkResponse
// @failure 400,404,500 {object} ErrorResponse
// @router /items/batch [post]
func (s *Server) processBatchItems(w http.ResponseWriter, req *http.Request) {
	var data *ProcessBatchItemsRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		log.Error().Err(fmt.Errorf("http server: processBatchItems: %w", err)).Msg("request decoding error")
		w.WriteHeader(http.StatusBadRequest)
		resp := ErrorResponse{Error: "request decoding error"}
		jData, _ := json.Marshal(resp)
		_, _ = w.Write(jData)
		return
	}

	if err := s.uc.Process(req.Context(), data.Items); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := ErrorResponse{Error: "items processing error"}
		jData, _ := json.Marshal(resp)
		_, _ = w.Write(jData)
		return
	}

	resp, _ := json.Marshal(OkResponse{Ok: true})
	_, _ = w.Write(resp)
}
