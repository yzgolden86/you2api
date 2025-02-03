package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Proxy struct {
	target     *url.URL
	proxy      *httputil.ReverseProxy
	timeoutMS  int
}

func NewProxy(targetURL string, timeoutMS int) (*Proxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		target:    target,
		timeoutMS: timeoutMS,
	}

	proxy.proxy = &httputil.ReverseProxy{
		Director: proxy.director,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(target),
			ResponseHeaderTimeout: time.Duration(timeoutMS) * time.Millisecond,
		},
	}

	return proxy, nil
}

func (p *Proxy) director(req *http.Request) {
	req.URL.Scheme = p.target.Scheme
	req.URL.Host = p.target.Host
	req.Host = p.target.Host
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
} 