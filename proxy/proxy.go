package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/osamikoyo/adori/config"
)

const AdoriHeader = "Adori-Proxy"

type Proxy struct {
	mutex                sync.RWMutex
	prefixToIps          map[string][]string
	ipToProxy            map[string]*httputil.ReverseProxy
	ipRequestCounters    map[string]uint
	healthChecks         map[string]bool
}

func NewProxy(cfg *config.Config) (*Proxy, error) {
	p := &Proxy{
		prefixToIps:       make(map[string][]string),
		ipToProxy:         make(map[string]*httputil.ReverseProxy),
		ipRequestCounters: make(map[string]uint),
		healthChecks:      make(map[string]bool),
	}

	for _, gateway := range cfg.ApiGateway {
		p.prefixToIps[gateway.Prefix] = append(p.prefixToIps[gateway.Prefix], gateway.SelfPath)

		targetURL, err := url.Parse(gateway.OutputAddr)
		if err != nil {
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		p.initProxyDirector(proxy, targetURL.Host)
		p.ipToProxy[gateway.SelfPath] = proxy
		p.ipRequestCounters[gateway.SelfPath] = 0
		p.healthChecks[gateway.SelfPath] = true 
	}

	return p, nil
}

func (p *Proxy) initProxyDirector(proxy *httputil.ReverseProxy, targetHost string) {
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Proxy", AdoriHeader)
		req.Host = targetHost
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Proxy-Server", "Optimized-Gateway")
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		
		target := r.URL.Host
		p.healthChecks[target] = false
		
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Proxy error occurred"))
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	var matchedPrefix string
	for prefix := range p.prefixToIps {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			matchedPrefix = prefix
			break
		}
	}
	
	if matchedPrefix == "" {
		http.NotFound(w, r)
		return
	}

	target := p.balance(matchedPrefix)
	if target == "" {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	p.mutex.Lock()
	p.ipRequestCounters[target]++
	p.mutex.Unlock()

	proxy := p.ipToProxy[target]
	proxy.ServeHTTP(w, r)
}

func (p *Proxy) balance(prefix string) string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	availableIps, exists := p.prefixToIps[prefix]
	if !exists || len(availableIps) == 0 {
		return ""
	}

	var healthyIps []string
	for _, ip := range availableIps {
		if healthy, ok := p.healthChecks[ip]; ok && healthy {
			healthyIps = append(healthyIps, ip)
		}
	}

	if len(healthyIps) == 0 {
		return ""
	}

	var selected string
	minRequests := ^uint(0)
	for _, ip := range healthyIps {
		if count := p.ipRequestCounters[ip]; count < minRequests {
			minRequests = count
			selected = ip
		}
	}

	return selected
}

func (p *Proxy) HealthCheck() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for ip := range p.healthChecks {
		resp, err := http.Get("http://" + ip + "/health")
		if err != nil || resp.StatusCode != http.StatusOK {
			p.healthChecks[ip] = false
		} else {
			p.healthChecks[ip] = true
		}
	}
}