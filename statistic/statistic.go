package statistic

import (
	"context"
	"os"

	"github.com/bytedance/sonic"
	"github.com/osamikoyo/adori/logger"
	"github.com/osamikoyo/adori/models"
	"go.uber.org/zap"
)

type StatisticWriter struct {
	input chan *models.StatisticChunk

	file   *os.File
	logger *logger.Logger
}

func NewStatisticWriter(
	filepath string,
	logger *logger.Logger,
	input chan *models.StatisticChunk,
) (*StatisticWriter, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		logger.Error("failed open statistic file",
			zap.String("filepath", filepath),
			zap.Error(err))

		return nil, err
	}

	return &StatisticWriter{
		input:  input,
		file:   file,
		logger: logger,
	}, nil
}

func (sw *StatisticWriter) Close() error {
	sw.logger.Info("stopping statistic writer...")

	return sw.file.Close()
}

func (sw *StatisticWriter) Listen(ctx context.Context) {
	sw.logger.Info("starting statistic writer...")

	encoder := sonic.ConfigDefault.NewEncoder(sw.file)

	for {
		select {
		case <-ctx.Done():
			sw.Close()
			return
		case chunk := <-sw.input:
			if err := encoder.Encode(chunk); err != nil {
				sw.logger.Error("failed encode chunk to statistic",
					zap.Any("chunk", chunk),
					zap.Error(err))

				continue
			}
		}
	}
}
