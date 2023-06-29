package gateway

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/webstrasuite/webstra-gateway/pb"
	"github.com/webstrasuite/webstra-gateway/proxy"
	"google.golang.org/grpc"
)

type Gateway struct {
	listenAddr      string
	authServiceAddr string
	e               *echo.Echo
	proxy           proxy.Proxier
	authClient      pb.AuthServiceClient
}

func New(addr, authServiceAddr string, proxy proxy.Proxier) (*Gateway, error) {
	// Initialise router
	e := echo.New()

	e.HideBanner = true

	// Logger config to skip logging of healthcheck calls
	loggerConfig := middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health"
		},
	}

	// Register custom logger and standard recovery middleware
	e.Use(middleware.Recover(), middleware.LoggerWithConfig(loggerConfig))

	// Initialize a (gRPC) authentication service client
	authClient, err := initAuthClient(authServiceAddr)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		listenAddr:      addr,
		authServiceAddr: authServiceAddr,
		e:               e,
		proxy:           proxy,
		authClient:      authClient,
	}, nil
}

func (g *Gateway) RegisterRoutes() {
	// Health check endpoint for k8s liveness/readiness
	g.e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Route any other requests through the reverse proxy / gateway
	g.e.Any("/api/*path", g.proxy.Handler(g.authClient))
}

func (g *Gateway) Start() {
	// Start server
	go func() {
		if err := g.e.Start(g.listenAddr); err != nil && err != http.ErrServerClosed {
			g.e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := g.e.Shutdown(ctx); err != nil {
		g.e.Logger.Fatal(err)
	}
}

func initAuthClient(addr string) (pb.AuthServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewAuthServiceClient(conn), nil
}
