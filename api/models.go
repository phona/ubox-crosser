package api

import "ubox-crosser/models/connection"

type ConnectionListResponse struct {
	Connections []*connection.Connection `json:"connections"`
	Pagination  Pagination              `json:"pagination"`
}

type ConnectionResponse struct {
	Connection *connection.Connection `json:"connection"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
