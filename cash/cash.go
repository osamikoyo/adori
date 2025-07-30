package cash

import (
	"errors"
	"sync"
	"time"

	"github.com/osamikoyo/adori/logger"
	"github.com/osamikoyo/adori/models"
	"go.uber.org/zap"
)

const AddIntervalMinutes = 10

var ErrNotFound = errors.New("not found")

type LocalCash struct {
	stop chan struct{}

	logger *logger.Logger

	mu    sync.RWMutex
	wg    sync.WaitGroup
	files map[string]*models.YadoriFile
}

func NewLocalCash(cleanupInterval time.Duration, logger *logger.Logger) *LocalCash {
	cash := &LocalCash{
		stop:   make(chan struct{}),
		files:  make(map[string]*models.YadoriFile),
		logger: logger,
	}

	cash.wg.Add(1)

	go func(inteval time.Duration) {
		defer cash.wg.Done()

		cash.clean(inteval)
	}(cleanupInterval)

	return cash
}

func (l *LocalCash) clean(inteval time.Duration) {
	t := time.NewTicker(inteval)
	defer t.Stop()

	for {
		select {
		case <-l.stop:
			return
		case <-t.C:
			l.mu.Lock()
			for path, file := range l.files {
				if file.ExpireAtTimestamp <= time.Now().Unix() {
					delete(l.files, path)

					l.logger.Info("file deleted from cash",
						zap.String("path", path))
				}
			}
			l.mu.Unlock()
		}
	}
}

func (l *LocalCash) Add(file *models.YadoriFile) {
	if file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	file.ExpireAtTimestamp = time.Now().Add(AddIntervalMinutes * time.Second).Unix()

	l.files[file.Path] = file
}

func (l *LocalCash) Get(path string) (*models.YadoriFile, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	file, ok := l.files[path]
	if !ok {
		return nil, ErrNotFound
	}

	return file, nil
}
