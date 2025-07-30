package defence

import (
	"net/http"
	"sync"

	"github.com/osamikoyo/adori/logger"
	"go.uber.org/zap"
)

const (
	RequestFromOneIpInSecond = 10
)

type Defence struct {
	mutex sync.RWMutex

	suspiciousIp []string
	IPtable      map[string]uint
	logger       *logger.Logger
}

func in(elem string, arr []string) bool {
	for _, e := range arr {
		if elem == e {
			return true
		}
	}

	return false
}

func (d *Defence) CheckRequestOK(r *http.Request) bool {
	ip := r.RemoteAddr

	if in(ip, d.suspiciousIp) {
		d.logger.Info("request from suspicious ip list",
			zap.String("ip", ip),
			zap.String("addr", r.RequestURI))

		return false
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.IPtable[ip]++

	if d.IPtable[ip] <= RequestFromOneIpInSecond {
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
