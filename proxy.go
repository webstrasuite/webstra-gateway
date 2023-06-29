package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
)

type Proxier interface {
	Proxy(echo.Context) error
	ExtractService(string) (string, error)
}

func proxy(p Proxier, ctx echo.Context) error {
	path := ctx.Request().URL.Path
	service, err := p.ExtractService(path)
	if err != nil {
		return err
	}

	serviceUrl, err := url.Parse(service)
	if err != nil {
		return err
	}

	// A few steps need to happen here:
	// -> Log the request (sending something to the logging service (probably rpc))
	// -> Check if the user is authenticated (probably rpc)
	// -> Check if the user is allowed to access this service (probably rpc)
	// -> If the user is allowed to access the service, proxy the request to the service
	//    likely will need to inject some user / permissions / role data.

	// potentially a proxy is not necessary and we can make the call as specified directly.

	createReverseProxy(serviceUrl).ServeHTTP(ctx.Response().Writer, ctx.Request())
	return nil
}

// Should take in user object which it can pass in request headers
func createReverseProxy(address *url.URL) *httputil.ReverseProxy {
	p := httputil.NewSingleHostReverseProxy(address)
	p.Director = func(request *http.Request) {
		request.Host = address.Host
		request.URL.Scheme = address.Scheme
		request.URL.Host = address.Host
		request.URL.Path = address.Path
		// Add request header(s) with information about the user like ID and Role

	}
	log.Println("hit")

	// Handle Responses by logging errors or changing status codes in production with p.ModifyResponse
	return p
}
