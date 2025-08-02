package core

import (
	"net/http"

	"github.com/osamikoyo/adori/cash"
	"github.com/osamikoyo/adori/defence"
	"github.com/osamikoyo/adori/logger"
	"github.com/osamikoyo/adori/statistic"
	"go.uber.org/zap"
)

type AdoriCore struct {
	cash      *cash.LocalCash
	defence   *defence.Defence
	statistic *statistic.StatisticClient

	logger *logger.Logger
}

func NewAdoriCore(
	cash *cash.LocalCash,
	defence *defence.Defence,
	statistic *statistic.StatisticClient,
	logger *logger.Logger,
) *AdoriCore {
	return &AdoriCore{
		cash: cash,
		defence: defence,
		statistic: statistic,
		logger: logger,
	}
}

func (ac *AdoriCore) CoreMiddlewareForHandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var status string

		defer ac.statistic.AddChunk(r, status)

		if !ac.defence.CheckRequestOK(r) {
			w.Write([]byte("bad request"))

			status = "bad"

			return
		}

		path := r.URL.Path

		cash, err := ac.cash.Get(path)
		if err == nil{
			w.Header().Set("Content-Type", http.DetectContentType([]byte(cash.Content)))

			w.Write([]byte(cash.Content))
			
			status = "ok (cashed)"

			return 
		}

		status = "ok"

		ac.logger.Info("new good request", zap.String("path", r.URL.Path))

		handler.ServeHTTP(w, r)
	}
}

func (ac *AdoriCore) CoreMiddlewareForHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){ 
		var status string

		defer ac.statistic.AddChunk(r, status)

		if !ac.defence.CheckRequestOK(r) {
			w.Write([]byte("bad request"))

			status = "bad"

			return
		}

		path := r.URL.Path

		cash, err := ac.cash.Get(path)
		if err == nil{
			w.Header().Set("Content-Type", http.DetectContentType([]byte(cash.Content)))

			w.Write([]byte(cash.Content))

			status = "ok (cashed)"

			return 
		}


		status = "ok"

		ac.logger.Info("new good request", zap.String("path", r.URL.Path))

		handler.ServeHTTP(w, r)
	})
}