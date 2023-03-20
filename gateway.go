package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Gateway function that extract the targeted service and proxies the request to it
func Gateway(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	service, err := ExtractService(path)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	serviceUrl, err := url.Parse(service)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO authentication happens here and pass user object to Create ReverseProxy
	CreateReverseProxy(serviceUrl).ServeHTTP(ctx.Writer, ctx.Request)

}

func ExtractService(path string) (string, error) {
	// The passed path should be /api/{serviceName}/{serviceNameSpace}/*
	split := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(split) <= 1 {
		return "", fmt.Errorf("failed to parse target service from path: %s", path)
	}
	serviceHost := fmt.Sprintf("svc-%s", split[1])
	serviceNameSpace := fmt.Sprintf("svc-%s", split[2])
	if serviceHost == "" {
		return "", fmt.Errorf("failed to parse target  from path: %s", path)
	}

	// Return the interal k8s address for the found service
	return fmt.Sprintf(
		"http://%s.%s:%d/api/%s",
		serviceHost, serviceNameSpace, 10000, strings.Join(split[3:], "/"),
	), nil
}

// Should take in user object which it can pass in request headers
func CreateReverseProxy(address *url.URL) *httputil.ReverseProxy {
	p := httputil.NewSingleHostReverseProxy(address)
	p.Director = func(request *http.Request) {
		request.Host = address.Host
		request.URL.Scheme = address.Scheme
		request.URL.Host = address.Host
		request.URL.Path = address.Path
		// Add request header(s) with information about the user like ID and Role

	}

	// Handle Responses by logging errors or changing status codes in production with p.ModifyResponse

	return p
}
