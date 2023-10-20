package http

import "data-bus/pkg/externalservice"

type ProcessBatchItemsRequest struct {
	Items []externalservice.Item `json:"items"`
}
