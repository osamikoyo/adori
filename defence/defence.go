package defence

import (
	"net/http"
	"strings"
	"sync"

	"github.com/osamikoyo/adori/config"
	"github.com/osamikoyo/adori/logger"
	"go.uber.org/zap"
)

const (
	RequestFromOneIpInSecond = 10
)

type Defence struct {
	mutex sync.RWMutex

	suspiciousIp []string
	badPathParts []string
	iptable      map[string]uint
	logger       *logger.Logger
}

func NewDefence(cfg *config.Config, logger *logger.Logger) *Defence {
	return &Defence{
		suspiciousIp: cfg.Defence.BlackList,
		badPathParts: cfg.Defence.BadRequestParts,
		iptable: make(map[string]uint),
	}
}

func in(elem string, arr []string) bool {
	for _, e := range arr {
		if elem == e {
			return true
		}
	}

	return false
}

func (d *Defence) haveBadRequestPath(r *http.Request) bool {
	path := r.URL.Path

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if in(part, d.badPathParts) {
			return true
		}
	}

	return false
}

func (d *Defence) CheckRequestOK(r *http.Request) bool {
	ip := r.RemoteAddr

	if d.haveBadRequestPath(r) {
		return false
	}

	if in(ip, d.suspiciousIp) {
		d.logger.Info("request from suspicious ip list",
			zap.String("ip", ip),
			zap.String("addr", r.RequestURI))

		return false
	}

	d.mutex.Lock()
	d.iptable[ip]++
	d.mutex.Unlock()

	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if d.iptable[ip] <= RequestFromOneIpInSecond {
		d.logger.Error("request from ok url",
			zap.String("ip", ip),
			zap.String("addr", r.RequestURI))

		return true
	} else {
		d.logger.Error("new ip added to suspicious list",
			zap.String("id", ip))

		d.suspiciousIp = append(d.suspiciousIp, ip)

		return false
	}
}
