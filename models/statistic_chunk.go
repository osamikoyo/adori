package models

import (
	"net/http"
	"time"
)

type StatisticChunk struct {
	RequestPath string    `json:"request_path"`
	IP          string    `json:"ip"`
	Timestamp   time.Time `json:"timestamp"`
	Protocol    string    `json:"protocol"`
	Status      string    `json:"status"`
}

func NewStatisticChunk(r *http.Request,status string) *StatisticChunk {
	return &StatisticChunk{
		RequestPath: r.URL.String(),
		IP:          r.RemoteAddr,
		Protocol:    r.Proto,
		Timestamp:   time.Now(),
		Status:      status,
	}
}
