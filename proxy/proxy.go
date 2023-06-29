package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/webstrasuite/webstra-gateway/pb"
)

const userIdKey = "userID"

type Proxier interface {
	Handler(authClient pb.AuthServiceClient) echo.HandlerFunc
	ExtractService(string) (string, error)
}

func proxy(authClient pb.AuthServiceClient, p Proxier) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
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
		resp, err := authClient.Validate(c.Request().Context(), &pb.ValidateRequest{Token: "fail"})
		if err != nil {
			return err
		}

		if resp.Status != http.StatusOK {
			c.String(int(resp.Status), resp.Error)
			return fmt.Errorf("user not authenticated")
		}

		// -> Check if the user is allowed to access this service (probably rpc)
		// -> If the user is allowed to access the service, proxy the request to the service
		//    likely will need to inject some user / permissions / role data.

		// potentially a proxy is not necessary and we can make the call as specified directly.

		createReverseProxy(serviceUrl, resp.UserId).ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}

// Should take in user object which it can pass in request headers
func createReverseProxy(address *url.URL, userID int64) *httputil.ReverseProxy {
	p := httputil.NewSingleHostReverseProxy(address)
	p.Director = func(request *http.Request) {
		request.Host = address.Host
		request.URL.Scheme = address.Scheme
		request.URL.Host = address.Host
		request.URL.Path = address.Path
		// Add request header(s) with information about the user like ID and Role
		request.Header.Add(userIdKey, fmt.Sprint(userID))
	}

	// Handle Responses by logging errors or changing status codes in production with p.ModifyResponse
	return p
}
