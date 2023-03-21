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

	// A few steps need to happen here:
	// -> Log the request (sending something to the logging service (probably rpc))
	// -> Check if the user is authenticated (probably rpc)
	// -> Check if the user is allowed to access this service (probably rpc)
	// -> If the user is allowed to access the service, proxy the request to the service
	//    likely will need to inject some user / permissions / role data.

	// potentially a proxy is not necessary and we can make the call as specified directly.

	CreateReverseProxy(serviceUrl).ServeHTTP(ctx.Writer, ctx.Request)
}

func ExtractService(path string) (string, error) {
	// The passed path should be /api/{serviceName}/*
	split := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(split) <= 1 {
		return "", fmt.Errorf("failed to parse target service from path: %s", path)
	}

	if split[0] != "api" {
		return "", fmt.Errorf("failed to parse target service from path: %s", path)
	}

	serviceHost := fmt.Sprintf("svc-%s", split[1])
	if serviceHost == "" {
		return "", fmt.Errorf("failed to parse target  from path: %s", path)
	}

	// Return the interal k8s address for the found service
	return fmt.Sprintf(
		"http://%s.svc.cluster.local/%s",
		serviceHost, strings.Join(split[2:], "/"),
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
