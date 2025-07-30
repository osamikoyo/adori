package models

import (
	"net/http"
	"time"
)

type StatisticChunk struct {
	RequestPath string    `json:"request_path"`
	IP          string    `json:"ip"`
	Timestamp   time.Time `json:"timestamp"`
	Commpressed bool      `json:"commpressed"`
	Status      string    `json:"status"`
}

func NewStatisticChunk(r *http.Request, commpressed bool, status string) *StatisticChunk {
	return &StatisticChunk{
		RequestPath: r.URL.Path,
		IP: r.RemoteAddr,
		Timestamp: time.Now(),
		Commpressed: commpressed,
		Status: status,
	}
}