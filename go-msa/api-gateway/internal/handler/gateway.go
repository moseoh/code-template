package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type GatewayHandler struct {
	authServiceURL     *url.URL
	businessServiceURL *url.URL
}

func NewGatewayHandler(authServiceURL, businessServiceURL string) (*GatewayHandler, error) {
	authURL, err := url.Parse(authServiceURL)
	if err != nil {
		return nil, err
	}

	businessURL, err := url.Parse(businessServiceURL)
	if err != nil {
		return nil, err
	}

	return &GatewayHandler{
		authServiceURL:     authURL,
		businessServiceURL: businessURL,
	}, nil
}

func (g *GatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var proxy *httputil.ReverseProxy
	path := r.URL.Path

	switch {
	case strings.HasPrefix(path, "/auth"):
		proxy = httputil.NewSingleHostReverseProxy(g.authServiceURL)
	case strings.HasPrefix(path, "/products"):
		proxy = httputil.NewSingleHostReverseProxy(g.businessServiceURL)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	proxy.ServeHTTP(w, r)
}
