package statistic

import (
	"net/http"

	"github.com/osamikoyo/adori/models"
)

type StatisticClient struct {
	output chan *models.StatisticChunk
}

func NewStatisticClient(output chan *models.StatisticChunk) *StatisticClient {
	return &StatisticClient{
		output: output,
	}
}

func (s *StatisticClient) AddChunk(r *http.Request, status string) {
	chunk := models.NewStatisticChunk(r,status)

	s.output <- chunk
}
