package http

type OkResponse struct {
	Ok bool `json:"ok"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
