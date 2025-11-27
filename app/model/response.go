package model

// APIResponse adalah struktur response API standar
type APIResponse struct {
	Status  string      `json:"status"`  // "success" atau "error"
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PaginationResponse adalah struktur untuk response dengan pagination
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int         `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}
