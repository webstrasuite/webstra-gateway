package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/webstrasuite/webstra-gateway/proxy"
)

type Router struct {
	listenAddr string
	e          *echo.Echo
	proxy      proxy.Proxier
}

func NewRouter(addr string, proxy proxy.Proxier) *Router {
	// Initialise router
	e := echo.New()

	// Logger config to skip logging of healthcheck calls
	loggerConfig := middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health"
		},
	}

	// Register custom logger and standard recovery middleware
	e.Use(middleware.Recover(), middleware.LoggerWithConfig(loggerConfig))

	return &Router{
		listenAddr: addr,
		e:          e,
		proxy:      proxy,
	}
}

func (r *Router) RegisterRoutes() {
	// Health check endpoint for k8s liveness/readiness
	r.e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Route any other requests through the reverse proxy / gateway
	r.e.Any("/api/*path", r.proxy.Handle)
}

func (r *Router) Start() {
	// Start server
	go func() {
		if err := r.e.Start(r.listenAddr); err != nil && err != http.ErrServerClosed {
			r.e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.e.Shutdown(ctx); err != nil {
		r.e.Logger.Fatal(err)
	}
}
