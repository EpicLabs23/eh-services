package models

type Zone struct {
	Name    string   `json:"name" binding:"required"`
	Records []Record `json:"records,omitempty"`
}

type Record struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl"`
}

type ZoneUpdate struct {
	Records []Record `json:"records" binding:"required"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
